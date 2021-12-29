package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"

	"reminder/pkg/utils"
)

/*
ReminderData represents whole reminder data-structure
*/
type ReminderData struct {
	User         *User  `json:"user"`
	Notes        Notes  `json:"notes"`
	Tags         Tags   `json:"tags"`
	DataFile     string `json:"data_file"`
	LastBackupAt int64  `json:"last_backup_at"`
	BaseStruct
}

// methods

// update data file
func (reminderData *ReminderData) UpdateDataFile() error {
	// update updated_at field
	// note that updated_at of a whole remider_data object is different
	// from corresponding field of each note
	reminderData.UpdatedAt = utils.CurrentUnixTimestamp()
	// marshal the data
	byteValue, err := json.MarshalIndent(&reminderData, "", "    ")
	utils.PrintErrorIfPresent(err)
	// commit the byte data to file
	err = ioutil.WriteFile(reminderData.DataFile, byteValue, 0755)
	utils.PrintErrorIfPresent(err)
	return err
}

// sort tags in-place and return slugs
func (reminderData *ReminderData) SortedTagSlugs() []string {
	// sort tags in place
	sort.Sort(reminderData.Tags)
	// fetch slugs and return
	return reminderData.Tags.Slugs()
}

// get tag with given slug
func (reminderData *ReminderData) TagFromSlug(slug string) *Tag {
	return reminderData.Tags.FromSlug(slug)
}

// get tags from tagIDs
func (reminderData *ReminderData) TagsFromIds(tagIDs []int) Tags {
	return reminderData.Tags.FromIds(tagIDs)
}

// get tag ids for given group
func (reminderData *ReminderData) TagIdsForGroup(group string) []int {
	return reminderData.Tags.IdsForGroup(group)
}

// get all notes with given tagID and given status
func (reminderData *ReminderData) FindNotesByTagId(tagID int, status string) Notes {
	return reminderData.Notes.WithTagIdAndStatus(tagID, status)
}

// get all notes with given tagSlug and given status
func (reminderData *ReminderData) FindNotesByTagSlug(tagSlug string, status string) Notes {
	tag := reminderData.TagFromSlug(tagSlug)
	return reminderData.FindNotesByTagId(tag.Id, status)
}

// update note's text
func (reminderData *ReminderData) UpateNoteText(note *Note, text string) error {
	err := note.UpdateText(text)
	if err != nil {
		return err
	}
	reminderData.UpdateDataFile()
	fmt.Println("Updated the data file")
	return nil
}

// update note's due date (complete by)
func (reminderData *ReminderData) UpdateNoteCompleteBy(note *Note, text string) error {
	err := note.UpdateCompleteBy(text)
	if err != nil {
		return err
	}
	reminderData.UpdateDataFile()
	fmt.Println("Updated the data file")
	return nil
}

// add note's comment
func (reminderData *ReminderData) AddNoteComment(note *Note, text string) error {
	err := note.AddComment(text)
	if err != nil {
		return err
	}
	reminderData.UpdateDataFile()
	fmt.Println("Updated the data file")
	return nil
}

// update note's tags
func (reminderData *ReminderData) UpdateNoteTags(note *Note, tagIDs []int) error {
	err := note.UpdateTags(tagIDs)
	if err != nil {
		return err
	}
	reminderData.UpdateDataFile()
	fmt.Println("Updated the data file")
	return nil
}

// update note's status
func (reminderData *ReminderData) UpdateNoteStatus(note *Note, status string) error {
	repeatTagIDs := reminderData.TagIdsForGroup("repeat")
	err := note.UpdateStatus(status, repeatTagIDs)
	if err != nil {
		return err
	}
	reminderData.UpdateDataFile()
	fmt.Println("Updated the data file")
	return nil
}

// register basic tags
func (reminderData *ReminderData) RegisterBasicTags() {
	if len(reminderData.Tags) == 0 {
		fmt.Println("Adding tags:")
		basicTags := FBasicTags()
		reminderData.Tags = basicTags
		reminderData.UpdateDataFile()
		fmt.Printf("Added basic tags: %v\n", reminderData.Tags)
	} else {
		fmt.Printf("%v Skipped registering basic tags as tag list is not empty\n", utils.Symbols["warning"])
	}
}

// register a new tag
func (reminderData *ReminderData) NewTagRegistration() (int, error) {
	// collect and ask info about the tag
	tagID := reminderData.nextPossibleTagId()

	promptTagSlug := utils.GeneratePrompt("tag_slug", "")
	promptTagGroup := utils.GeneratePrompt("tag_group", "")

	tag, err := FNewTag(tagID, promptTagSlug, promptTagGroup)

	// validate and save data
	if err != nil {
		utils.PrintErrorIfPresent(err)
		return 0, err
	} else {
		err, _ = reminderData.newTagAppend(tag), tagID
	}
	return tagID, err
}

// (private) get next possible tagID
func (reminderData *ReminderData) nextPossibleTagId() int {
	allTags := reminderData.Tags
	allTagsLen := len(allTags)
	return allTagsLen
}

// (private) append a new tag
func (reminderData *ReminderData) newTagAppend(tag *Tag) error {
	// check if tag's slug is already present
	isNewSlug := true
	for _, existingTag := range reminderData.Tags {
		if existingTag.Slug == tag.Slug {
			isNewSlug = false
		}
	}
	if !isNewSlug {
		fmt.Printf("%v The tag already exists!\n", utils.Symbols["error"])
		return errors.New("Tag Already Exists")
	}
	// go ahead and append
	fmt.Printf("Added Tag: %v\n", *tag)
	reminderData.Tags = append(reminderData.Tags, tag)
	reminderData.UpdateDataFile()
	return nil
}

// register new note
func (reminderData *ReminderData) NewNoteRegistration(tagIDs []int) (*Note, error) {
	// collect info about the note
	if tagIDs == nil {
		// assuming each note with have on average 2 tags
		tagIDs = make([]int, 0, 2)
	}
	promptNoteText := utils.GeneratePrompt("note_text", "")
	note, err := FNewNote(tagIDs, promptNoteText)
	// validate and save data
	if err != nil {
		utils.PrintErrorIfPresent(err)
		return note, err
	} else {
		err = reminderData.newNoteAppend(note)
	}
	return note, nil
}

// (private) append a new note
func (reminderData *ReminderData) newNoteAppend(note *Note) error {
	fmt.Printf("Added Note: %v\n", *note)
	reminderData.Notes = append(reminderData.Notes, note)
	reminderData.UpdateDataFile()
	return nil
}

// return current status
func (reminderData *ReminderData) Stats() string {
	reportTemplate := `
Stats of "{{.DataFile}}"
  - Number of Tags: {{.Tags | len}}
  - Pending Notes: {{.Notes | numPending}}/{{.Notes | numAll}}
`
	funcMap := template.FuncMap{
		"numPending": func(notes Notes) int { return len(notes.WithStatus("pending")) },
		"numAll":     func(notes Notes) int { return len(notes) },
	}
	return utils.TemplateResult(reportTemplate, funcMap, *reminderData)
}

// create timestamped backup
func (reminderData *ReminderData) CreateBackup() string {
	// get backup file name
	ext := path.Ext(reminderData.DataFile)
	dstFile := reminderData.DataFile[:len(reminderData.DataFile)-len(ext)] + "_backup_" + strconv.FormatInt(int64(utils.CurrentUnixTimestamp()), 10) + ext
	lnFile := reminderData.DataFile[:len(reminderData.DataFile)-len(ext)] + "_backup_latest" + ext
	fmt.Printf("Creating backup at %q\n", dstFile)
	// create backup
	byteValue, err := ioutil.ReadFile(reminderData.DataFile)
	utils.PrintErrorIfPresent(err)
	err = ioutil.WriteFile(dstFile, byteValue, 0644)
	utils.PrintErrorIfPresent(err)
	// create alias of latest backup
	fmt.Printf("Creating synlink at %q\n", lnFile)
	executable, _ := exec.LookPath("ln")
	cmd := &exec.Cmd{
		Path:   executable,
		Args:   []string{executable, "-f", dstFile, lnFile},
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
	err = cmd.Run()
	utils.PrintErrorIfPresent(err)
	return dstFile
}

// auto backup
func (reminderData *ReminderData) AutoBackup(gapSecs int64) string {
	var dstFile string
	currentTime := utils.CurrentUnixTimestamp()
	lastBackup := reminderData.LastBackupAt
	gap := currentTime - lastBackup
	fmt.Printf("Automatic Backup Gap = %vs/%vs\n", gap, gapSecs)
	if gap >= gapSecs {
		dstFile = reminderData.CreateBackup()
		reminderData.LastBackupAt = currentTime
		reminderData.UpdateDataFile()
	} else {
		fmt.Printf("Skipping automatic backup\n")
	}
	return dstFile
}

// method (recursive) to ask tagIDs that are to be associated with a note
// it also registers tags for you, if user asks
func (reminderData *ReminderData) AskTagIds(tagIDs []int) []int {
	var err error
	var tagID int
	// ask user to select tag
	optionIndex, _ := utils.AskOption(append(reminderData.SortedTagSlugs(), fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag")
	if optionIndex == -1 {
		return []int{}
	}
	// get tagID
	if optionIndex == len(reminderData.SortedTagSlugs()) {
		// add new tag
		tagID, err = reminderData.NewTagRegistration()
	} else {
		// existing tag selected
		tagID = reminderData.Tags[optionIndex].Id
		err = nil
	}
	// update tagIDs
	if (err == nil) && (!utils.IntPresentInSlice(tagID, tagIDs)) {
		tagIDs = append(tagIDs, tagID)
	}
	// check with user if another tag is to be added
	promtTagAnother := utils.GeneratePrompt("tag_another", "")
	promptText, err := promtTagAnother.Run()
	utils.PrintErrorIfPresent(err)
	promptText = strings.ToLower(promptText)
	nextTag := false
	for _, yes := range []string{"yes", "y"} {
		if yes == promptText {
			nextTag = true
		}
	}
	if nextTag {
		return reminderData.AskTagIds(tagIDs)
	}
	return tagIDs
}

// method to print note and display options
func (reminderData *ReminderData) PrintNoteAndAskOptions(note *Note) string {
	fmt.Print(note.ExternalText(reminderData))
	_, noteOption := utils.AskOption([]string{fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"),
		fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"),
		fmt.Sprintf("%v %v", utils.Symbols["noAction"], "Do nothing"),
		fmt.Sprintf("%v %v", utils.Symbols["upVote"], "Mark as done"),
		fmt.Sprintf("%v %v", utils.Symbols["downVote"], "Mark as pending"),
		fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"),
		fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"),
		fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text")},
		"Select Action")
	fmt.Println("Do you want to update the note?")
	switch noteOption {
	case fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"):
		promptCommment := utils.GeneratePrompt("note_comment", "")
		promptText, err := promptCommment.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.AddNoteComment(note, promptText)
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"):
		return "main-menu"
	case fmt.Sprintf("%v %v", utils.Symbols["noAction"], "Do nothing"):
		fmt.Println("No changes made")
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["upVote"], "Mark as done"):
		_ = reminderData.UpdateNoteStatus(note, "done")
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["downVote"], "Mark as pending"):
		_ = reminderData.UpdateNoteStatus(note, "pending")
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"):
		promptCompleteBy := utils.GeneratePrompt("note_completed_by", "")
		promptText, err := promptCompleteBy.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.UpdateNoteCompleteBy(note, promptText)
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text"):
		promptNoteTextWithDefault := utils.GeneratePrompt("note_text", note.Text)
		promptText, err := promptNoteTextWithDefault.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.UpateNoteText(note, promptText)
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"):
		tagIDs := reminderData.AskTagIds([]int{})
		if len(tagIDs) > 0 {
			reminderData.UpdateNoteTags(note, tagIDs)
			fmt.Print(note.ExternalText(reminderData))
		} else {
			fmt.Printf("%v Skipping updating note with empty tagIDs list\n", utils.Symbols["warning"])
		}
	}
	return "stay"
}

// method (recursively) to print notes interactively
// if tagID >= 0, the passed notes is ignored, otherwise tagID and status
// is used to get notes
func (reminderData *ReminderData) PrintNotesAndAskOptions(notes Notes, tagID int, status string) error {
	// fetch latest notes if `tagID` is valid
	if tagID >= 0 {
		notes = reminderData.FindNotesByTagId(tagID, status)
	}
	// sort notes
	sort.Sort(Notes(notes))
	texts := notes.ExternalTexts(utils.TerminalWidth() - 50)
	// ask user to select a note
	promptText := ""
	if tagID >= 0 {
		promptText = fmt.Sprintf("Select Note (for the tag %v %v)", utils.Symbols["tag"], reminderData.Tags[tagID].Slug)
	} else {
		promptText = fmt.Sprintf("Select Note")
	}
	noteIndex, _ := utils.AskOption(append(texts, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Note")), promptText)
	if noteIndex == -1 {
		return errors.New("The noteIndex is invalid!")
	}
	// create new note or show note options
	if noteIndex == len(texts) {
		// add new note
		if tagID < 0 {
			return errors.New("The passed tagID is invalid!")
		}
		note, err := reminderData.NewNoteRegistration([]int{tagID})
		if err != nil {
			utils.PrintErrorIfPresent(err)
			return err
		}
		var updatedNotes Notes
		updatedNotes = append(updatedNotes, note)
		updatedNotes = append(updatedNotes, notes...)
		reminderData.PrintNotesAndAskOptions(updatedNotes, tagID, status)
	} else {
		// ask options about select note
		note := notes[noteIndex]
		action := reminderData.PrintNoteAndAskOptions(note)
		if action == "stay" {
			reminderData.PrintNotesAndAskOptions(notes, tagID, status)
		}
	}
	return nil
}

// functions

// function to return default data file path
func FDefaultDataFile() string {
	return path.Join(os.Getenv("HOME"), "reminder", "data.json")
}

// function to make sure dataFilePath exists
func FMakeSureFileExists(dataFilePath string) {
	_, err := os.Stat(dataFilePath)
	if err != nil {
		fmt.Printf("Error finding existing data file: %v\n", err)
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("Generating new data file %v.\n", dataFilePath)
			os.MkdirAll(path.Dir(dataFilePath), 0751)
			reminderData := *FBlankReminder()
			reminderData.DataFile = dataFilePath
			err = reminderData.UpdateDataFile()
		}
	}
	utils.PrintErrorIfPresent(err)
}

// function to create blank ReminderData object
func FBlankReminder() *ReminderData {
	fmt.Println("Initializing the data file. Please provide following data.")
	promptUserName := utils.GeneratePrompt("user_name", "")
	name, err := promptUserName.Run()
	utils.PrintErrorIfPresent(err)
	promptUserEmail := utils.GeneratePrompt("user_email", "")
	emailID, err := promptUserEmail.Run()
	utils.PrintErrorIfPresent(err)
	return &ReminderData{
		User:     &User{Name: name, EmailId: emailID},
		Notes:    Notes{},
		Tags:     Tags{},
		DataFile: FDefaultDataFile(),
	}
}

// function to read data file
func FReadDataFile(dataFilePath string) *ReminderData {
	var reminderData ReminderData
	// read byte data from file
	byteValue, err := ioutil.ReadFile(dataFilePath)
	utils.PrintErrorIfPresent(err)
	// parse json data
	err = json.Unmarshal(byteValue, &reminderData)
	utils.PrintErrorIfPresent(err)
	// close the file
	return &reminderData
}
