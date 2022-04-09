package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"reminder/pkg/utils"
	"strings"
)

// functions

// printNoteField function prints the given field of a note.
func printNoteField(fieldName string, fieldValue interface{}) string {
	var strs []string
	fieldDynamicType := fmt.Sprintf("%T", fieldValue)
	if fieldDynamicType == "[]string" {
		items := fieldValue.([]string)
		strs = append(strs, fmt.Sprintf("  |  %12v:\n", fieldName))
		if items != nil {
			for _, v := range items {
				strs = append(strs, fmt.Sprintf("  |  %12v:  %v\n", "", v))
			}
		}
	} else {
		strs = append(strs, fmt.Sprintf("  |  %12v:  %v\n", fieldName, fieldValue))
	}
	return strings.Join(strs, "")
}

// NewNote function provides prompt to registre a new Note.
func NewNote(tagIDs []int, promptNoteText Prompter) (*Note, error) {
	note := &Note{
		Comments:   Comments{},
		Status:     "pending",
		CompleteBy: 0,
		TagIds:     tagIDs,
		BaseStruct: BaseStruct{
			CreatedAt: utils.CurrentUnixTimestamp(),
			UpdatedAt: utils.CurrentUnixTimestamp()},
		// Text:       noteText,
	}
	noteText, err := promptNoteText.Run()
	note.Text = utils.TrimString(noteText)
	if err != nil || strings.Contains(note.Text, "^C") {
		return note, err
	}
	if len(utils.TrimString(note.Text)) == 0 {
		// this should never be encountered because of validation in earlier step
		fmt.Printf("%v Skipping adding note with empty text\n", utils.Symbols["warning"])
		return note, errors.New("Note's text is empty")
	}
	return note, nil
}

// BasicTags function returns an array of basic tags
// which can be used for initial setup of the application.
// Here some of the tags will have special meaning/functionality
// such as repeat-annually and repeat-monthly.
func BasicTags() Tags {
	basicTagsMap := []map[string]string{
		{"slug": "current", "group": ""},
		{"slug": "priority-urgent", "group": "priority"},
		{"slug": "priority-medium", "group": "priority"},
		{"slug": "priority-low", "group": "priority"},
		{"slug": "repeat-annually", "group": "repeat"},
		{"slug": "repeat-monthly", "group": "repeat"},
		{"slug": "tips", "group": "tips"},
	}
	var basicTags Tags
	for index, tagMap := range basicTagsMap {
		tag := Tag{
			Id:    index,
			Slug:  tagMap["slug"],
			Group: tagMap["group"],
			BaseStruct: BaseStruct{
				CreatedAt: utils.CurrentUnixTimestamp(),
				UpdatedAt: utils.CurrentUnixTimestamp()},
		}
		basicTags = append(basicTags, &tag)
	}
	return basicTags
}

// NewTag funciton provides prompt for creating new Tag.
func NewTag(tagID int, promptTagSlug Prompter, promptTagGroup Prompter) (*Tag, error) {
	tag := &Tag{
		Id: tagID,
		BaseStruct: BaseStruct{
			CreatedAt: utils.CurrentUnixTimestamp(),
			UpdatedAt: utils.CurrentUnixTimestamp()},
		// Slug:      tagSlug,
		// Group:     tagGroup,
	}
	// ask for tag slug
	tagSlug, err := promptTagSlug.Run()
	tag.Slug = utils.TrimString(tagSlug)
	tag.Slug = strings.ToLower(tag.Slug)
	// in case of error or Ctrl-c as input, don't create the tag
	if err != nil || strings.Contains(tag.Slug, "^C") {
		return tag, err
	}
	if len(utils.TrimString(tag.Slug)) == 0 {
		// this should never be encountered because of validation in earlier step
		fmt.Printf("%v Skipping adding tag with empty slug\n", utils.Symbols["warning"])
		err := errors.New("Tag's slug is empty")
		return tag, err
	}
	// ask for tag's group
	tagGroup, err := promptTagGroup.Run()
	if err != nil {
		return tag, err
	}
	tag.Group = strings.ToLower(tagGroup)
	// return successful tag
	return tag, nil
}

// DefaultDataFile function returns default data file path.
func DefaultDataFile() string {
	return path.Join(os.Getenv("HOME"), "reminder", "data.json")
}

// MakeSureFileExists function makes sure that the dataFilePath exists.
func MakeSureFileExists(dataFilePath string) error {
	_, err := os.Stat(dataFilePath)
	if err != nil {
		fmt.Printf("Error finding existing data file: %v\n", err)
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("Try generating new data file %v.\n", dataFilePath)
			err := os.MkdirAll(path.Dir(dataFilePath), 0751)
			if err != nil {
				return err
			}
			reminderData := *BlankReminder()
			reminderData.DataFile = dataFilePath
			return reminderData.UpdateDataFile("")
		}
		utils.PrintErrorIfPresent(err)
		return err
	}
	return nil
}

// BlankReminder function creates blank ReminderData object.
func BlankReminder() *ReminderData {
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
		DataFile: DefaultDataFile(),
	}
}

// ReadDataFile function reads data file.
func ReadDataFile(dataFilePath string) *ReminderData {
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
