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

/*
A Notes is a slice of Note objects.

By default it is sorted by its CreatedAt field.
*/
type Notes []*Note

func (c Notes) Len() int           { return len(c) }
func (c Notes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Notes) Less(i, j int) bool { return c[i].UpdatedAt > c[j].UpdatedAt }

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
	searchableText = append(searchableText, note.Text)
	searchableText = append(searchableText, note.Summary)
	searchableText = append(searchableText, strings.Join(commentsText, ""))
	// return searchable text for note a string
	return strings.Join(searchableText, " ")
}

// AddComment adds a new comment to note.
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

// UpdateText updates note's text.
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

// UpdateSummary updates note's summary.
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

// UpdateCompleteBy updates note's due date.
// If input is "nil", the existing due date is cleared.
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

// UpdateTags updates note's tags.
func (note *Note) UpdateTags(tagIDs []int) error {
	note.TagIds = tagIDs
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note with tags")
	// never expecting an error here
	return nil
}

// UpdateStatus updates note's status.
func (note *Note) UpdateStatus(status string, repeatTagIDs []int) error {
	noteIDsWithRepeat := utils.GetCommonMembersIntSlices(note.TagIds, repeatTagIDs)
	if len(noteIDsWithRepeat) != 0 {
		fmt.Printf("%v Update skipped as one of the associated tag is a \"repeat\" group tag \n", utils.Symbols["warning"])
		return nil
	}
	if note.Status == status {
		fmt.Printf("%v Update skipped as there were no changes\n", utils.Symbols["warning"])
		return nil
	}
	note.Status = status
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note")
	return nil
}

// ToggleMainFlag toggles note's main flag.
func (note *Note) ToggleMainFlag() error {
	note.IsMain = !(note.IsMain)
	note.UpdatedAt = utils.CurrentUnixTimestamp()
	fmt.Println("Updated the note's priority")
	return nil
}

// ExternalTexts returns display text (that is, external representation) of list of notes
// with width of each note is truncated to maxStrLen.
// It returns empty []string if there are no notes.
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

// WithStatus filters notes with given status (such as "pending" status).
// It returns empty Notes if no matching Note is found (even when given status doesn't exist).
func (notes Notes) WithStatus(status string) Notes {
	var result Notes
	for _, note := range notes {
		if note.Status == status {
			result = append(result, note)
		}
	}
	return result
}

// OnlyMain filters notes which are set as main.
// It returns empty Notes if no main notes is found.
func (notes Notes) OnlyMain() Notes {
	var result Notes
	for _, note := range notes {
		if note.IsMain {
			result = append(result, note)
		}
	}
	return result
}

// WithTagIdAndStatus returns all notes with given tagID and given status.
// It returns empty Notes if no matching Note is found (even when given tagID or status doesn't exist).
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
