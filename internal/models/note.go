package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"reminder/pkg/utils"
)

/*
Note represents a task (a TO-DO item)

A note can be main or main-note (incidental)
*/
type Note struct {
	Text       string   `json:"text"`
	Comments   Comments `json:"comments"`
	Summary    string   `json:"summary"`
	Status     string   `json:"status"`
	TagIds     []int    `json:"tag_ids"`
	IsMain     bool     `json:"is_priority"`
	CompleteBy int64    `json:"complete_by"`
	BaseStruct
}

type Notes []*Note

func (c Notes) Len() int           { return len(c) }
func (c Notes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Notes) Less(i, j int) bool { return c[i].UpdatedAt > c[j].UpdatedAt }

// provide basic string representation (actually a slice of strings) of a note
// with each element of slice representing certain field of the note
func (note *Note) Strings() []string {
	// allocating 10 members before hand, considering there will be around 10 status fields
	strs := make([]string, 0, 10)
	strs = append(strs, fPrintNoteField("Text", note.Text))
	strs = append(strs, fPrintNoteField("Comments", note.Comments.Strings()))
	strs = append(strs, fPrintNoteField("Summary", note.Summary))
	strs = append(strs, fPrintNoteField("Status", note.Status))
	strs = append(strs, fPrintNoteField("Tags", note.TagIds))
	strs = append(strs, fPrintNoteField("IsMain", note.IsMain))
	strs = append(strs, fPrintNoteField("CompleteBy", utils.UnixTimestampToLongTimeStr(note.CompleteBy)))
	strs = append(strs, fPrintNoteField("CreatedAt", utils.UnixTimestampToLongTimeStr(note.CreatedAt)))
	strs = append(strs, fPrintNoteField("UpdatedAt", utils.UnixTimestampToLongTimeStr(note.UpdatedAt)))
	return strs
}

// print note with its tags slugs
// this is used as final external reprensentation for display of a single note
func (note *Note) ExternalText(reminderData *ReminderData) string {
	var strs []string
	strs = append(strs, fmt.Sprintln("Note Details: -------------------------------------------------"))
	basicStrs := note.Strings()
	// replace tag ids with tag slugs
	tagsStr := fPrintNoteField("Tags", reminderData.TagsFromIds(note.TagIds).Slugs())
	basicStrs[4] = tagsStr
	// create final list of strings
	strs = append(strs, basicStrs...)
	return strings.Join(strs, "")
}

// provide string representation for searching
// we want to perform full text search on Text and Comments of a note
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
	mainFlag := "incidental"
	if note.IsMain {
	  mainFlag = "main"
	}
	filters := fmt.Sprintf("| %-11s | %-8s |", mainFlag, note.Status)
	// get a complete searchable text array for note
	var searchableText []string
	searchableText = append(searchableText, filters)
	searchableText = append(searchableText, note.Text)
	searchableText = append(searchableText, note.Summary)
	searchableText = append(searchableText, strings.Join(commentsText, ""))
	// return searchable text for note a string
	return strings.Join(searchableText, " ")
}

// add new comment to note
func (note *Note) AddComment(text string) error {
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping adding comment with empty text\n", utils.Symbols["warning"])
		return errors.New("Note's comment text is empty")
	}
	comment := &Comment{Text: text, BaseStruct: BaseStruct{CreatedAt: utils.CurrentUnixTimestamp()}}
	note.Comments = append(note.Comments, comment)
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note")
	return nil
}

// update note's text
func (note *Note) UpdateText(text string) error {
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping updating note with empty text\n", utils.Symbols["warning"])
		return errors.New("Note's text is empty")
	}
	note.Text = text
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note")
	return nil
}

// update note's summary
func (note *Note) UpdateSummary(text string) error {
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping updating note with empty summary\n", utils.Symbols["warning"])
		return errors.New("Note's summary is empty")
	}
	note.Summary = text
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note")
	return nil
}

// update note's due date
// if input is "nil", the existing due date is cleared
func (note *Note) UpdateCompleteBy(text string) error {
	// handle edge-case of empty text
	if len(utils.TrimString(text)) == 0 {
		fmt.Printf("%v Skipping updating note with empty due date\n", utils.Symbols["warning"])
		return errors.New("Note's due date is empty")
	}
	// handle edge-case for clearning the existing due date
	if text == "nil" {
		note.CompleteBy = 0
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		fmt.Println("Cleared the due date from the note")
		return nil
	}
	// happy-path
	format := "2-1-2006"
	timeValue, _ := time.Parse(format, text)
	note.CompleteBy = int64(timeValue.Unix())
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note with new due date")
	return nil
}

// update note's tags
func (note *Note) UpdateTags(tagIDs []int) error {
	note.TagIds = tagIDs
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note with tags")
	// never expecting an error here
	return nil
}

// update note's status
func (note *Note) UpdateStatus(status string, repeatTagIDs []int) error {
	noteIDsWithRepeat := utils.GetCommonMembersIntSlices(note.TagIds, repeatTagIDs)
	if len(noteIDsWithRepeat) != 0 {
		fmt.Printf("%v Update skipped as one of the associated tag is a \"repeat\" group tag \n", utils.Symbols["warning"])
	} else if note.Status != status {
		note.Status = status
		note.UpdatedAt = utils.CurrentUnixTimestamp()
		fmt.Println("Updated the note")
	} else {
		fmt.Printf("%v Update skipped as there were no changes\n", utils.Symbols["warning"])
	}
	return nil
}

// toggle note's main flag
func (note *Note) ToggleMain() error {
	note.IsMain = !(note.IsMain)
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note's priority")
	return nil
}

// get display text (that is, external representation) of list of notes
// with width of each note is truncated to maxStrLen
// return empty []string if there are no notes
func (notes Notes) ExternalTexts(maxStrLen int) []string {
	// assuming there are at least (on average) 100s of notes
	allTexts := make([]string, 0, 100)
	for _, note := range notes {
		noteText := note.Text
		if maxStrLen > 0 {
			if len(noteText) > maxStrLen {
				noteText = fmt.Sprintf("%v%v", noteText[0:(maxStrLen-3)], "...")
			}
		}
		noteText = fmt.Sprintf("%*v {C:%02d, S:%v, D:%v}", -maxStrLen, noteText, len(note.Comments), strings.ToUpper(note.Status[0:1]), utils.UnixTimestampToShortTimeStr(note.CompleteBy))
		allTexts = append(allTexts, noteText)
	}
	return allTexts
}

// filter notes with given status (such as "pending" status)
// return empty Notes if no matching Note is found (even when given status doesn't exist)
func (notes Notes) WithStatus(status string) Notes {
	var result Notes
	for _, note := range notes {
		if note.Status == status {
			result = append(result, note)
		}
	}
	return result
}

// filter notes which are set as main
// return empty Notes if no main notes is found
func (notes Notes) OnlyMain() Notes {
	var result Notes
	for _, note := range notes {
		if note.IsMain {
			result = append(result, note)
		}
	}
	return result
}

// get all notes with given tagID and given status
// return empty Notes if no matching Note is found (even when given tagID or status doesn't exist)
func (notes Notes) WithTagIdAndStatus(tagID int, status string) Notes {
	notesWithStatus := notes.WithStatus(status)
	var result Notes
	for _, note := range notesWithStatus {
		if utils.IntPresentInSlice(tagID, note.TagIds) {
			result = append(result, note)
		}
	}
	return result
}

// functions

// function to print given field of a note
func fPrintNoteField(fieldName string, fieldValue interface{}) string {
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

// prompt for new Note
func FNewNote(tagIDs []int, promptNoteText Prompter) (*Note, error) {
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
