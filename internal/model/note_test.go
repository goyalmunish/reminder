package model_test

import (
	"errors"
	"strings"
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
	gc "google.golang.org/api/calendar/v3"
)

func TestNoteStrings(t *testing.T) {
	utils.Location = utils.UTCLocation()
	comments := model.Comments{&model.Comment{Text: "c1:\n- line 1\n\n- line 2\n- line 3 with \" and < characters"}, &model.Comment{Text: "c2"}, &model.Comment{Text: "c3"}}
	note := &model.Note{Text: "dummy text with \" and < characters", Comments: comments, Status: model.NoteStatus_Pending, Summary: "summary heading:\n- line 1\n- line 2", TagIds: []int{1, 2}, CompleteBy: 1609669235}
	want := `[  |          Text:  dummy text with " and < characters
   |      Comments:
  |              :  nil | c1:
  |                     - line 1
  |                     - line 2
  |                     - line 3 with " and < characters
  |              :  nil | c2
  |              :  nil | c3
   |       Summary:  summary heading:
  |                     - line 1
  |                     - line 2
   |        Status:  pending
   |          Tags:  [1 2]
   |        IsMain:  false
   |    CompleteBy:  Sunday, 03-Jan-21 10:20:35 UTC
   |     CreatedAt:  nil
   |     UpdatedAt:  nil
]`
	text, _ := note.Strings()
	utils.AssertEqual(t, text, want)
}

func TestNoteExternalText(t *testing.T) {
	utils.Location = utils.UTCLocation()
	comments := model.Comments{&model.Comment{Text: "c < 1"}, &model.Comment{Text: "c > 2"}, &model.Comment{Text: "c & \" 3"}}
	note := &model.Note{Text: "dummy < > \" text", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1, 2}, CompleteBy: 1609669235}
	var tags model.Tags
	tags = append(tags, &model.Tag{Id: 0, Slug: "tag_0", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 1, Slug: "tag_1", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "tag_2", Group: "tag_group2"})
	reminderData := &model.ReminderData{Tags: tags}
	want := `Note Details: -------------------------------------------------
  |          Text:  dummy < > " text
  |      Comments:
  |              :  nil | c < 1
  |              :  nil | c > 2
  |              :  nil | c & " 3
  |       Summary:  
  |        Status:  pending
  |          Tags:
  |              :  tag_1
  |              :  tag_2
  |        IsMain:  false
  |    CompleteBy:  Sunday, 03-Jan-21 10:20:35 UTC
  |     CreatedAt:  nil
  |     UpdatedAt:  nil
`
	text, _ := note.ExternalText(reminderData)
	utils.AssertEqual(t, text, want)
}

func TestNoteSafeExtText(t *testing.T) {
	utils.Location = utils.UTCLocation()
	comments := model.Comments{&model.Comment{Text: "c < 1"}, &model.Comment{Text: "c > 2"}, &model.Comment{Text: "c & \" 3"}}
	note := &model.Note{Text: "dummy < > \" text", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1, 2}, CompleteBy: 1609669235}
	var tags model.Tags
	tags = append(tags, &model.Tag{Id: 0, Slug: "tag_0", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 1, Slug: "tag_1", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "tag_2", Group: "tag_group2"})
	reminderData := &model.ReminderData{Tags: tags}
	want := `Note Details: -------------------------------------------------
  |          Text:  dummy < > " text
  |       Summary:  
  |        Status:  pending
  |          Tags:
  |              :  tag_1
  |              :  tag_2
  |        IsMain:  false
  |    CompleteBy:  Sunday, 03-Jan-21 10:20:35 UTC
  |     CreatedAt:  nil
  |     UpdatedAt:  nil
`
	text, _ := note.SafeExtText(reminderData)
	utils.AssertEqual(t, text, want)
}

func TestNoteSearchableText(t *testing.T) {
	// case 1
	comments := model.Comments{&model.Comment{Text: "c1"}}
	note := model.Note{Text: "a beautiful cat", Comments: comments, Status: model.NoteStatus_Pending, TagIds: []int{1, 2}, CompleteBy: 1609669231}
	got, _ := note.SearchableText()
	utils.AssertEqual(t, got, "| incidental | pending   | ├ a beautiful cat ┤  [nil | c1]")
	// case 2
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}}
	note = model.Note{Text: "a cute dog", Comments: comments, Status: model.NoteStatus_Done, TagIds: []int{1, 2}, CompleteBy: 1609669232}
	got, _ = note.SearchableText()
	utils.AssertEqual(t, got, "| incidental | done      | ├ a cute dog ┤  [nil | c1, nil | foo bar, nil | c3]")
	// case 3
	comments = model.Comments{}
	note = model.Note{Text: "a cute dog", Comments: comments}
	got, _ = note.SearchableText()
	utils.AssertEqual(t, got, "| incidental |           | ├ a cute dog ┤  [no-comments]")
	// case 4
	comments = model.Comments{}
	note = model.Note{Text: "first line\nsecondline\nthird line", Comments: comments}
	got, _ = note.SearchableText()
	utils.AssertEqual(t, got, "| incidental |           | ├ first line NWL secondline NWL third line ┤  [no-comments]")
	// case 5
	comments = model.Comments{&model.Comment{Text: "c1"}}
	note = model.Note{Text: "a beautiful cat", Comments: comments, Status: model.NoteStatus_Suspended, TagIds: []int{1, 2}, CompleteBy: 1609669231}
	got, _ = note.SearchableText()
	utils.AssertEqual(t, got, "| incidental | suspended | ├ a beautiful cat ┤  [nil | c1]")
}

func TestNoteAddComment(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "1", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// add comments
	// case 1
	err := note1.AddComment("test comment 1")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, len(note1.Comments), 1)
	utils.AssertEqual(t, strings.Contains(note1.Comments[0].Text, "test comment 1"), true)
	// case 2
	err = note1.AddComment("test comment 2")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, len(note1.Comments), 2)
	utils.AssertEqual(t, strings.Contains(note1.Comments[1].Text, "test comment 2"), true)
	// case 3
	err = note1.AddComment("")
	utils.AssertEqual(t, strings.Contains(err.Error(), "Note's comment text is empty"), true)
	utils.AssertEqual(t, len(note1.Comments), 2)
	utils.AssertEqual(t, strings.Contains(note1.Comments[1].Text, "test comment 2"), true)
}

func TestNoteUpdateTags(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// update TagIds
	// case 1
	tagIds := []int{2, 5}
	err := note1.UpdateTags(tagIds)
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.TagIds, tagIds)
	// case 2
	tagIds = []int{}
	err = note1.UpdateTags(tagIds)
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.TagIds, tagIds)
}

func TestNoteUpdateStatus(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// update TagIds
	// case 1
	err := note1.UpdateStatus(model.NoteStatus_Done, []int{1, 2, 3})
	utils.AssertEqual(t, err, errors.New("Note is part of a \"repeat\" group"))
	utils.AssertEqual(t, note1.Status, model.NoteStatus_Pending)
	// case 2
	err = note1.UpdateStatus(model.NoteStatus_Done, []int{5, 6, 7})
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.Status, model.NoteStatus_Done)
	// case 3
	err = note1.UpdateStatus(model.NoteStatus_Pending, []int{5, 6, 7})
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.Status, model.NoteStatus_Pending)
}

func TestNoteUpdateText(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// update text
	// case 1
	err := note1.UpdateText("updated text 1")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.Text, "updated text 1")
	// case 2
	err = note1.UpdateText("")
	utils.AssertEqual(t, strings.Contains(err.Error(), "Note's text is empty"), true)
	utils.AssertEqual(t, note1.Text, "updated text 1")
}

func TestNoteUpdateSummary(t *testing.T) {
	// create notes
	note1 := model.Note{Summary: "original summary", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// update summary
	// case 1
	err := note1.UpdateSummary("updated summary 1")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.Summary, "updated summary 1")
	// case 2
	err = note1.UpdateSummary("")
	utils.AssertEqual(t, strings.Contains(err.Error(), "Note's summary is empty"), true)
	utils.AssertEqual(t, note1.Summary, "updated summary 1")
}

func TestNoteUpdateCompleteBy(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	utils.AssertEqual(t, note1.CompleteBy, 0)
	// update complete_by
	// case 1
	err := note1.UpdateCompleteBy("15-12-2021")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.CompleteBy, 1639526400) // Wed Dec 15 2021 00:00:00 GMT+0000
	// case 2
	err = note1.UpdateCompleteBy("31-12-2022")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.CompleteBy, 1672444800) // Sat Dec 31 2022 00:00:00 GMT+0000
	// case 3
	err = note1.UpdateCompleteBy("nil")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.CompleteBy, 0)
}

func TestNoteRepeatType(t *testing.T) {
	repeatAnnuallyTagId := 3
	repeatMonthlyTagId := 4
	// create notes
	note1 := model.Note{Text: "original text1", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	note2 := model.Note{Text: "original text2", Status: model.NoteStatus_Done, TagIds: []int{3, 5}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	note3 := model.Note{Text: "original text3", Status: model.NoteStatus_Done, TagIds: []int{2, 6}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// assert repeat type
	utils.AssertEqual(t, note1.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), "M")
	utils.AssertEqual(t, note2.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), "A")
	utils.AssertEqual(t, note3.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), "-")
	utils.AssertEqual(t, note3.RepeatType(0, 0), "-")
}

func TestNoteToggleMainFlag(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// update TagIds
	// case 1
	originalPriority := note1.IsMain
	err := note1.ToggleMainFlag()
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, originalPriority != note1.IsMain, true)
	// case 2
	originalPriority = note1.IsMain
	err = note1.ToggleMainFlag()
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, originalPriority != note1.IsMain, true)
}

func TestGoogleCalendarEvent(t *testing.T) {
	tagger := TestTagger{}
	var tests = []struct {
		name          string // has to be string
		note          model.Note
		inputRATID    int
		inputRMTID    int
		inputTimezone string
		inputTagger   model.Tagger
		want          *gc.Event
		wantErr       error
		wantedErr     bool
	}{
		{
			name:          "general case 1",
			note:          model.Note{Text: "original text", Status: model.NoteStatus_Pending, TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}},
			inputRATID:    1,
			inputRMTID:    3,
			inputTimezone: "Australia/Melbourne",
			inputTagger:   tagger,
			want: &gc.Event{
				Summary: "[reminder] original text",
			},
			wantedErr: false,
		},
	}
	for position, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			got, err := subtest.note.GoogleCalendarEvent(subtest.inputRATID, subtest.inputRMTID, subtest.inputTimezone, tagger)
			if (err != nil) != subtest.wantedErr {
				t.Fatalf("GoogleCalendarEvent case %q (position=%d) with input <%+v> returns error <%v>; wantError <%v>", subtest.name, position, subtest.note, err, subtest.wantErr)
			}
			if got.Summary != subtest.want.Summary {
				t.Errorf("GoogleCalendarEvent case %q (position=%d) with input <%+v> returns <%+v>; want <%+v>", subtest.name, position, subtest.note, got, subtest.want)
			}
		})
	}
}
