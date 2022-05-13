package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"reminder/pkg/utils"
)

/*
A Note represents a task.

A note can be main or incidental.
A note can be multiple tags, and a tag can be assocaited with mutiple notes.
*/
type Note struct {
	Text       string   `json:"text"`
	Comments   Comments `json:"comments"`
	Summary    string   `json:"summary"`
	Status     string   `json:"status"`
	TagIds     []int    `json:"tag_ids"`
	IsMain     bool     `json:"is_main"`
	CompleteBy int64    `json:"complete_by"`
	BaseStruct
}

// Type returns type of the note: main or incidental.
func (note *Note) Type() string {
	if note.IsMain {
		return "main"
	}
	return "incidental"
}

// Strings provides basic string representation (as a slice of strings) of a note
// with each element of slice representing certain field of the note.
func (note *Note) Strings() []string {
	// allocating 10 members before hand, considering there will be around 10 status fields
	strs := make([]string, 0, 10)
	strs = append(strs, printNoteField("Text", note.Text))
	strs = append(strs, printNoteField("Comments", note.Comments.Strings()))
	strs = append(strs, printNoteField("Summary", note.Summary))
	strs = append(strs, printNoteField("Status", note.Status))
	strs = append(strs, printNoteField("Tags", note.TagIds))
	strs = append(strs, printNoteField("IsMain", note.IsMain))
	strs = append(strs, printNoteField("CompleteBy", utils.UnixTimestampToLongTimeStr(note.CompleteBy)))
	strs = append(strs, printNoteField("CreatedAt", utils.UnixTimestampToLongTimeStr(note.CreatedAt)))
	strs = append(strs, printNoteField("UpdatedAt", utils.UnixTimestampToLongTimeStr(note.UpdatedAt)))
	return strs
}

// ExternalText prints a note with its tags slugs.
// This is used as final external reprensentation for display of a single note.
func (note *Note) ExternalText(reminderData *ReminderData) string {
	var strs []string
	strs = append(strs, fmt.Sprintln("Note Details: -------------------------------------------------"))
	basicStrs := note.Strings()
	// replace tag ids with tag slugs
	tagsStr := printNoteField("Tags", reminderData.TagsFromIds(note.TagIds).Slugs())
	basicStrs[4] = tagsStr
	// create final list of strings
	strs = append(strs, basicStrs...)
	return strings.Join(strs, "")
}

// SearchableText provides string representation of the object.
// It is used while performing full text search on Text and Comments of a note.
func (note *Note) SearchableText() string {
	// get comments text array for note
	var commentsText []string
	commentsText = append(commentsText, "[")
	if len(note.Comments) == 0 {
		commentsText = append(commentsText, "no-comments")
	} else {
		commentsText = append(commentsText, strings.Join(note.Comments.Strings(), ", "))
	}
	commentsText = append(commentsText, "]")
	// get filters
	filters := fmt.Sprintf("| %-10s | %-7s |", note.Type(), note.Status)
	// get a complete searchable text array for note
	var searchableText []string
	searchableText = append(searchableText, filters)
	searchableText = append(searchableText, fmt.Sprintf("├ %s ┤", note.Text))
	searchableText = append(searchableText, note.Summary)
	searchableText = append(searchableText, strings.Join(commentsText, ""))
	// form a single string
	text := strings.Join(searchableText, " ")
	// address some special characters
	text = strings.ReplaceAll(text, "\n", " NWL ")
	// return searchable text for note a string
	return text
}

// AddComment adds a new comment to note.
func (note *Note) AddComment(text string) error {
	if len(utils.TrimString(text)) == 0 {
		return errors.New("Note's comment text is empty")
	}
	comment := &Comment{Text: text, BaseStruct: BaseStruct{CreatedAt: utils.CurrentUnixTimestamp()}}
	note.Comments = append(note.Comments, comment)
	defer fmt.Println("Added the comment")
	// update the UpdatedAt as well
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	return nil
}

// UpdateTags updates note's tags.
func (note *Note) UpdateTags(tagIDs []int) error {
	note.TagIds = tagIDs
	defer fmt.Println("Updated the note with tags")
	// update the UpdatedAt as well
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	return nil
}

// UpdateStatus updates note's status ("done"/"pending").
// Status of a note tag with repeat tag cannot be mared as "done".
func (note *Note) UpdateStatus(status string, repeatTagIDs []int) error {
	noteIDsWithRepeat := utils.GetCommonMembersIntSlices(note.TagIds, repeatTagIDs)
	if len(noteIDsWithRepeat) != 0 {
		return errors.New("Note is part of a \"repeat\" group")
	}
	if note.Status == status {
		return errors.New("Desired status is same as existing one")
	}
	// happy path
	note.Status = status
	defer fmt.Println("Updated the status")
	// update the UpdatedAt as well
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	return nil
}

// UpdateText updates note's text.
// Once updated, the text cannot be made empty.
func (note *Note) UpdateText(text string) error {
	if len(utils.TrimString(text)) == 0 {
		return errors.New("Note's text is empty")
	}
	// happy path
	note.Text = text
	defer fmt.Println("Updated the text")
	// update the UpdatedAt as well
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	return nil
}

// UpdateSummary updates note's summary.
// If input is "nil", the existing summary is cleared.
func (note *Note) UpdateSummary(text string) error {
	if len(utils.TrimString(text)) == 0 {
		return errors.New("Note's summary is empty")
	}
	// happy path
	if text == "nil" {
		note.Summary = ""
		defer fmt.Println("Cleared the due date from the note")
	} else {
		note.Summary = text
		defer fmt.Println("Updated the summary")
	}
	// update the UpdatedAt as well
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	return nil
}

// UpdateCompleteBy updates note's due date.
// If input is "nil", the existing due date is cleared.
func (note *Note) UpdateCompleteBy(text string) error {
	// handle edge-case of empty text
	if len(utils.TrimString(text)) == 0 {
		return errors.New("Note's due date is empty")
	}
	// happy path
	if text == "nil" {
		note.CompleteBy = 0
		defer fmt.Println("Cleared the due date from the note")
	} else {
		format := "2-1-2006"
		// set current year as year if year part is missing
		timeSplit := strings.Split(text, "-")
		if len(timeSplit) == 2 {
			year, err := utils.YearForDueDateDDMM(text)
			if err != nil {
				return err
			}
			text = fmt.Sprintf("%s-%d", text, year)
		}
		// parse and set the date
		timeValue, _ := time.Parse(format, text)
		note.CompleteBy = int64(timeValue.Unix())
		defer fmt.Println("Updated the note with new due date")
	}
	// update the UpdatedAt as well
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	return nil
}

// RepeatType return - (Not-repeat), A (Annual-Repeat), or M (Monthly-Repeat) string
// representing repeate-type of the note
func (note *Note) RepeatType(repeatAnnuallyTagId int, repeatMonthlyTagId int) string {
	repeat := "-" // non-repeat
	if utils.IntPresentInSlice(repeatAnnuallyTagId, note.TagIds) {
		repeat = "A"
	} else if utils.IntPresentInSlice(repeatMonthlyTagId, note.TagIds) {
		repeat = "M"
	}
	return repeat
}

// ToggleMainFlag toggles note's main flag.
func (note *Note) ToggleMainFlag() error {
	note.IsMain = !(note.IsMain)
	defer fmt.Println("Toggled the note's main/incedental flag")
	// update the UpdatedAt as well
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	return nil
}
