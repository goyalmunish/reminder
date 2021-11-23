package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"

	"reminder/pkg/utils"
)

type ReminderData struct {
	User      *User  `json:"user"`
	Notes     Notes  `json:"notes"`
	Tags      Tags   `json:"tags"`
	DataFile  string `json:"data_file"`
	UpdatedAt int64  `json:"updated_at"`
}

// methods

// method to update data file
func (reminderData *ReminderData) UpdateDataFile() error {
	// update updated_at field
	reminderData.UpdatedAt = utils.CurrentUnixTimestamp()
	// marshal the data
	byteValue, err := json.MarshalIndent(&reminderData, "", "    ")
	utils.PrintErrorIfPresent(err)
	// commit the byte data to file
	err = ioutil.WriteFile(reminderData.DataFile, byteValue, 0755)
	utils.PrintErrorIfPresent(err)
	return err
}

// method to get slugs of all tags
func (reminderData *ReminderData) TagsSlugs() []string {
	// sort tags in place
	sort.Sort(reminderData.Tags)
	// fetch slugs and return
	return reminderData.Tags.Slugs()
}

// method to get tag with given slug
func (reminderData *ReminderData) TagFromSlug(slug string) *Tag {
	return reminderData.Tags.FromSlug(slug)
}

// method to get tags from tagIDs
func (reminderData *ReminderData) TagsFromIds(tagIDs []int) Tags {
	return reminderData.Tags.FromIds(tagIDs)
}

// method to get tag ids for given group
func (reminderData *ReminderData) TagIdsForGroup(group string) []int {
	allTags := reminderData.Tags
	var tagIDs []int
	for _, tag := range allTags {
		if tag.Group == group {
			tagIDs = append(tagIDs, tag.Id)
		}
	}
	return tagIDs
}

// method to get next possible tagID
func (reminderData *ReminderData) NextPossibleTagId() int {
	allTags := reminderData.Tags
	allTagsLen := len(allTags)
	return allTagsLen
}

// method to get all notes with given tagID and given status
func (reminderData *ReminderData) NotesWithTagId(tagID int, status string) Notes {
	allNotes := reminderData.Notes.WithStatus(status)
	var notes Notes
	for _, note := range allNotes {
		if utils.IntPresentInSlice(tagID, note.TagIds) {
			notes = append(notes, note)
		}
	}
	return Notes(notes)
}

// method to register basic tags
func (reminderData *ReminderData) RegisterBasicTags() {
	if len(reminderData.Tags) == 0 {
		fmt.Println("Adding tags:")
		basicTags := FBasicTags()
		reminderData.Tags = basicTags
		reminderData.UpdateDataFile()
		fmt.Printf("Added basic tags: %v\n", reminderData.Tags)
	} else {
		fmt.Printf("%v Skipped registering basic tags as tag list is not empty\n", utils.Symbols["error"])
	}
}

// method to register a new tag
func (reminderData *ReminderData) NewTagRegistration() (error, int) {
	// create a tag object
	tagID := reminderData.NextPossibleTagId()
	tag := FNewTag(tagID)
	return reminderData.NewTagAppend(tag), tagID
}

// method to append a new tag
func (reminderData *ReminderData) NewTagAppend(tag *Tag) error {
	// validations
	// check if tag's slug is empty
	if len(utils.TrimString(tag.Slug)) == 0 {
		fmt.Printf("%v Skipping adding tag with empty slug\n", utils.Symbols["error"])
		return errors.New("Tag's slug is empty")
	}
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
	fmt.Println("Tag: ", *tag)
	reminderData.Tags = append(reminderData.Tags, tag)
	reminderData.UpdateDataFile()
	return nil
}

// method to append a new note
func (reminderData *ReminderData) NewNoteAppend(note *Note) error {
	// validations
	if len(utils.TrimString(note.Text)) == 0 {
		fmt.Printf("%v Skipping adding note with empty text\n", utils.Symbols["error"])
		return errors.New("Note's text is empty")
	}
	// go ahead and append
	fmt.Println("Note: ", *note)
	reminderData.Notes = append(reminderData.Notes, note)
	reminderData.UpdateDataFile()
	return nil
}

// method to add note's comment
func (reminderData *ReminderData) AddNoteComment(note *Note, text string) error {
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping adding comment with empty text\n", utils.Symbols["error"])
		return errors.New("Note's comment text is empty")
	} else {
		text := "(" + strconv.Itoa(int(utils.CurrentUnixTimestamp())) + "): " + text
		note.Comments = append(note.Comments, text)
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		reminderData.UpdateDataFile()
		fmt.Println("Updated the note")
		return nil
	}
}

// method to update note status
func (reminderData *ReminderData) UpdateNoteStatus(note *Note, status string) {
	repeatTagIDs := reminderData.TagIdsForGroup("repeat")
	noteIDsWithRepeat := utils.GetCommonMembersIntSlices(note.TagIds, repeatTagIDs)
	if len(noteIDsWithRepeat) != 0 {
		fmt.Printf("%v Update skipped as one of the associated tag is a \"repeat\" group tag \n", utils.Symbols["error"])
	} else if note.Status != status {
		note.Status = status
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		reminderData.UpdateDataFile()
		fmt.Println("Updated the note")
	} else {
		fmt.Printf("%v Update skipped as there were no changes\n", utils.Symbols["error"])
	}
}

// method to update note text
func (reminderData *ReminderData) UpateNoteText(note *Note, text string) error {
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping updating note with empty text\n", utils.Symbols["error"])
		return errors.New("Note's text is empty")
	} else {
		note.Text = text
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		reminderData.UpdateDataFile()
		fmt.Println("Updated the note")
		return nil
	}
}

// method to update note's due date
// if input is "nil", the existing due date is cleared
func (reminderData *ReminderData) UpdateNoteCompleteBy(note *Note, text string) error {
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping updating note with empty text\n", utils.Symbols["error"])
		return errors.New("Note's due date is empty")
	} else if text == "nil" {
		note.CompleteBy = 0
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		reminderData.UpdateDataFile()
		fmt.Println("Cleared the due date from the note")
		return nil
	} else {
		// fmt.Println(text)
		format := "2006-1-2"
		timeValue, _ := time.Parse(format, text)
		note.CompleteBy = int64(timeValue.Unix())
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		reminderData.UpdateDataFile()
		fmt.Println("Updated the note with new due date")
		return nil
	}
}

// method to update note tags
func (reminderData *ReminderData) UpdateNoteTags(note *Note, tagIDs []int) {
	note.TagIds = tagIDs
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	reminderData.UpdateDataFile()
	fmt.Println("Updated the note")
}

// method (recursive) to ask tagIDs that are to be associated with a note
// it also registers tags for you, if user asks
func (reminderData *ReminderData) AskTagIds(tagIDs []int) []int {
	var err error
	var tagID int
	// make sure reminderData.Tags is sorted
	sort.Sort(reminderData.Tags)
	// ask user to select tag
	optionIndex, _ := utils.AskOption(append(reminderData.TagsSlugs(), fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag")
	if optionIndex == -1 {
		return []int{}
	}
	// get tagID
	if optionIndex == len(reminderData.TagsSlugs()) {
		// add new tag
		err, tagID = reminderData.NewTagRegistration()
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
	prompt := promptui.Prompt{
		Label:    "Add another tag: yes/no (default: no):",
		Validate: utils.ValidateString,
	}
	promptText, err := prompt.Run()
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
	_, noteOption := utils.AskOption([]string{fmt.Sprintf("%v %v", utils.Symbols["noAction"], "Do nothing"),
		fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"),
		fmt.Sprintf("%v %v", utils.Symbols["upVote"], "Mark as done"),
		fmt.Sprintf("%v %v", utils.Symbols["downVote"], "Mark as pending"),
		fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"),
		fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"),
		fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"),
		fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text")},
		"Select Action")
	fmt.Println("Do you want to update the note?")
	switch noteOption {
	case fmt.Sprintf("%v %v", utils.Symbols["noAction"], "Do nothing"):
		fmt.Println("No changes made")
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"):
		return "main-menu"
	case fmt.Sprintf("%v %v", utils.Symbols["upVote"], "Mark as done"):
		reminderData.UpdateNoteStatus(note, "done")
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["downVote"], "Mark as pending"):
		reminderData.UpdateNoteStatus(note, "pending")
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"):
		prompt := promptui.Prompt{
			Label:    "Due Date (format: YYYY-MM-DD), or enter nil to clear existing value",
			Validate: utils.ValidateDateString,
		}
		promptText, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.UpdateNoteCompleteBy(note, promptText)
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"):
		prompt := promptui.Prompt{
			Label:    "New Comment",
			Validate: utils.ValidateNonEmptyString,
		}
		promptText, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.AddNoteComment(note, promptText)
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text"):
		prompt := promptui.Prompt{
			Label:    "New Text",
			Default:  note.Text,
			Validate: utils.ValidateNonEmptyString,
		}
		promptText, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.UpateNoteText(note, promptText)
		fmt.Print(note.ExternalText(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"):
		tagIDs := reminderData.AskTagIds([]int{})
		if len(tagIDs) > 0 {
			reminderData.UpdateNoteTags(note, tagIDs)
			fmt.Print(note.ExternalText(reminderData))
		} else {
			fmt.Printf("%v Skipping updating note with empty tagIDs list\n", utils.Symbols["error"])
		}
	}
	return "stay"
}

// method (recursively) to print notes interactively
func (reminderData *ReminderData) PrintNotesAndAskOptions(notes Notes, tagID int) error {
	// sort notes
	sort.Sort(Notes(notes))
	texts := notes.ExternalTexts(utils.TerminalWidth() - 50)
	// ask user to select a note
	fmt.Println("Note: An added note appears immidiately, but if a note is moved, refresh the list by going to main menu and come back.")
	noteIndex, _ := utils.AskOption(append(texts, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Note")), "Select Note")
	if noteIndex == -1 {
		return errors.New("The noteIndex is invalid!")
	}
	// create new note or show note options
	if noteIndex == len(texts) {
		// add new note
		if tagID >= 0 {
			note := FNewNote([]int{tagID})
			err := reminderData.NewNoteAppend(note)
			if err == nil {
				var updatedNotes Notes
				updatedNotes = append(updatedNotes, note)
				updatedNotes = append(updatedNotes, notes...)
				reminderData.PrintNotesAndAskOptions(updatedNotes, tagID)
			}
			utils.PrintErrorIfPresent(err)
		} else {
			return errors.New("The passed tagID is invalid!")
		}
	} else {
		// ask options about select note
		note := notes[noteIndex]
		action := reminderData.PrintNoteAndAskOptions(note)
		if action == "stay" {
			reminderData.PrintNotesAndAskOptions(notes, tagID)
		}
	}
	return nil
}

func (reminderData *ReminderData) Stats() string {
	var stats []string
	if len(reminderData.Tags) > 0 {
		stats = append(stats, fmt.Sprintf("\nStats of %q\n", reminderData.DataFile))
		stats = append(stats, fmt.Sprintf("%4vNumber of Tags: %v\n", "- ", len(reminderData.Tags)))
		stats = append(stats, fmt.Sprintf("%4vPending Notes: %v/%v\n", "- ", len(reminderData.Notes.WithStatus("pending")), len(reminderData.Notes)))
	}
	stats_str := ""
	for _, elem := range stats {
		stats_str += elem
	}
	return stats_str
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
	prompt := promptui.Prompt{
		Label:    "User Name",
		Validate: utils.ValidateNonEmptyString,
	}
	name, err := prompt.Run()
	utils.PrintErrorIfPresent(err)
	prompt = promptui.Prompt{
		Label:    "User Email",
		Validate: utils.ValidateNonEmptyString,
	}
	emailID, err := prompt.Run()
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
