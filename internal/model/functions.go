package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/goyalmunish/reminder/pkg/utils"
)

// appendMultiLineField prints a multi-line string with first line as its heading
func appendMultiLineField(fieldName, multiLineString string, appendTo []string) []string {
	data := strings.Split(multiLineString, "\n")
	heading, subItems := data[0], data[1:]
	appendTo = append(appendTo, fmt.Sprintf("  |  %12v:  %v\n", fieldName, heading))
	for _, e := range subItems {
		e = strings.TrimSpace(e)
		if e != "" {
			appendTo = append(appendTo, fmt.Sprintf("  |  %18v %v\n", "", e))
		}
	}
	return appendTo
}

// appendSimpleField appends a simple (string or similar) field
func appendSimpleField(fieldName, fieldValue interface{}, appendTo []string) []string {
	appendTo = append(appendTo, fmt.Sprintf("  |  %12v:  %v\n", fieldName, fieldValue))
	return appendTo
}

// printNoteField function prints the given field of a note.
func printNoteField(fieldName string, fieldValue interface{}) string {
	var strs []string
	fieldDynamicType := fmt.Sprintf("%T", fieldValue)
	// construct the display text based on value type
	if fieldDynamicType == "[]string" {
		items := fieldValue.([]string)
		strs = append(strs, fmt.Sprintf("  |  %12v:\n", fieldName))
		for _, v := range items {
			if strings.Contains(v, "\n") {
				// useful for multi-line comments
				strs = appendMultiLineField("", v, strs)
			} else {
				strs = appendSimpleField("", v, strs)
			}
		}
	} else if fieldDynamicType == "string" {
		value := fieldValue.(string)
		if strings.Contains(value, "\n") {
			// useful for multi-line summary
			strs = appendMultiLineField(fieldName, value, strs)
		} else {
			strs = appendSimpleField(fieldName, fieldValue, strs)
		}
	} else {
		strs = appendSimpleField(fieldName, fieldValue, strs)
	}
	return strings.Join(strs, "")
}

// NewNote function provides prompt to register a new Note, and returns its answer.
func NewNote(tagIDs []int, useText string) (*Note, error) {
	var noteText string
	var err error
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
	if useText == "" {
		noteText, err = utils.GeneratePrompt("note_text", "")
		if err != nil {
			return note, err
		}
	} else {
		noteText = useText
	}
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
// Pass useSlug and/or useGroup to use given values instead of prompting user.
func NewTag(tagID int, useSlug string, useGroup string) (*Tag, error) {
	var err error
	var tagSlug string
	var tagGroup string
	tag := &Tag{
		Id: tagID,
		BaseStruct: BaseStruct{
			CreatedAt: utils.CurrentUnixTimestamp(),
			UpdatedAt: utils.CurrentUnixTimestamp()},
		// Slug:      tagSlug,
		// Group:     tagGroup,
	}
	// ask for tag slug
	if useSlug == "" {
		tagSlug, err = utils.GeneratePrompt("tag_slug", "")
		// in case of error or Ctrl-c as input, don't create the tag
		if err != nil || strings.Contains(tag.Slug, "^C") {
			return tag, err
		}
	} else {
		tagSlug = useSlug
	}
	tag.Slug = utils.TrimString(tagSlug)
	tag.Slug = strings.ToLower(tag.Slug)
	if len(utils.TrimString(tag.Slug)) == 0 {
		// this should never be encountered because of validation in earlier step
		fmt.Printf("%v Skipping adding tag with empty slug\n", utils.Symbols["warning"])
		err := errors.New("Tag's slug is empty")
		return tag, err
	}
	// ask for tag's group
	if useGroup == "" {
		tagGroup, err = utils.GeneratePrompt("tag_group", "")
		if err != nil {
			return tag, err
		}
	} else {
		tagGroup = useGroup
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
		utils.PrintError(err)
		return err
	}
	return nil
}

// BlankReminder function creates blank ReminderData object.
func BlankReminder() *ReminderData {
	fmt.Println("Initializing the data file. Please provide following data.")
	name, err := utils.GeneratePrompt("user_name", "")
	utils.PrintError(err)
	emailID, err := utils.GeneratePrompt("user_email", "")
	utils.PrintError(err)
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
	utils.PrintError(err)
	// parse json data
	err = json.Unmarshal(byteValue, &reminderData)
	utils.PrintError(err)
	// close the file
	return &reminderData
}
