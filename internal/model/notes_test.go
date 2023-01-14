package model_test

import (
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func TestNotesExternalTexts(t *testing.T) {
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, "[]", "[]")
	// add notes
	comments := model.Comments{&model.Comment{Text: "c1"}}
	notes = append(notes, &model.Note{Text: "beautiful little cat", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1, 2}, CompleteBy: 1609669231})
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}, &model.Comment{Text: "baz"}}
	notes = append(notes, &model.Note{Text: "cute brown dog", Comments: comments, Status: model.NoteStatus_Done, TagIds: []int{1, 2}, CompleteBy: 1609669232})
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "f b"}, &model.Comment{Text: "c4"}, &model.Comment{Text: "b"}}
	notes = append(notes, &model.Note{Text: "cbd", Comments: comments, Status: model.NoteStatus_Suspended, TagIds: []int{1}, CompleteBy: 1609669235})
	// case 2
	got := notes.ExternalTexts(0, 0, 0)
	want := "[beautiful little cat {R: -, C:01, S:P, D:03-Jan-21} cute brown dog {R: -, C:04, S:D, D:03-Jan-21} cbd {R: -, C:04, S:S, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 3
	got = notes.ExternalTexts(5, 0, 0)
	want = "[be... {R: -, C:01, S:P, D:03-Jan-21} cu... {R: -, C:04, S:D, D:03-Jan-21} cbd   {R: -, C:04, S:S, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 4
	got = notes.ExternalTexts(15, 0, 0)
	want = "[beautiful li... {R: -, C:01, S:P, D:03-Jan-21} cute brown dog  {R: -, C:04, S:D, D:03-Jan-21} cbd             {R: -, C:04, S:S, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 5
	got = notes.ExternalTexts(25, 0, 0)
	want = "[beautiful little cat      {R: -, C:01, S:P, D:03-Jan-21} cute brown dog            {R: -, C:04, S:D, D:03-Jan-21} cbd                       {R: -, C:04, S:S, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
}

func TestNotesWithStatus(t *testing.T) {
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, notes.WithStatus(model.NoteStatus_Pending), model.Notes{})
	// add some notes
	comments := model.Comments{&model.Comment{Text: "c1"}}
	note1 := model.Note{Text: "big fat cat", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1, 2}, CompleteBy: 1609669231}
	notes = append(notes, &note1)
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}}
	note2 := model.Note{Text: "cute brown dog", Comments: comments, Status: model.NoteStatus_Done, TagIds: []int{1, 3}, CompleteBy: 1609669232}
	notes = append(notes, &note2)
	comments = model.Comments{&model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}}
	note3 := model.Note{Text: "little hamster", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1}, CompleteBy: 1609669233}
	notes = append(notes, &note3)
	// case 2 (with an invalid status)
	utils.AssertEqual(t, notes.WithStatus("no-such-status"), model.Notes{})
	// case 3 (with valid status "pending")
	got := notes.WithStatus(model.NoteStatus_Pending)
	want := model.Notes{&note1, &note3}
	utils.AssertEqual(t, got, want)
	// case 4 (with valid status "done")
	got = notes.WithStatus(model.NoteStatus_Done)
	want = model.Notes{&note2}
	utils.AssertEqual(t, got, want)
}

func TestNotesWithCompleteBy(t *testing.T) {
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, notes.WithCompleteBy(), model.Notes{})
	// add some notes
	comments := model.Comments{&model.Comment{Text: "c1"}}
	note1 := model.Note{Text: "big fat cat", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1, 2}, CompleteBy: 1609669231}
	notes = append(notes, &note1)
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}}
	note2 := model.Note{Text: "cute brown dog", Comments: comments, Status: model.NoteStatus_Done, TagIds: []int{1, 3}, CompleteBy: 1609669232}
	notes = append(notes, &note2)
	comments = model.Comments{&model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}}
	note3 := model.Note{Text: "little hamster", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1}}
	notes = append(notes, &note3)
	// case 3 (with only few notes to be filtered in)
	got := notes.WithCompleteBy()
	want := model.Notes{&note1, &note2}
	utils.AssertEqual(t, got, want)
}

func TestNotesOnlyMain(t *testing.T) {
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, notes.OnlyMain(), model.Notes{})
	// add some notes
	comments := model.Comments{&model.Comment{Text: "c1"}}
	note1 := model.Note{Text: "big fat cat", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1, 2}, CompleteBy: 1609669231}
	notes = append(notes, &note1)
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}}
	note2 := model.Note{Text: "cute brown dog", Comments: comments, Status: model.NoteStatus_Done, TagIds: []int{1, 3}, IsMain: true, CompleteBy: 1609669232}
	notes = append(notes, &note2)
	comments = model.Comments{&model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}}
	note3 := model.Note{Text: "little hamster", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1}}
	notes = append(notes, &note3)
	// case 3 (with only few notes to be filtered in)
	got := notes.OnlyMain()
	want := model.Notes{&note2}
	utils.AssertEqual(t, got, want)
}

func TestNotesWithTagIdAndStatus(t *testing.T) {
	// var tags model.Tags
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, notes.WithTagIdAndStatus(2, model.NoteStatus_Pending), model.Notes{})
	// creating tags
	// tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	// tags = append(tags, &tag1)
	// tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	// tags = append(tags, &tag2)
	// tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	// tags = append(tags, &tag3)
	// tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	// tags = append(tags, &tag4)
	// create notes
	note1 := model.Note{Text: "1", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	notes = append(notes, &note1)
	note2 := model.Note{Text: "2", Status: model.NoteStatus_Pending, TagIds: []int{2, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}}
	notes = append(notes, &note2)
	note3 := model.Note{Text: "3", Status: model.NoteStatus_Done, TagIds: []int{2}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}}
	notes = append(notes, &note3)
	note4 := model.Note{Text: "4", Status: model.NoteStatus_Done, TagIds: []int{}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}}
	notes = append(notes, &note4)
	note5 := model.Note{Text: "5", Status: model.NoteStatus_Pending, BaseStruct: model.BaseStruct{UpdatedAt: 1600000005}}
	notes = append(notes, &note5)
	note6 := model.Note{Text: "6", Status: model.NoteStatus_Suspended, TagIds: []int{1}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000006}}
	notes = append(notes, &note6)
	// case 2
	utils.AssertEqual(t, notes.WithTagIdAndStatus(2, model.NoteStatus_Pending), []*model.Note{&note2})
	// case 3
	utils.AssertEqual(t, notes.WithTagIdAndStatus(2, model.NoteStatus_Done), []*model.Note{&note3})
	// case 4
	utils.AssertEqual(t, notes.WithTagIdAndStatus(4, model.NoteStatus_Pending), []*model.Note{&note1, &note2})
	// case 5
	utils.AssertEqual(t, notes.WithTagIdAndStatus(1, model.NoteStatus_Done), []*model.Note{})
	// case 6
	utils.AssertEqual(t, notes.WithTagIdAndStatus(1, model.NoteStatus_Suspended), []*model.Note{&note6})
}
