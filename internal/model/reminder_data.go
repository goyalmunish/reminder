package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goyalmunish/reminder/pkg/calendar"
	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"
	gc "google.golang.org/api/calendar/v3"
)

const EnableCalendar bool = true

/*
A ReminderData represents the whole reminder data-structure.
*/
type ReminderData struct {
	User         *User  `json:"user"`
	Notes        Notes  `json:"notes"`
	Tags         Tags   `json:"tags"`
	DataFile     string `json:"data_file"`
	LastBackupAt int64  `json:"last_backup_at"`
	MutexLock    bool   `json:"mutex_lock"`
	BaseStruct
}

// SyncCalendar syncs pending notes to Cloud Calendar.
func (rd *ReminderData) SyncCalendar(calOptions *calendar.Options) {
	if !EnableCalendar {
		logger.Warn("Google Calendar is disabled.")
		return
	}

	// Get calendar service
	logger.Info("Retrieve the Cloud Calendar Service.")
	srv, err := calendar.GetCalendarService(calOptions)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to retrieve Calendar client: %v", err))
	}

	logger.Info("Fetch the list of all upcoming Calendar Events with each type of recurring event as single unit.")
	// Get list of all upcomming events, with recurring events as a
	// unit (and not as separate single events).
	currentTime := time.Now()
	tStart := currentTime.Format(time.RFC3339)
	tStop := currentTime.AddDate(5, 0, 0).Format(time.RFC3339) // 5 years from now
	existingEvents, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(false).
		TimeMin(tStart).
		TimeMax(tStop).
		MaxResults(250). // 250 is default and is maximum value; if there are more than 250 events, then we'll have to paginate
		Do()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to retrieve the events: %v", err))
	}

	// Get Cloud Calendar details
	fmt.Println(calendar.EventsDetails(existingEvents))

	// Iterating through the Cloud Calendar Events
	fmt.Printf("Listing upcoming calendar events (%v):\n", len(existingEvents.Items))
	fmt.Println("Note: It may take some time for Google Calendar read API to pick up the recently added events, or Calendar trash https://calendar.google.com/calendar/u/0/r/trash can also cause issues.")
	if len(existingEvents.Items) == 0 {
		logger.Warn("No upcoming events found.")
	} else {
		for _, item := range existingEvents.Items {
			if item.Summary == "" {
				continue
			}
			owned := false
			if strings.HasPrefix(item.Summary, calendar.TitlePrefix) {
				owned = true
			}
			fmt.Printf("  - %v | %v | owned=%v\n", item.Summary, item.Recurrence, owned)
		}
	}
	fmt.Println()

	// Note: Your might like to check your Calendar trash
	// https://calendar.google.com/calendar/u/0/r/trash and clear it from
	// time to time. Otherwise, this obstructs with visibility of newly
	// added event in the API call.
	logger.Info("Fetch all the events (max 250) registered by reminder app (and some other matching the query).")
	reminderEvents, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(false).
		Q(calendar.TitlePrefix).
		Do()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to retrieve the events: %v", err))
	}
	fmt.Printf("Listing matching events (%v) and deleting the ones registered by reminder app:\n", len(reminderEvents.Items))
	if len(reminderEvents.Items) == 0 {
		logger.Warn("No registered events found.")
	} else {
		deletionCount := 0
		for _, item := range reminderEvents.Items {
			owned := false
			if strings.HasPrefix(item.Summary, calendar.TitlePrefix) {
				owned = true
			}
			fmt.Printf("  - %q -- %q (owned=%v)\n", calendar.EventString(item), item.Id, owned)
			if owned {
				if calOptions.DryMode {
					logger.Warn(fmt.Sprintf("Dry mode is enabled; skipping deletion of the event %q.", item.Id))
					// continue with next iteration
					continue
				}
				if err := srv.Events.Delete("primary", item.Id).Do(); err != nil {
					logger.Fatal(fmt.Sprintf("Couldn't delete the Calendar event %q | %q | %v", item.Id, item.Summary, err))
				} else {
					deletionCount += 1
					fmt.Printf("    - Deleted the Calendar event %q | %q\n", item.Id, calendar.EventString(item))
				}
			}
		}
		if deletionCount > 0 {
			fmt.Printf("\nWaring! Deletion count is %v; you might like to clear your trash manually by visiting https://calendar.google.com/calendar/u/0/r/trash\n", deletionCount)
		}
	}

	// Add events to Cloud Calendar
	newEvents := rd.GoogleCalendarEvents(existingEvents.TimeZone, rd)
	waitFor := 10 * time.Second
	fmt.Printf("\nSyncing %v events (pending onces with due-date) to Google Calendar. Hit Ctrl-c if you don't want to do it at the moment. The process will wait for %v.\n", len(newEvents), waitFor)
	fmt.Printf("waiting for %v...\n", waitFor)
	time.Sleep(waitFor)
	fmt.Println("Starting the syncing process.")
	for _, event := range newEvents {
		if calOptions.DryMode {
			logger.Warn(fmt.Sprintf("Dry mode is enabled; skipping insertion of event %q.\n", calendar.EventString(event)))
			// continue with next iteration
			continue
		}
		_, err = srv.Events.Insert("primary", event).Do()
		if err != nil {
			logger.Error(err)
		}
		logger.Info(fmt.Sprintf("Synced the event %q.\n", calendar.EventString(event)))
	}
	fmt.Println("Done with syncing process.")
}

// GoogleCalendarEvents returns list of Google Calendar Events.
func (rd *ReminderData) GoogleCalendarEvents(timezoneIANA string, reminderData *ReminderData) []*gc.Event {
	// get all pending notes
	allNotes := rd.Notes
	relevantNotes := allNotes.WithStatus(NoteStatus_Pending).WithCompleteBy()
	// construct Cloud Events
	repeatAnnuallyTag := rd.TagFromSlug("repeat-annually")
	repeatMonthlyTag := rd.TagFromSlug("repeat-monthly")
	var events []*gc.Event
	for _, note := range relevantNotes {
		event := note.GoogleCalendarEvent(repeatAnnuallyTag.Id, repeatMonthlyTag.Id, timezoneIANA, reminderData)
		events = append(events, event)
	}
	return events
}

// CreateDataFile creates data file with current state of `rd`.
// The msg is any additional message to be printed.
func (rd *ReminderData) CreateDataFile(msg string) error {
	// update UpdatedAt field
	// note that UpdatedAt of a whole ReminderData object is different
	// from corresponding field of each note
	currentTime := utils.CurrentUnixTimestamp()
	rd.UpdatedAt = currentTime
	rd.CreatedAt = currentTime
	// marshal the data
	byteValue, err := json.MarshalIndent(&rd, "", "    ")
	if err != nil {
		return err
	}
	// persist the byte data to file
	err = os.WriteFile(rd.DataFile, byteValue, 0755)
	if err != nil {
		return err
	}
	if msg != "" {
		logger.Info(msg)
	}
	logger.Info(fmt.Sprintf("Created the data file at %v!", rd.UpdatedAt))
	return nil
}

// UpdateDataFile updates data file with current state of `rd`.
// The msg is any additional message to be printed.
func (rd *ReminderData) UpdateDataFile(msg string) error {
	var conflictError error
	// if UpdatedAt timestamp of currently loaded data is not same as timestamp persisted on datafile
	// some other process would have updated the data file, which can lead to data inconsistency
	persistedData, err := ReadDataFile(rd.DataFile, true)
	if err != nil {
		return err
	}
	persistedTimestamp := persistedData.UpdatedAt
	inMemoryTimestamp := rd.UpdatedAt
	currentTimestamp := utils.CurrentUnixTimestamp()
	logger.Info(fmt.Sprintf("In-memory timestamp: %v, Persisted timestamp: %v, Current Timestamp (being persisted): %v", inMemoryTimestamp, persistedTimestamp, currentTimestamp))
	if persistedTimestamp != inMemoryTimestamp {
		newFilePath := fmt.Sprintf("%s_CONFLICT_%d", rd.DataFile, currentTimestamp)
		rd.DataFile = newFilePath
		logger.Error(fmt.Sprintf("It seems another instance of the application updated the data file; the data will instead be saved to confict file %q.", newFilePath))
		conflictError = ErrorConflictFile
	}
	// update UpdatedAt field
	// note that UpdatedAt of a whole ReminderData object is different
	// from corresponding field of each note
	rd.UpdatedAt = currentTimestamp
	// marshal the data
	// Refer https://pkg.go.dev/encoding/json#MarshalIndent
	// Note: String values encoded as JSON strings are coerced to valid
	// UTF-8, replacing invalid bytes with Unicode replacement rune. So
	// that the JSON will be safe to embed inside HTML <script> tags, the
	// string is encoded using HTMLEscape.
	// For example, a text such as `comment with < and "` will be written
	// as `"comment with \u003c and \"` but it will read back same as the
	// original string
	byteValue, err := json.MarshalIndent(&rd, "", "    ")
	if err != nil {
		return err
	}
	// persist the byte data to file
	err = os.WriteFile(rd.DataFile, byteValue, 0755)
	if err != nil {
		return err
	}
	if msg != "" {
		logger.Info(msg)
	}
	logger.Info(fmt.Sprintf("Updated the data file %q at %v!", rd.DataFile, rd.UpdatedAt))
	if conflictError != nil {
		return conflictError
	}
	return nil
}

// SortedTagSlugs sorts the tags in-place and return slugs.
// Empty Tags is returned if there are no tags.
func (rd *ReminderData) SortedTagSlugs() []string {
	// sort tags in place
	sort.Sort(rd.Tags)
	// fetch slugs and return
	return rd.Tags.Slugs()
}

// TagFromSlug returns tag with given slug.
func (rd *ReminderData) TagFromSlug(slug string) *Tag {
	return rd.Tags.FromSlug(slug)
}

// TagsFromIds returns tags from tagIDs.
func (rd *ReminderData) TagsFromIds(tagIDs []int) Tags {
	return rd.Tags.FromIds(tagIDs)
}

// TagIdsForGroup gets tag ids for given group.
func (rd *ReminderData) TagIdsForGroup(group string) []int {
	return rd.Tags.IdsForGroup(group)
}

// FindNotesByTagId gets all notes with given tagID and given status.
func (rd *ReminderData) FindNotesByTagId(tagID int, status NoteStatus) Notes {
	return rd.Notes.WithTagIdAndStatus(tagID, status)
}

// FindNotesByTagSlug gets all notes with given tagSlug and given status.
func (rd *ReminderData) FindNotesByTagSlug(tagSlug string, status NoteStatus) Notes {
	tag := rd.TagFromSlug(tagSlug)
	// return empty Notes object for nil `tag`
	if tag == nil {
		return Notes{}
	}
	return rd.FindNotesByTagId(tag.Id, status)
}

// UpdateNoteText updates note's text.
func (rd *ReminderData) UpdateNoteText(note *Note, text string) error {
	err := note.UpdateText(text)
	if err != nil {
		return err
	}
	return rd.UpdateDataFile("")
}

// UpdateNoteSummary updates the note's summary.
func (rd *ReminderData) UpdateNoteSummary(note *Note, text string) error {
	err := note.UpdateSummary(text)
	if err != nil {
		return err
	}
	return rd.UpdateDataFile("")
}

// UpdateNoteCompleteBy updates the note's due date (complete by).
func (rd *ReminderData) UpdateNoteCompleteBy(note *Note, text string) error {
	err := note.UpdateCompleteBy(text)
	if err != nil {
		return err
	}
	return rd.UpdateDataFile("")
}

// AddNoteComment adds note's comment.
func (rd *ReminderData) AddNoteComment(note *Note, text string) error {
	err := note.AddComment(text)
	if err != nil {
		return err
	}
	return rd.UpdateDataFile("")
}

// UpdateNoteTags updates note's tags.
func (rd *ReminderData) UpdateNoteTags(note *Note, tagIDs []int) error {
	err := note.UpdateTags(tagIDs)
	if err != nil {
		return err
	}
	return rd.UpdateDataFile("")
}

// UpdateNoteStatus updates note's status.
func (rd *ReminderData) UpdateNoteStatus(note *Note, status NoteStatus) error {
	repeatTagIDs := rd.TagIdsForGroup("repeat")
	err := note.UpdateStatus(status, repeatTagIDs)
	if err != nil {
		return err
	}
	return rd.UpdateDataFile("")
}

// ToggleNoteMainFlag toggles note's priority.
func (rd *ReminderData) ToggleNoteMainFlag(note *Note) error {
	err := note.ToggleMainFlag()
	if err != nil {
		return err
	}
	return rd.UpdateDataFile("")
}

// RegisterBasicTags registers basic tags.
func (rd *ReminderData) RegisterBasicTags() error {
	if len(rd.Tags) != 0 {
		msg := fmt.Sprintf("%v Skipped registering basic tags as tag list is not empty\n", utils.Symbols["warning"])
		return errors.New(msg)
	}
	basicTags := BasicTags()
	rd.Tags = basicTags
	msg := fmt.Sprintf("Added basic tags: %+v\n", rd.Tags)
	return rd.UpdateDataFile(msg)
}

// ListTags prompts a list of all tags (and their notes underneath).
// Like utils.AskOptions, it prints any encountered error, and returns that error just for information.
func (rd *ReminderData) ListTags() error {
	// function to return a tag sumbol
	// keep different tag symbol for empty tags
	tagSymbol := func(tagSlug string) string {
		PendingNote := rd.FindNotesByTagSlug(tagSlug, NoteStatus_Pending)
		if len(PendingNote) > 0 {
			return utils.Symbols["tag"]
		} else {
			return utils.Symbols["zzz"]
		}
	}
	// get list of tags with their emojis
	// assuming there are at least 20 tags (on average)
	allTagSlugsWithEmoji := make([]string, 0, 20)
	for _, tagSlug := range rd.SortedTagSlugs() {
		allTagSlugsWithEmoji = append(allTagSlugsWithEmoji, fmt.Sprintf("%v %v", tagSymbol(tagSlug), tagSlug))
	}
	// ask user to select a tag
	tagIndex, _, err := utils.AskOption(append(allTagSlugsWithEmoji, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag: ")
	if (err != nil) || (tagIndex == -1) {
		// do nothing, just exit
		return err
	}
	// check if user wants to add a new tag
	if tagIndex == len(rd.SortedTagSlugs()) {
		// add new tag
		_, err = rd.NewTagRegistration()
		if err != nil {
			return err
		}
		return nil
	}
	// operate on the selected a tag, and display both main and non-main notes
	tag := rd.Tags[tagIndex]
	err = rd.PrintNotesAndAskOptions(Notes{}, "pending_tag_notes", tag.Id, "default")
	if err != nil {
		utils.LogError(err)
		// go back to ListTags
		err = rd.ListTags()
		if err != nil {
			return err
		}
	}
	return nil
}

// SearchNotes searches throught all notes.
// Like utils.AskOptions, it prints any encountered error, and returns that error just for information.
func (rd *ReminderData) SearchNotes() error {
	// get texts of all notes
	sort.Sort(rd.Notes)
	allNotes := rd.Notes
	// assuming the search shows 25 items in general
	allTexts := make([]string, 0, 25)
	for _, note := range allNotes {
		allTexts = append(allTexts, note.SearchableText())
	}
	// function to search across notes
	searchNotes := func(filterValue string, optValue string, optIndex int) bool {
		filterValue = strings.ToLower(filterValue)
		noteText := allTexts[optIndex]
		return strings.Contains(strings.ToLower(noteText), filterValue)
	}
	// display prompt
	fmt.Printf("Searching through a total of %v notes:\n", len(allTexts))
	index, err := utils.GenerateNoteSearchSelect(utils.ChopStrings(allTexts, utils.TerminalWidth()-10), searchNotes)
	if err != nil {
		return err
	}
	if index >= 0 {
		note := allNotes[index]
		action := rd.PrintNoteAndAskOptions(note)
		if action == "stay" {
			// no action was selected for the note, go one step back
			err := rd.SearchNotes()
			if err != nil {
				return err
			}
		}
	}
	return err
}

// NotesApprachingDueDate fetches all pending notes which are urgent.
// It accepts view as an argument with "default" or "long" as acceptable values
// Note: NotesApprachingDueDate is dangerous as it manipulates the due date (CompleteBy) date of repeating tags
// which can cause persitence of manupulated dates, if the returned data is persisted.
func (rd *ReminderData) NotesApprachingDueDate(view string) Notes {
	allNotes := rd.Notes
	pendingNotes := allNotes.WithStatus(NoteStatus_Pending)
	// assuming there are at least 100 notes (on average)
	currentNotes := make([]*Note, 0, 100)
	repeatTagIDs := rd.TagIdsForGroup("repeat")
	// populating currentNotes
	for _, note := range pendingNotes {
		noteIDsWithRepeat := utils.GetCommonMembersIntSlices(note.TagIds, repeatTagIDs)
		// first process notes WITHOUT tag with group "repeat"
		// start showing such notes 7 days in advance from their due date, and until they are marked done
		minDay := note.CompleteBy - 7*24*60*60
		if view == "long" {
			minDay = note.CompleteBy - 365*24*60*60
		}
		currentTimestamp := utils.CurrentUnixTimestamp()
		if (len(noteIDsWithRepeat) == 0) && (note.CompleteBy != 0) && (currentTimestamp >= minDay) {
			currentNotes = append(currentNotes, note)
		}
		// check notes with tag with group "repeat"
		// start showing notes with "repeat-annually" 7 days in advance
		// start showing notes with "repeat-monthly" 3 days in advance
		// don't show such notes after their due date is past by 2 day
		if (len(noteIDsWithRepeat) > 0) && (note.CompleteBy != 0) {
			// check for repeat-annually tag
			// note: for the CompleteBy date of the note, we accept only date
			// so, even if there is a time element recorded the the timestamp,
			// we ignore it
			repeatAnnuallyTag := rd.TagFromSlug("repeat-annually")
			repeatMonthlyTag := rd.TagFromSlug("repeat-monthly")
			if (repeatAnnuallyTag != nil) && utils.IntPresentInSlice(repeatAnnuallyTag.Id, note.TagIds) {
				_, noteMonth, noteDay := utils.UnixTimestampToTime(note.CompleteBy).Date()
				noteTimestampCurrent := utils.UnixTimestampForCorrespondingCurrentYear(int(noteMonth), noteDay)
				noteTimestampPrevious := noteTimestampCurrent - 365*24*60*60
				noteTimestampNext := noteTimestampCurrent + 365*24*60*60
				daysBefore := int64(3) // days before to start showing the note
				daysAfter := int64(7)  // days after until to show the note
				if view == "long" {
					daysBefore = int64(365)
				}
				shouldDisplay, matchingTimestamp := utils.IsTimeForRepeatNote(noteTimestampCurrent, noteTimestampPrevious, noteTimestampNext, daysBefore, daysAfter)
				// temporarity update note's timestamp
				note.CompleteBy = matchingTimestamp
				if shouldDisplay {
					currentNotes = append(currentNotes, note)
				}
			}
			// check for repeat-monthly tag
			if (repeatMonthlyTag != nil) && utils.IntPresentInSlice(repeatMonthlyTag.Id, note.TagIds) {
				_, _, noteDay := utils.UnixTimestampToTime(note.CompleteBy).Date()
				noteTimestampCurrent := utils.UnixTimestampForCorrespondingCurrentYearMonth(noteDay)
				noteTimestampPrevious := noteTimestampCurrent - 30*24*60*60
				noteTimestampNext := noteTimestampCurrent + 30*24*60*60
				daysBefore := int64(1) // days beofre to start showing the note
				daysAfter := int64(3)  // days after until to show the note
				if view == "long" {
					daysBefore = int64(31)
				}
				shouldDisplay, matchingTimestamp := utils.IsTimeForRepeatNote(noteTimestampCurrent, noteTimestampPrevious, noteTimestampNext, daysBefore, daysAfter)
				// temporarity update note's timestamp
				note.CompleteBy = matchingTimestamp
				if shouldDisplay {
					currentNotes = append(currentNotes, note)
				}
			}
		}
	}
	// return unsorted list
	return currentNotes
}

// NewTagRegistration registers a new tag.
func (rd *ReminderData) NewTagRegistration() (int, error) {
	// collect and ask info about the tag
	tagID := rd.nextPossibleTagId()

	tag, err := NewTag(tagID, "", "")

	// validate and save data
	if err != nil {
		return 0, err
	} else {
		err, _ = rd.newTagAppend(tag), tagID
	}
	return tagID, err
}

// nextPossibleTagId gets next possible tagID.
func (rd *ReminderData) nextPossibleTagId() int {
	allTags := rd.Tags
	allTagsLen := len(allTags)
	return allTagsLen
}

// newTagAppend appends a new tag.
func (rd *ReminderData) newTagAppend(tag *Tag) error {
	// check if tag's slug is already present
	isNewSlug := true
	for _, existingTag := range rd.Tags {
		if existingTag.Slug == tag.Slug {
			isNewSlug = false
		}
	}
	if !isNewSlug {
		return errors.New("Tag Already Exists")
	}
	// go ahead and append
	logger.Info(fmt.Sprintf("Added Tag: %v\n", *tag))
	rd.Tags = append(rd.Tags, tag)
	return rd.UpdateDataFile("")
}

// NewNoteRegistration registers new note.
// The note is saved to the data file.
func (rd *ReminderData) NewNoteRegistration(tagIDs []int) (*Note, error) {
	// collect info about the note
	if tagIDs == nil {
		// assuming each note with have on average 2 tags
		tagIDs = make([]int, 0, 2)
	}
	note, err := NewNote(tagIDs, "")
	// validate and save data
	if err != nil {
		return note, err
	}
	err = rd.newNoteAppend(note)
	if err != nil {
		return note, err
	}
	return note, nil
}

// newNoteAppend appends a new note.
// The note is saved to the data file.
func (rd *ReminderData) newNoteAppend(note *Note) error {
	logger.Info(fmt.Sprintf("Adding Note: %+v\n", *note))
	rd.Notes = append(rd.Notes, note)
	return rd.UpdateDataFile("")
}

// Stats returns current status.
func (rd *ReminderData) Stats() string {
	reportTemplate := `
Stats of "{{.DataFile}}":
  - Number of Tags:  {{.Tags | len}}
  - Pending Notes:   {{.Notes | numPending}}/{{.Notes | numAll}}
  - Suspended Notes: {{.Notes | numSuspended}}
  - Done Notes:      {{.Notes | numDone}}
`
	funcMap := template.FuncMap{
		"numPending":   func(notes Notes) int { return len(notes.WithStatus(NoteStatus_Pending)) },
		"numSuspended": func(notes Notes) int { return len(notes.WithStatus(NoteStatus_Suspended)) },
		"numDone":      func(notes Notes) int { return len(notes.WithStatus(NoteStatus_Done)) },
		"numAll":       func(notes Notes) int { return len(notes) },
	}
	return utils.TemplateResult(reportTemplate, funcMap, *rd)
}

// CreateBackup creates timestamped backup.
// It returns path of the data file.
// Like utils.AskOptions, it prints any encountered error, but doesn't return the error.
func (rd *ReminderData) CreateBackup() (string, error) {
	// get backup file name
	ext := path.Ext(rd.DataFile)
	dstFile := rd.DataFile[:len(rd.DataFile)-len(ext)] + "_backup_" + strconv.FormatInt(int64(utils.CurrentUnixTimestamp()), 10) + ext
	lnFile := rd.DataFile[:len(rd.DataFile)-len(ext)] + "_backup_latest" + ext
	logger.Info(fmt.Sprintf("Creating backup at %q.\n", dstFile))
	// create backup
	byteValue, err := os.ReadFile(rd.DataFile)
	if err != nil {
		return dstFile, err
	}
	err = os.WriteFile(dstFile, byteValue, 0644)
	if err != nil {
		return dstFile, err
	}
	// create alias of latest backup
	logger.Info(fmt.Sprintf("Creating symlink at %q.\n", lnFile))
	executable, _ := exec.LookPath("ln")
	cmd := &exec.Cmd{
		Path:   executable,
		Args:   []string{executable, "-sf", dstFile, lnFile},
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
	err = cmd.Run()
	if err != nil {
		return dstFile, err
	}
	return dstFile, nil
}

// DisplayDataFile displays the data file.
// Like utils.AskOptions, it prints any encountered error, and returns that error just for information.
func (rd *ReminderData) DisplayDataFile() error {
	fmt.Printf("Printing contents (and if possible, its difference since last backup) of %q:\n", rd.DataFile)
	ext := path.Ext(rd.DataFile)
	lnFile := rd.DataFile[:len(rd.DataFile)-len(ext)] + "_backup_latest" + ext
	err := utils.PerformWhich("wdiff")
	if err != nil {
		fmt.Printf("%v Warning: `wdiff` command is not available\n", utils.Symbols["error"])
		err = utils.PerformCat(rd.DataFile)
	} else {
		err = utils.PerformFilePresence(lnFile)
		if err != nil {
			fmt.Printf("Warning: `%v` file is not available yet\n", lnFile)
			err = utils.PerformCat(rd.DataFile)
		} else {
			err = utils.PerformCwdiff(lnFile, rd.DataFile)
		}
	}
	return err
}

// AutoBackup does auto backup.
func (rd *ReminderData) AutoBackup(gapSecs int64) (string, error) {
	var dstFile string
	currentTime := utils.CurrentUnixTimestamp()
	lastBackup := rd.LastBackupAt
	gap := currentTime - lastBackup
	logger.Info(fmt.Sprintf("Automatic Backup Gap = %vs/%vs\n", gap, gapSecs))
	if gap < gapSecs {
		logger.Info(fmt.Sprintln("Skipping automatic backup."))
		return dstFile, nil
	}
	dstFile, _ = rd.CreateBackup()
	rd.LastBackupAt = currentTime
	return dstFile, rd.UpdateDataFile("")
}

// AskTagIds (recursive) ask tagIDs that are to be associated with a note..
// It also registers tags for you, if user asks.
func (rd *ReminderData) AskTagIds(tagIDs []int) []int {
	var err error
	var tagID int
	// ask user to select tag
	optionIndex, _, _ := utils.AskOption(append(rd.SortedTagSlugs(), fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag: ")
	if optionIndex == -1 {
		return []int{}
	}
	// get tagID
	if optionIndex == len(rd.SortedTagSlugs()) {
		// add new tag
		tagID, err = rd.NewTagRegistration()
	} else {
		// existing tag selected
		tagID = rd.Tags[optionIndex].Id
		err = nil
	}
	// update tagIDs
	if (err == nil) && (!utils.IntPresentInSlice(tagID, tagIDs)) {
		tagIDs = append(tagIDs, tagID)
	}
	// check with user if another tag is to be added
	promptText, err := utils.GeneratePrompt("tag_another", "")
	utils.LogError(err)
	promptText = strings.ToLower(promptText)
	nextTag := false
	for _, yes := range []string{"yes", "y"} {
		if yes == promptText {
			nextTag = true
		}
	}
	if nextTag {
		return rd.AskTagIds(tagIDs)
	}
	return tagIDs
}

// PrintNoteAndAskOptions prints note and display options.
// Like utils.AskOptions, it prints any encountered error, but doesn't returns that error just for information.
// It return string representing workflow direction.
func (rd *ReminderData) PrintNoteAndAskOptions(note *Note) string {
	fmt.Print(note.ExternalText(rd))
	_, noteOption, _ := utils.AskOption([]string{
		fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"),
		fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"),
		fmt.Sprintf("%v %v", utils.Symbols["noAction"], "Do nothing"),
		fmt.Sprintf("%v %v", utils.Symbols["upVote"], "Mark as done"),
		fmt.Sprintf("%v %v", utils.Symbols["zzz"], "Mark as suspended"),
		fmt.Sprintf("%v %v", utils.Symbols["downVote"], "Mark as pending"),
		fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"),
		fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"),
		fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text"),
		fmt.Sprintf("%v %v", utils.Symbols["glossary"], "Update summary"),
		fmt.Sprintf("%v %v", utils.Symbols["hat"], "Toggle main/incidental")},
		"Select Action: ")
	switch noteOption {
	case fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"):
		promptText, err := utils.GeneratePrompt("note_comment", "")
		utils.LogError(err)
		err = rd.AddNoteComment(note, promptText)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"):
		return "main-menu"
	case fmt.Sprintf("%v %v", utils.Symbols["noAction"], "Do nothing"):
		fmt.Println("No changes made")
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["upVote"], "Mark as done"):
		err := rd.UpdateNoteStatus(note, NoteStatus_Done)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["zzz"], "Mark as suspended"):
		err := rd.UpdateNoteStatus(note, NoteStatus_Suspended)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["downVote"], "Mark as pending"):
		err := rd.UpdateNoteStatus(note, NoteStatus_Pending)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"):
		promptText, err := utils.GeneratePrompt("note_completed_by", "")
		utils.LogError(err)
		err = rd.UpdateNoteCompleteBy(note, promptText)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text"):
		promptText, err := utils.GeneratePrompt("note_text", note.Text)
		utils.LogError(err)
		err = rd.UpdateNoteText(note, promptText)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["glossary"], "Update summary"):
		promptText, err := utils.GeneratePrompt("note_summary", note.Summary)
		utils.LogError(err)
		err = rd.UpdateNoteSummary(note, promptText)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	case fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"):
		tagIDs := rd.AskTagIds([]int{})
		if len(tagIDs) > 0 {
			err := rd.UpdateNoteTags(note, tagIDs)
			utils.LogError(err)
			fmt.Print(note.ExternalText(rd))
		} else {
			fmt.Printf("%v Skipping updating note with empty tagIDs list\n", utils.Symbols["warning"])
		}
	case fmt.Sprintf("%v %v", utils.Symbols["hat"], "Toggle main/incidental"):
		err := rd.ToggleNoteMainFlag(note)
		utils.LogError(err)
		fmt.Print(note.ExternalText(rd))
	}
	return "stay"
}

// PrintNotesAndAskOptions (recursively) prints notes interactively.
// In some cases, updated list notes will be fetched, so blank notes can be passed in those cases.
// Unless notes are to be fetched, the passed `status` doesn't make sense, so in such cases it can be passed as "fake".
// Like utils.AskOptions, it prints any encountered error, and returns that error just for information.
// It accepts following values for `display_mode`:
// - "done_notes": fetch only done notes
// - "suspended_notes": fetch only suspended notes
// - "pending_tag_notes": fetch pending notes with given tagID
// - "pending_only_main_notes": fetch pending notes with IsMain set as true
// - "pending_approaching_notes": fetch pending notes with approaching due date
// - "pending_long_view_notes": fetch long-view (52 weeks) of pending notes
// - "passed_notes": use passed notes
func (rd *ReminderData) PrintNotesAndAskOptions(notes Notes, display_mode string, tagID int, sortBy string) error {
	// check if passed notes is to be used or to fetch latest notes
	switch display_mode {
	case "done_notes":
		// ignore the passed notes
		// fetch all the done notes
		notes = rd.Notes.WithStatus(NoteStatus_Done)
		fmt.Printf("A total of %v notes marked as 'done':\n", len(notes))
	case "suspended_notes":
		// ignore the passed notes
		// fetch all the done notes
		notes = rd.Notes.WithStatus(NoteStatus_Suspended)
		fmt.Printf("A total of %v notes marked as 'suspended':\n", len(notes))
	case "pending_tag_notes":
		// this is for listing all notes associated with given tag, with asked status
		// fetch pending notes with given tagID
		notes = rd.FindNotesByTagId(tagID, NoteStatus_Pending)
	case "pending_only_main_notes":
		// this is for listing all main notes, with asked status
		notes = rd.Notes.OnlyMain()
		countAllMain := len(notes)
		notes = notes.WithStatus(NoteStatus_Pending)
		countPendingMain := len(notes)
		fmt.Printf("A total of %v/%v notes flagged as 'main':\n", countPendingMain, countAllMain)
	case "pending_long_view_notes":
		// look ahead a year (52 weeks)
		notes = rd.NotesApprachingDueDate("long")
	case "pending_approaching_notes":
		// this is for listing all notes approaching due date
		// fetch notes approaching due date
		fmt.Println("Note: A note can be in 'pending', 'suspended' or 'done' status.")
		fmt.Println("Note: Notes marked as 'pending' are special and they show up everywhere, whereas notes with other status only show up in 'Search' or under their dedicated menu.")
		fmt.Println("Note: Following are the pending notes with due date:")
		fmt.Println("      - within a week or already crossed (for non repeat-annually or repeat-monthly)")
		fmt.Println("      - within 3 days for repeat-annually and a week post due date (ignoring its year)")
		fmt.Println("      - within 1 day for repeat-monthly and 3 days post due date (ignoring its year and month)")
		fmt.Println("Note: The process may automatically adjust CompleteBy (due-date) with MM-YYYY for monthly repeating notes and YYYY for yearly repeating ones. This is done as part of search algorithm, and it does not impacts on any visibility of those notes.")
		notes = rd.NotesApprachingDueDate("default")
	case "passed_notes":
		// use passed notes
		// this is used wthen this function is called recursively
		// this is useful in circumstances where for example, a note's tags are updated for a note under a tag in which case
		// otherwise the note will immediately disappear if the updated tags list doesn't include the original tag
		fmt.Printf("Note: Using passed notes; the list will not be refreshed immediately!\n")
		fmt.Printf("Note: You must not run multiple instances of the app on same data file!\n")
	default:
		logger.Error("Error: Unreachable code")
	}

	// sort notes
	switch sortBy {
	case "due-date":
		sort.Sort(NotesByDueDate(notes))
	case "default":
		sort.Sort(Notes(notes))
	}
	repeatAnnuallyTag := rd.TagFromSlug("repeat-annually")
	repeatMonthlyTag := rd.TagFromSlug("repeat-monthly")
	texts := notes.ExternalTexts(utils.TerminalWidth()-50, repeatAnnuallyTag.Id, repeatMonthlyTag.Id)

	// ask user to select a note
	promptText := ""
	if tagID >= 0 {
		promptText = fmt.Sprintf("Select Note (for the tag %v %v): ", utils.Symbols["tag"], rd.TagsFromIds([]int{tagID})[0].Slug)
	} else {
		promptText = "Select Note: "
	}
	noteIndex, _, err := utils.AskOption(append(texts, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Note")), promptText)
	if (err != nil) || (noteIndex == -1) {
		return err
	}

	// create new note
	if noteIndex == len(texts) {
		// add new note
		if tagID < 0 {
			return errors.New("The passed tagID is invalid!")
		}
		note, err := rd.NewNoteRegistration([]int{tagID})
		if err != nil {
			return err
		}
		var updatedNotes Notes
		updatedNotes = append(updatedNotes, note)
		updatedNotes = append(updatedNotes, notes...)
		err = rd.PrintNotesAndAskOptions(updatedNotes, "passed_notes", tagID, sortBy)
		if err != nil {
			return err
		}
		return nil
	}

	// ask options about the selected note
	note := notes[noteIndex]
	action := rd.PrintNoteAndAskOptions(note)
	if action == "stay" {
		// no action was selected for the note, go one step back
		err = rd.PrintNotesAndAskOptions(notes, "passed_notes", tagID, sortBy)
		if err != nil {
			return err
		}
	}
	return nil
}
