package model

import (
	"fmt"
	"strings"

	"reminder/pkg/utils"
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
			note.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), len(note.Comments), strings.ToUpper(note.Status[0:1]), utils.UnixTimestampToShortTimeStr(note.CompleteBy))
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
