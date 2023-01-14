package model_test

import (
	"sort"
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func TestNotesByDueDate(t *testing.T) {
	var notes []*model.Note
	notes = append(notes, &model.Note{Text: "1", Status: model.NoteStatus_Pending, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: 1800000003})
	notes = append(notes, &model.Note{Text: "2", Status: model.NoteStatus_Pending, BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}, CompleteBy: 1800000004})
	notes = append(notes, &model.Note{Text: "3", Status: model.NoteStatus_Done, BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}, CompleteBy: 1800000002})
	notes = append(notes, &model.Note{Text: "4", Status: model.NoteStatus_Done, BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}, CompleteBy: 1800000001})
	sort.Sort(model.NotesByDueDate(notes))
	var gotTexts []string
	for _, value := range notes {
		gotTexts = append(gotTexts, value.Text)
	}
	wantTexts := []string{"4", "3", "1", "2"}
	utils.AssertEqual(t, gotTexts, wantTexts)
}
