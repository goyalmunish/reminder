package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"reminder/pkg/utils"
)

type ReminderData struct {
	User      *User   `json:"user"`
	Notes     []*Note `json:"notes"`
	Tags      []*Tag  `json:"tags"`
	UpdatedAt int64   `json:"updated_at"`
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
	email_id, err := prompt.Run()
	utils.PrintErrorIfPresent(err)
	return &ReminderData{
		User:  &User{Name: name, EmailId: email_id},
		Notes: []*Note{},
		Tags:  []*Tag{},
	}
}

// function/methods related to serialization

// function to make sure data_file_path exists
func FMakeSureFileExists(data_file_path string) {
	_, err := os.Stat(data_file_path)
	if os.IsNotExist(err) {
		reminderData := *FBlankReminder()
		err = reminderData.UpdateDataFile(data_file_path)
	}
	utils.PrintErrorIfPresent(err)
}

// function to read data file
func FReadDataFile(data_file_path string) *ReminderData {
	var reminderData ReminderData
	// read byte data from file
	byteValue, err := ioutil.ReadFile(data_file_path)
	utils.PrintErrorIfPresent(err)
	// parse json data
	err = json.Unmarshal(byteValue, &reminderData)
	utils.PrintErrorIfPresent(err)
	// close the file
	return &reminderData
}

// method to update data file
func (reminderData *ReminderData) UpdateDataFile(data_file_path string) error {
	// update updated_at field
	reminderData.UpdatedAt = utils.CurrentUnixTimestamp()
	// marshal the data
	byteValue, err := json.MarshalIndent(&reminderData, "", "    ")
	utils.PrintErrorIfPresent(err)
	// commit the byte data to file
	err = ioutil.WriteFile(data_file_path, byteValue, 0755)
	utils.PrintErrorIfPresent(err)
	return err
}

// methods depending on other models

// method to get slugs of all tags
func (reminderData *ReminderData) TagsSlugs() []string {
	// sort tags in place
	sort.Sort(BySlug(reminderData.Tags))
	// fetch sluts and return
	return FTagsSlugs(reminderData.Tags)
}

// method to get tags from tag_ids
func (reminderData *ReminderData) TagsFromIds(tag_ids []int) []*Tag {
	all_tags := reminderData.Tags
	var tags []*Tag
	for _, tag_id := range tag_ids {
		for _, tag := range all_tags {
			if tag_id == tag.Id {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

// method to get tag with given slug
func (reminderData *ReminderData) TagFromSlug(slug string) *Tag {
	all_tags := reminderData.Tags
	for _, tag := range all_tags {
		if tag.Slug == slug {
			return tag
		}
	}
	return nil
}

// method to get tag ids for given group
func (reminderData *ReminderData) TagIdsForGroup(group string) []int {
	all_tags := reminderData.Tags
	var tag_ids []int
	for _, tag := range all_tags {
		if tag.Group == group {
			tag_ids = append(tag_ids, tag.Id)
		}
	}
	return tag_ids
}

// method to get next possible tag_id
func (reminderData *ReminderData) NextPossibleTagId() int {
	all_tags := reminderData.Tags
	all_tags_len := len(all_tags)
	return all_tags_len
}

// method to get all notes with given tag_id and given status
func (reminderData *ReminderData) NotesWithTagId(tag_id int, status string) []*Note {
	all_notes := FNotesWithStatus(reminderData.Notes, status)
	var notes []*Note
	for _, note := range all_notes {
		if utils.IntInSlice(tag_id, note.TagIds) {
			notes = append(notes, note)
		}
	}
	return notes
}

// method to register basic tags
func (reminderData *ReminderData) RegisterBasicTags() {
	if len(reminderData.Tags) == 0 {
		fmt.Println("Adding tags:")
		basic_tags := FBasicTags()
		reminderData.Tags = basic_tags
		reminderData.UpdateDataFile(DataFile)
		fmt.Println("Added!")
	} else {
		fmt.Printf("%v Skipped registering basic tags as tag list is not empty\n", utils.Symbols["error"])
	}
}

// method to register a new tag
func (reminderData *ReminderData) NewTagRegistration() (error, int) {
	// create a tag object
	tag_id := reminderData.NextPossibleTagId()
	tag := FNewTag(tag_id)
	return reminderData.NewTagAppend(tag), tag_id
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
	is_new_slug := true
	for _, existing_tag := range reminderData.Tags {
		if existing_tag.Slug == tag.Slug {
			is_new_slug = false
		}
	}
	if !is_new_slug {
		fmt.Printf("%v The tag already exists!\n", utils.Symbols["error"])
		return errors.New("Tag Already Exists")
	}
	// go ahead and append
	fmt.Println("Tag: ", *tag)
	reminderData.Tags = append(reminderData.Tags, tag)
	reminderData.UpdateDataFile(DataFile)
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
	reminderData.UpdateDataFile(DataFile)
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
		reminderData.UpdateDataFile(DataFile)
		fmt.Println("Updated the note")
		return nil
	}
}

// method to update note status
func (reminderData *ReminderData) UpdateNoteStatus(note *Note, status string) {
	repeat_tag_ids := reminderData.TagIdsForGroup("repeat")
	note_ids_with_repeat := utils.GetCommonIntMembers(note.TagIds, repeat_tag_ids)
	if len(note_ids_with_repeat) != 0 {
		fmt.Printf("%v Update skipped as one of the associated tag is a \"repeat\" group tag \n", utils.Symbols["error"])
	} else if note.Status != status {
		note.Status = status
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		reminderData.UpdateDataFile(DataFile)
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
		reminderData.UpdateDataFile(DataFile)
		fmt.Println("Updated the note")
		return nil
	}
}

// method to update note's due date
func (reminderData *ReminderData) UpdateNoteCompleteBy(note *Note, text string) error {
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping updating note with empty text\n", utils.Symbols["error"])
		return errors.New("Note's due date is empty")
	} else {
		format := "2006-1-2"
		time_value, _ := time.Parse(format, text)
		note.CompleteBy = int64(time_value.Unix())
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		reminderData.UpdateDataFile(DataFile)
		fmt.Println("Updated the note")
		return nil
	}
}

// method to update note tags
func (reminderData *ReminderData) UpdateNoteTags(note *Note, tag_ids []int) {
	note.TagIds = tag_ids
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	reminderData.UpdateDataFile(DataFile)
	fmt.Println("Updated the note")
}

// method (recursive) to ask tag_ids that are to be associated with a note
// it also registers tags for you, if user asks
func (reminderData *ReminderData) AskTagIds(tag_ids []int) []int {
	var err error
	var tag_id int
	// make sure reminderData.Tags is sorted
	sort.Sort(BySlug(reminderData.Tags))
	// ask user to select tag
	option_index, _ := utils.AskOption(append(reminderData.TagsSlugs(), fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag")
	if option_index == -1 {
		return []int{}
	}
	// get tag_id
	if option_index == len(reminderData.TagsSlugs()) {
		// add new tag
		err, tag_id = reminderData.NewTagRegistration()
	} else {
		// existing tag selected
		tag_id = reminderData.Tags[option_index].Id
		err = nil
	}
	// update tag_ids
	if (err == nil) && (!utils.IntInSlice(tag_id, tag_ids)) {
		tag_ids = append(tag_ids, tag_id)
	}
	// check with user if another tag is to be added
	prompt := promptui.Prompt{
		Label:    "Add another tag (default: no):",
		Validate: utils.ValidateString,
	}
	prompt_text, err := prompt.Run()
	utils.PrintErrorIfPresent(err)
	prompt_text = strings.ToLower(prompt_text)
	next_tag := false
	for _, yes := range []string{"yes", "y"} {
		if yes == prompt_text {
			next_tag = true
		}
	}
	if next_tag {
		return reminderData.AskTagIds(tag_ids)
	}
	return tag_ids
}

// method to print note and display options
func (reminderData *ReminderData) PrintNoteAndAskOptions(note *Note) string {
	fmt.Print(note.StringRepr(reminderData))
	_, note_option := utils.AskOption([]string{fmt.Sprintf("%v %v", utils.Symbols["no_action"], "Do nothing"),
		fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"),
		fmt.Sprintf("%v %v", utils.Symbols["up_vote"], "Mark as done"),
		fmt.Sprintf("%v %v", utils.Symbols["down_vote"], "Mark as pending"),
		fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"),
		fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"),
		fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"),
		fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text")},
		"Select Action")
	fmt.Println("Do you want to update the note?")
	switch note_option {
	case fmt.Sprintf("%v %v", utils.Symbols["no_action"], "Do nothing"):
		fmt.Println("No changes made")
		fmt.Print(note.StringRepr(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["home"], "Exit to main menu"):
		return "main-menu"
	case fmt.Sprintf("%v %v", utils.Symbols["up_vote"], "Mark as done"):
		reminderData.UpdateNoteStatus(note, "done")
		fmt.Print(note.StringRepr(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["down_vote"], "Mark as pending"):
		reminderData.UpdateNoteStatus(note, "pending")
		fmt.Print(note.StringRepr(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["calendar"], "Update due date"):
		prompt := promptui.Prompt{
			Label:    "Due Date (YYYY-MM-DD)",
			Validate: utils.ValidateDateString,
		}
		prompt_text, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.UpdateNoteCompleteBy(note, prompt_text)
		fmt.Print(note.StringRepr(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["comment"], "Add comment"):
		prompt := promptui.Prompt{
			Label:    "New Comment",
			Validate: utils.ValidateNonEmptyString,
		}
		prompt_text, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.AddNoteComment(note, prompt_text)
		fmt.Print(note.StringRepr(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["text"], "Update text"):
		prompt := promptui.Prompt{
			Label:    "New Text",
			Default:  note.Text,
			Validate: utils.ValidateNonEmptyString,
		}
		prompt_text, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		reminderData.UpateNoteText(note, prompt_text)
		fmt.Print(note.StringRepr(reminderData))
	case fmt.Sprintf("%v %v", utils.Symbols["tag"], "Update tags"):
		tag_ids := reminderData.AskTagIds([]int{})
		if len(tag_ids) > 0 {
			reminderData.UpdateNoteTags(note, tag_ids)
			fmt.Print(note.StringRepr(reminderData))
		} else {
			fmt.Printf("%v Skipping updating note with empty tag_ids list\n", utils.Symbols["error"])
		}
	}
	return "stay"
}

// method (recursively) to print notes interactively
func (reminderData *ReminderData) PrintNotesAndAskOptions(notes []*Note, tag_id int) error {
	// sort notes
	sort.Sort(ByUpdatedAt(notes))
	texts := FNotesTexts(notes, 150)
	// ask user to select a note
	fmt.Println("Note: An added note appears immidiately, but if a note is moved, refresh the list by going to main menu and come back.")
	note_index, _ := utils.AskOption(append(texts, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Note")), "Select Note")
	if note_index == -1 {
		return errors.New("The note_index is invalid!")
	}
	// create new note or show note options
	if note_index == len(texts) {
		// add new note
		if tag_id >= 0 {
			note := FNewNote([]int{tag_id})
			err := reminderData.NewNoteAppend(note)
			if err == nil {
				var updated_notes []*Note
				updated_notes = append(updated_notes, note)
				updated_notes = append(updated_notes, notes...)
				reminderData.PrintNotesAndAskOptions(updated_notes, tag_id)
			}
			utils.PrintErrorIfPresent(err)
		} else {
			return errors.New("The passed tag_id is invalid!")
		}
	} else {
		// ask options about select note
		note := notes[note_index]
		action := reminderData.PrintNoteAndAskOptions(note)
		if action == "stay" {
			reminderData.PrintNotesAndAskOptions(notes, tag_id)
		}
	}
	return nil
}
