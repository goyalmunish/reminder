package models

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"

	"reminder/pkg/utils"
)

type Note struct {
	Text       string   `json:"text"`
	Comments   []string `json:"comments"`
	Status     string   `json:"status"`
	TagIds     []int    `json:"tag_ids"`
	CompleteBy int64    `json:"complete_by"`
	CreatedAt  int64    `json:"created_at"`
	UpdatedAt  int64    `json:"updated_at"`
}

// method to provide basic string representation of a note
func (note *Note) String() []string {
	var strs []string
	strs = append(strs, fPrintNoteField("Text", note.Text))
	strs = append(strs, fPrintNoteField("Comments", note.Comments))
	strs = append(strs, fPrintNoteField("Status", note.Status))
	strs = append(strs, fPrintNoteField("Tags", note.TagIds))
	strs = append(strs, fPrintNoteField("CompleteBy", utils.UnixTimestampToLongTimeStr(note.CompleteBy)))
	strs = append(strs, fPrintNoteField("CreatedAt", utils.UnixTimestampToLongTimeStr(note.CreatedAt)))
	strs = append(strs, fPrintNoteField("UpdatedAt", utils.UnixTimestampToLongTimeStr(note.UpdatedAt)))
	return strs
}

// method to print note with its tags slugs
func (note *Note) ExternalRepr(reminderData *ReminderData) string {
	var strs []string
	strs = append(strs, fmt.Sprintln("Note Details: -------------------------------------------------"))
	basicStrs := note.String()
	tagsStr := fPrintNoteField("Tags", FTagsSlugs(reminderData.TagsFromIds(note.TagIds)))
	basicStrs[3] = tagsStr
	strs = append(strs, basicStrs...)
	return strings.Join(strs, "")
}

// method providing string representation for searching
// we want to perform full text search on Text and Comments of a note
func (note *Note) SearchableText() string {
	// get comments text array for note
	var commentsText []string
	commentsText = append(commentsText, "[")
	if len(note.Comments) == 0 {
		commentsText = append(commentsText, "no-comments")
	} else {
		commentsText = append(commentsText, strings.Join(note.Comments, ", "))
	}
	commentsText = append(commentsText, "]")
	// get a complete searchable text array for note
	var searchableText []string
	searchableText = append(searchableText, note.Text)
	searchableText = append(searchableText, strings.Join(commentsText, ""))
	// return searchable text for note a string
	return strings.Join(searchableText, " ")
}

type Notes []*Note

func (c Notes) Len() int           { return len(c) }
func (c Notes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Notes) Less(i, j int) bool { return c[i].UpdatedAt > c[j].UpdatedAt }

// get info-texts
func (notes Notes) Texts(maxStrLen int) []string {
	var allTexts []string
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
func (notes Notes) WithStatus(status string) Notes {
	var result Notes
	for _, note := range notes {
		if note.Status == status {
			result = append(result, note)
		}
	}
	return result
}

// function to print given field of a note
func fPrintNoteField(fieldName string, fieldValue interface{}) string {
	var strs []string
	fieldDynamicType := fmt.Sprintf("%T", fieldValue)
	if fieldDynamicType == "[]string" {
		comments := fieldValue.([]string)
		if comments != nil {
			for _, v := range comments {
				strs = append(strs, fmt.Sprintf("  |  %12v:  %v\n", "", v))
			}
		}
	} else {
		strs = append(strs, fmt.Sprintf("  |  %12v:  %v\n", fieldName, fieldValue))
	}
	return strings.Join(strs, "")
}

// prompt for new Note
func FNewNote(tagIDs []int) *Note {
	prompt := promptui.Prompt{
		Label:    "Note Text",
		Validate: utils.ValidateNonEmptyString,
	}
	noteText, err := prompt.Run()
	utils.PrintErrorIfPresent(err)
	noteText = utils.TrimString(noteText)
	return &Note{
		Text:       noteText,
		Comments:   *new([]string),
		Status:     "pending",
		CompleteBy: 0,
		TagIds:     tagIDs,
		CreatedAt:  utils.CurrentUnixTimestamp(),
		UpdatedAt:  utils.CurrentUnixTimestamp(),
	}
}
