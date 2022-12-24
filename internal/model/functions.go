package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/mail"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"
	"github.com/rivo/tview"
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
func NewNote(ctx context.Context, tagIDs []int, useText string) (*Note, error) {
	var noteText string
	var err error
	note := &Note{
		Comments:   Comments{},
		Status:     NoteStatus_Pending,
		CompleteBy: 0,
		TagIds:     tagIDs,
		BaseStruct: BaseStruct{
			CreatedAt: utils.CurrentUnixTimestamp(),
			UpdatedAt: utils.CurrentUnixTimestamp()},
		// Text:       noteText,
	}
	note.SetContext(ctx)
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
func NewTag(ctx context.Context, tagID int, useSlug string, useGroup string) (*Tag, error) {
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
	tag.SetContext(ctx)
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

// MakeSureFileExists function makes sure that the dataFilePath exists.
func MakeSureFileExists(ctx context.Context, dataFilePath string, askUserInput bool) error {
	_, err := os.Stat(dataFilePath)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("Error finding existing data file: %v\n", err))
		if errors.Is(err, fs.ErrNotExist) {
			logger.Info(ctx, fmt.Sprintf("Trying generating new data file %q.\n", dataFilePath))
			err := os.MkdirAll(path.Dir(dataFilePath), 0751)
			if err != nil {
				return err
			}
			reminderData := *BlankReminder(askUserInput, dataFilePath)
			reminderData.DataFile = dataFilePath
			reminderData.SetContext(ctx)
			reminderData.RegisterBasicTags()
			return reminderData.UpdateDataFile("")
		}
		return err
	}
	return nil
}

// BlankReminder function creates blank ReminderData object.
func BlankReminder(askUserInput bool, dataFilePath string) *ReminderData {
	var name string
	var emailID string
	fmt.Println("Initializing the data file. Please provide following data:")
	app := tview.NewApplication()
	reminderData := &ReminderData{
		User:     &User{Name: name, EmailId: emailID},
		Notes:    Notes{},
		Tags:     Tags{},
		DataFile: dataFilePath,
	}

	if !askUserInput {
		return reminderData
	}

	form := tview.NewForm().
		AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
		AddInputField("Name", "", 20, func(textToCheck string, lastChar rune) bool {
			return unicode.IsLetter(lastChar)
		}, func(text string) {
			name = text
		}).
		AddInputField("Email", "", 20, func(textToCheck string, lastChar rune) bool {
			// validation that needs to run on acceptance function
			var symbol_at rune = '\u0040'
			var symbol_dot rune = '\u002E'
			if unicode.IsLetter(lastChar) || unicode.IsDigit(lastChar) || lastChar == symbol_at || lastChar == symbol_dot {
				return true
			}
			return false
		}, func(text string) {
			emailID = text
		})
	form = form.
		AddButton("Confirm", func() {
			// run validations that can be done only on completed fields
			// if inputs are fine, close the form

			// validate emailID
			emailField := form.GetFormItemByLabel("Email")
			app.SetFocus(emailField)
			_, err := mail.ParseAddress(emailID)
			if err != nil {
				// don't stop
				return
			}
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter details: ").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
	reminderData.User = &User{Name: name, EmailId: emailID}
	return reminderData
}

// ReadDataFile function reads data file as instance of `ReminderData`
func ReadDataFile(ctx context.Context, dataFilePath string) *ReminderData {
	var reminderData ReminderData
	// read byte data from file
	byteValue, err := os.ReadFile(dataFilePath)
	utils.LogError(ctx, err)
	// parse json data
	err = json.Unmarshal(byteValue, &reminderData)
	utils.LogError(ctx, err)
	// set context
	reminderData.SetContext(ctx)
	for _, note := range reminderData.Notes {
		note.SetContext(ctx)
	}
	for _, tag := range reminderData.Tags {
		tag.SetContext(ctx)
	}
	// close the file
	return &reminderData
}
