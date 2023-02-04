package model

import (
	"fmt"
	"strings"

	"github.com/goyalmunish/reminder/pkg/utils"
)

/*
A Notes is a slice of Note objects.

By default it is sorted by its CreatedAt field.
*/
type Notes []*Note

func (c Notes) Len() int           { return len(c) }
func (c Notes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Notes) Less(i, j int) bool { return c[i].UpdatedAt > c[j].UpdatedAt }

// ExternalTexts returns display text (that is, external representation) of list of notes
// with width of each note is truncated to maxStrLen.
// It returns empty []string if there are no notes.
// Note: You may use repeatAnnuallyTagId and repeatMonthlyTagId as 0, if they are not required
// In the output: R means "repeat-type", C means "number of comments", S means "status", and D means "due date"
func (notes Notes) ExternalTexts(maxStrLen int, repeatAnnuallyTagId int, repeatMonthlyTagId int) []string {
	// assuming there are at least (on average) 100s of notes
	allTexts := make([]string, 0, 100)
	for _, note := range notes {
		noteText := note.Text
		if maxStrLen > 0 {
			if len(noteText) > maxStrLen {
				noteText = fmt.Sprintf("%v%v", noteText[0:(maxStrLen-3)], "...")
			}
		}
		noteText = fmt.Sprintf(
			"%*v {R: %s, C:%02d, S:%v, D:%v}", -maxStrLen, noteText,
			note.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), len(note.Comments), strings.ToUpper(string(note.Status)[0:1]), utils.UnixTimestampToShortTimeStr(note.CompleteBy))
		allTexts = append(allTexts, noteText)
	}
	return allTexts
}

// PopulateTempDueDate popultes tempDueDate field of note from its CompleteBy field.
func (notes Notes) PopulateTempDueDate() {
	for _, note := range notes {
		note.tempDueDate = note.CompleteBy
	}
}

// WithStatus filters-in notes with given status (such as "pending" status).
// It returns empty Notes if no matching Note is found (even when given status doesn't exist).
func (notes Notes) WithStatus(status NoteStatus) Notes {
	var result Notes
	for _, note := range notes {
		if note.Status == status {
			result = append(result, note)
		}
	}
	return result
}

// WithCompleteBy filters-in only notes with non-nil CompleteBy filed of the notes.
// It returns empty Notes if no matching Note is found (even when given status doesn't exist).
func (notes Notes) WithCompleteBy() Notes {
	var result Notes
	for _, note := range notes {
		if note.CompleteBy > 0 {
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
func (notes Notes) WithTagIdAndStatus(tagID int, status NoteStatus) Notes {
	notesWithStatus := notes.WithStatus(status)
	var result Notes
	for _, note := range notesWithStatus {
		if utils.IsMemberOfSlice(tagID, note.TagIds) {
			result = append(result, note)
		}
	}
	return result
}
