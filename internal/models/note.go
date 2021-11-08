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

type NotesByUpdatedAt []*Note

func (c NotesByUpdatedAt) Len() int           { return len(c) }
func (c NotesByUpdatedAt) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c NotesByUpdatedAt) Less(i, j int) bool { return c[i].UpdatedAt > c[j].UpdatedAt }

// method to provide basic string representation of a note
func (note *Note) String() []string {
	var strs []string
	strs = append(strs, FPrintNoteField("Text", note.Text))
	strs = append(strs, FPrintNoteField("Comments", note.Comments))
	strs = append(strs, FPrintNoteField("Status", note.Status))
	strs = append(strs, FPrintNoteField("Tags", note.TagIds))
	strs = append(strs, FPrintNoteField("CompleteBy", utils.UnixTimestampToLongTimeStr(note.CompleteBy)))
	strs = append(strs, FPrintNoteField("CreatedAt", utils.UnixTimestampToLongTimeStr(note.CreatedAt)))
	strs = append(strs, FPrintNoteField("UpdatedAt", utils.UnixTimestampToLongTimeStr(note.UpdatedAt)))
	return strs
}

// method to print note with its tags slugs
func (note *Note) StringRepr(reminderData *ReminderData) string {
	var strs []string
	strs = append(strs, fmt.Sprintln("Note Details: -------------------------------------------------"))
	basic_strs := note.String()
	tags_str := FPrintNoteField("Tags", FTagsSlugs(reminderData.TagsFromIds(note.TagIds)))
	basic_strs[3] = tags_str
	strs = append(strs, basic_strs...)
	return strings.Join(strs, "")
}

// method providing string representation for searching
// we want to full text search on Text and Comments of a note
func (note *Note) SearchText() string {
	// get comments text array for note
	var comments_text []string
	comments_text = append(comments_text, "[")
	if len(note.Comments) == 0 {
		comments_text = append(comments_text, "no-comments")
	} else {
		comments_text = append(comments_text, strings.Join(note.Comments, ", "))
	}
	comments_text = append(comments_text, "]")
	// get complete search text array for note
	var search_text []string
	search_text = append(search_text, note.Text)
	search_text = append(search_text, strings.Join(comments_text, ""))
	// return search text as for note a string
	return strings.Join(search_text, " ")
}

// function to print given field of a note
func FPrintNoteField(field_name string, field_value interface{}) string {
	var strs []string
	field_dynamic_type := fmt.Sprintf("%T", field_value)
	if field_dynamic_type == "[]string" {
		comments := field_value.([]string)
		if comments != nil {
			for _, v := range comments {
				strs = append(strs, fmt.Sprintf("  |  %12v:  %v\n", "", v))
			}
		}
	} else {
		strs = append(strs, fmt.Sprintf("  |  %12v:  %v\n", field_name, field_value))
	}
	return strings.Join(strs, "")
}

// get info texts of given notes
func FNotesTexts(notes []*Note, max_str_len int) []string {
	var all_texts []string
	for _, note := range notes {
		note_text := note.Text
		if max_str_len > 0 {
			if len(note_text) > max_str_len {
				note_text = fmt.Sprintf("%v%v", note_text[0:(max_str_len-3)], "...")
			}
		}
		note_text = fmt.Sprintf("%*v {C:%02d, S:%v, D:%v}", -max_str_len, note_text, len(note.Comments), strings.ToUpper(note.Status[0:1]), utils.UnixTimestampToShortTimeStr(note.CompleteBy))
		all_texts = append(all_texts, note_text)
	}
	return all_texts
}

// filter notes with given status (such as "pending" status)
func FNotesWithStatus(notes []*Note, status string) []*Note {
	var result []*Note
	for _, note := range notes {
		if note.Status == status {
			result = append(result, note)
		}
	}
	return result
}

// prompt for new Note
func FNewNote(tag_ids []int) *Note {
	prompt := promptui.Prompt{
		Label:    "Note Text",
		Validate: utils.ValidateNonEmptyString,
	}
	note_text, err := prompt.Run()
	utils.PrintErrorIfPresent(err)
	note_text = utils.TrimString(note_text)
	return &Note{
		Text:       note_text,
		Comments:   *new([]string),
		Status:     "pending",
		CompleteBy: 0,
		TagIds:     tag_ids,
		CreatedAt:  utils.CurrentUnixTimestamp(),
		UpdatedAt:  utils.CurrentUnixTimestamp(),
	}
}
