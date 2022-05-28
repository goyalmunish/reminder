package model_test

import (
	"errors"

	// "fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
	"time"

	// "github.com/golang/mock/gomock"

	model "reminder/internal/model"
	utils "reminder/pkg/utils"
)

// mocks

type MockPromptTagSlug struct {
}
type MockPromptTagGroup struct {
}
type MockPromptNoteText struct {
}

func (prompt *MockPromptTagSlug) Run() (string, error) {
	return "test_tag_slug", nil
}

func (prompt *MockPromptTagGroup) Run() (string, error) {
	return "test_tag_group", nil
}

func (prompt *MockPromptNoteText) Run() (string, error) {
	return "a random note text", nil
}

// test examples

func TestDataFile(t *testing.T) {
	defaultDataFilePath := model.DefaultDataFile()
	utils.AssertEqual(t, strings.HasPrefix(defaultDataFilePath, "/"), true)
	utils.AssertEqual(t, strings.HasSuffix(defaultDataFilePath, ".json"), true)
}

func TestUser(t *testing.T) {
	got := model.User{Name: "Test User", EmailId: "user@test.com"}
	want := "{Name: Test User, EmailId: user@test.com}"
	utils.AssertEqual(t, got, want)
}

func TestTag(t *testing.T) {
	// case 1: general case
	got := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	want := "tag_group1#a#1"
	utils.AssertEqual(t, got, want)
	// case 2: blank group
	got = model.Tag{Id: 1, Slug: "a", Group: ""}
	want = "#a#1"
	utils.AssertEqual(t, got, want)
	// case 3: omitted group
	got = model.Tag{Id: 1, Slug: "a"}
	want = "#a#1"
	utils.AssertEqual(t, got, want)
}

func TestTags(t *testing.T) {
	var tags model.Tags
	tags = append(tags, &model.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	sort.Sort(model.Tags(tags))
	var got []int
	for _, value := range tags {
		got = append(got, value.Id)
	}
	want := []int{1, 4, 3, 2}
	utils.AssertEqual(t, got, want)
}

func TestSlugs(t *testing.T) {
	var tags model.Tags
	utils.AssertEqual(t, tags, "[]")
	// case 1 (no tags)
	utils.AssertEqual(t, tags.Slugs(), "[]")
	// case 2 (non-empty tags)
	tags = append(tags, &model.Tag{Id: 1, Slug: "tag_1", Group: "tag_group"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "tag_2", Group: "tag_group"})
	tags = append(tags, &model.Tag{Id: 3, Slug: "tag_3", Group: "tag_group"})
	got := tags.Slugs()
	want := "[tag_1 tag_2 tag_3]"
	utils.AssertEqual(t, got, want)
}

func TestFromSlug(t *testing.T) {
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1 (passing non-existing slug)
	utils.AssertEqual(t, tags.FromSlug("no_such_slug"), nil)
	// case 2 (passing tag which is part of another tag as well)
	utils.AssertEqual(t, tags.FromSlug("a"), &tag1)
	// case 3
	utils.AssertEqual(t, tags.FromSlug("a1"), &tag2)
}

func TestFromIds(t *testing.T) {
	var tags model.Tags
	// creating tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "z", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "c", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1 (passing blank tagIDs)
	tagIDs := []int{}
	gotSlugs := tags.FromIds(tagIDs)
	wantSlugs := model.Tags{}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 2 (no matching tagIDs)
	tagIDs = []int{100, 101}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = model.Tags{}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 3 (two matching tagIDs)
	tagIDs = []int{1, 3}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = model.Tags{&tag1, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 4
	tagIDs = []int{1, 4, 2, 3}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = model.Tags{&tag1, &tag4, &tag2, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
}

func TestIdsForGroup(t *testing.T) {
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1 (group with no such name)
	utils.AssertEqual(t, tags.IdsForGroup("tag_group_NO"), []int{})
	// case 1 (group with multiple tags)
	utils.AssertEqual(t, tags.IdsForGroup("tag_group1"), []int{1, 2, 3})
	// case 2 (group with single tag)
	utils.AssertEqual(t, tags.IdsForGroup("tag_group2"), []int{4})
}

func TestBasicTags(t *testing.T) {
	basicTags := model.BasicTags()
	slugs := basicTags.Slugs()
	want := "[current priority-urgent priority-medium priority-low repeat-annually repeat-monthly tips]"
	utils.AssertEqual(t, slugs, want)
}

func TestNotes(t *testing.T) {
	var notes []*model.Note
	notes = append(notes, &model.Note{Text: "1", Status: "pending", BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}})
	notes = append(notes, &model.Note{Text: "2", Status: "pending", BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}})
	notes = append(notes, &model.Note{Text: "3", Status: "done", BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}})
	notes = append(notes, &model.Note{Text: "4", Status: "done", BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}})
	sort.Sort(model.Notes(notes))
	var gotTexts []string
	for _, value := range notes {
		gotTexts = append(gotTexts, value.Text)
	}
	wantTexts := []string{"2", "3", "4", "1"}
	utils.AssertEqual(t, gotTexts, wantTexts)
}

func TestNotesByDueDate(t *testing.T) {
	var notes []*model.Note
	notes = append(notes, &model.Note{Text: "1", Status: "pending", BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: 1800000003})
	notes = append(notes, &model.Note{Text: "2", Status: "pending", BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}, CompleteBy: 1800000004})
	notes = append(notes, &model.Note{Text: "3", Status: "done", BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}, CompleteBy: 1800000002})
	notes = append(notes, &model.Note{Text: "4", Status: "done", BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}, CompleteBy: 1800000001})
	sort.Sort(model.NotesByDueDate(notes))
	var gotTexts []string
	for _, value := range notes {
		gotTexts = append(gotTexts, value.Text)
	}
	wantTexts := []string{"4", "3", "1", "2"}
	utils.AssertEqual(t, gotTexts, wantTexts)
}

func TestNoteStrings(t *testing.T) {
	utils.Location = utils.UTCLocation()
	comments := model.Comments{&model.Comment{Text: "c1:\n- line 1\n\n- line 2\n- line 3"}, &model.Comment{Text: "c2"}, &model.Comment{Text: "c3"}}
	note := &model.Note{Text: "dummy text", Comments: comments, Status: "pending", Summary: "summary heading:\n- line 1\n- line 2", TagIds: []int{1, 2}, CompleteBy: 1609669235}
	want := `[  |          Text:  dummy text
   |      Comments:
  |              :  [nil] c1:
  |                     - line 1
  |                     - line 2
  |                     - line 3
  |              :  [nil] c2
  |              :  [nil] c3
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
	utils.AssertEqual(t, note.Strings(), want)
}

func TestExternalText(t *testing.T) {
	utils.Location = utils.UTCLocation()
	comments := model.Comments{&model.Comment{Text: "c < 1"}, &model.Comment{Text: "c > 2"}, &model.Comment{Text: "c & 3"}}
	note := &model.Note{Text: "dummy text", Comments: comments, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669235}
	var tags model.Tags
	tags = append(tags, &model.Tag{Id: 0, Slug: "tag_0", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 1, Slug: "tag_1", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "tag_2", Group: "tag_group2"})
	reminderData := &model.ReminderData{Tags: tags}
	want := `Note Details: -------------------------------------------------
  |          Text:  dummy text
  |      Comments:
  |              :  [nil] c &lt; 1
  |              :  [nil] c &gt; 2
  |              :  [nil] c &amp; 3
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
	utils.AssertEqual(t, note.ExternalText(reminderData), want)
}

func TestSearchableText(t *testing.T) {
	// case 1
	comments := model.Comments{&model.Comment{Text: "c1"}}
	note := model.Note{Text: "a beautiful cat", Comments: comments, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231}
	got := note.SearchableText()
	utils.AssertEqual(t, got, "| incidental | pending | ├ a beautiful cat ┤  [[nil] c1]")
	// case 2
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}}
	note = model.Note{Text: "a cute dog", Comments: comments, Status: "done", TagIds: []int{1, 2}, CompleteBy: 1609669232}
	got = note.SearchableText()
	utils.AssertEqual(t, got, "| incidental | done    | ├ a cute dog ┤  [[nil] c1, [nil] foo bar, [nil] c3]")
	// case 3
	comments = model.Comments{}
	note = model.Note{Text: "a cute dog", Comments: comments}
	got = note.SearchableText()
	utils.AssertEqual(t, got, "| incidental |         | ├ a cute dog ┤  [no-comments]")
	// case 4
	comments = model.Comments{}
	note = model.Note{Text: "first line\nsecondline\nthird line", Comments: comments}
	got = note.SearchableText()
	utils.AssertEqual(t, got, "| incidental |         | ├ first line NWL secondline NWL third line ┤  [no-comments]")
}

func TestExternalTexts(t *testing.T) {
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, "[]", "[]")
	// add notes
	comments := model.Comments{&model.Comment{Text: "c1"}}
	notes = append(notes, &model.Note{Text: "beautiful little cat", Comments: comments, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231})
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}, &model.Comment{Text: "baz"}}
	notes = append(notes, &model.Note{Text: "cute brown dog", Comments: comments, Status: "done", TagIds: []int{1, 2}, CompleteBy: 1609669232})
	// case 2
	got := notes.ExternalTexts(0, 0, 0)
	want := "[beautiful little cat {R: -, C:01, S:P, D:03-Jan-21} cute brown dog {R: -, C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 3
	got = notes.ExternalTexts(5, 0, 0)
	want = "[be... {R: -, C:01, S:P, D:03-Jan-21} cu... {R: -, C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 4
	got = notes.ExternalTexts(15, 0, 0)
	want = "[beautiful li... {R: -, C:01, S:P, D:03-Jan-21} cute brown dog  {R: -, C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 5
	got = notes.ExternalTexts(25, 0, 0)
	want = "[beautiful little cat      {R: -, C:01, S:P, D:03-Jan-21} cute brown dog            {R: -, C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
}

func TestWithStatus(t *testing.T) {
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, notes.WithStatus("pending"), model.Notes{})
	// add some notes
	comments := model.Comments{&model.Comment{Text: "c1"}}
	note1 := model.Note{Text: "big fat cat", Comments: comments, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231}
	notes = append(notes, &note1)
	comments = model.Comments{&model.Comment{Text: "c1"}, &model.Comment{Text: "foo bar"}}
	note2 := model.Note{Text: "cute brown dog", Comments: comments, Status: "done", TagIds: []int{1, 3}, CompleteBy: 1609669232}
	notes = append(notes, &note2)
	comments = model.Comments{&model.Comment{Text: "foo bar"}, &model.Comment{Text: "c3"}}
	note3 := model.Note{Text: "little hamster", Comments: comments, Status: "pending", TagIds: []int{1}, CompleteBy: 1609669233}
	notes = append(notes, &note3)
	// case 2 (with an invalid status)
	utils.AssertEqual(t, notes.WithStatus("no-such-status"), model.Notes{})
	// case 3 (with valid status "pending")
	got := notes.WithStatus("pending")
	want := model.Notes{&note1, &note3}
	utils.AssertEqual(t, got, want)
	// case 4 (with valid status "done")
	got = notes.WithStatus("done")
	want = model.Notes{&note2}
	utils.AssertEqual(t, got, want)
}

func TestWithTagIdAndStatus(t *testing.T) {
	// var tags model.Tags
	var notes model.Notes
	// case 1 (no notes)
	utils.AssertEqual(t, notes.WithTagIdAndStatus(2, "pending"), model.Notes{})
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
	note1 := model.Note{Text: "1", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	notes = append(notes, &note1)
	note2 := model.Note{Text: "2", Status: "pending", TagIds: []int{2, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}}
	notes = append(notes, &note2)
	note3 := model.Note{Text: "3", Status: "done", TagIds: []int{2}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}}
	notes = append(notes, &note3)
	note4 := model.Note{Text: "4", Status: "done", TagIds: []int{}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}}
	notes = append(notes, &note4)
	note5 := model.Note{Text: "5", Status: "pending", BaseStruct: model.BaseStruct{UpdatedAt: 1600000005}}
	notes = append(notes, &note5)
	// case 2
	utils.AssertEqual(t, notes.WithTagIdAndStatus(2, "pending"), []*model.Note{&note2})
	// case 3
	utils.AssertEqual(t, notes.WithTagIdAndStatus(2, "done"), []*model.Note{&note3})
	// case 4
	utils.AssertEqual(t, notes.WithTagIdAndStatus(4, "pending"), []*model.Note{&note1, &note2})
	// case 5
	utils.AssertEqual(t, notes.WithTagIdAndStatus(1, "done"), []*model.Note{})
}

func TestAddComment(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "1", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
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

func TestUpdateText(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
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

func TestUpdateSummary(t *testing.T) {
	// create notes
	note1 := model.Note{Summary: "original summary", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
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

func TestUpdateCompleteBy(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	utils.AssertEqual(t, note1.CompleteBy, 0)
	// update complete_by
	// case 1
	err := note1.UpdateCompleteBy("15-12-2021")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.CompleteBy, 1639526400)
	// case 2
	err = note1.UpdateCompleteBy("nil")
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.CompleteBy, 0)
}

func TestUpdateTags(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
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

func TestUpdateStatus(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// update TagIds
	// case 1
	err := note1.UpdateStatus("done", []int{1, 2, 3})
	utils.AssertEqual(t, err, errors.New("Note is part of a \"repeat\" group"))
	utils.AssertEqual(t, note1.Status, "pending")
	// case 2
	err = note1.UpdateStatus("done", []int{5, 6, 7})
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.Status, "done")
	// case 3
	err = note1.UpdateStatus("pending", []int{5, 6, 7})
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, note1.Status, "pending")
}

func TestRepeatType(t *testing.T) {
	repeatAnnuallyTagId := 3
	repeatMonthlyTagId := 4
	// create notes
	note1 := model.Note{Text: "original text1", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	note2 := model.Note{Text: "original text2", Status: "done", TagIds: []int{3, 5}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	note3 := model.Note{Text: "original text3", Status: "done", TagIds: []int{2, 6}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	// assert repeat type
	utils.AssertEqual(t, note1.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), "M")
	utils.AssertEqual(t, note2.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), "A")
	utils.AssertEqual(t, note3.RepeatType(repeatAnnuallyTagId, repeatMonthlyTagId), "-")
	utils.AssertEqual(t, note3.RepeatType(0, 0), "-")
}

func TestToggleMainFlag(t *testing.T) {
	// create notes
	note1 := model.Note{Text: "original text", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
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

func TestMakeSureFileExists(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))

	// make sure file doesn't exists already
	_, err := os.Stat(dataFilePath)
	utils.AssertEqual(t, err != nil, true)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), true)
	// attempt to create the file and required dirs, when the file doesn't exist already
	_ = model.MakeSureFileExists(dataFilePath)
	// prove that the file was created
	stats, err := os.Stat(dataFilePath)
	utils.AssertEqual(t, err != nil, false)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), false)

	// make sure that the existing file is not replaced
	modificationTime := stats.ModTime()
	// attempt to create the file and required dirs, when the file does exist already
	time.Sleep(10 * time.Millisecond)
	_ = model.MakeSureFileExists(dataFilePath)
	utils.AssertEqual(t, err != nil, false)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), false)
	stats, _ = os.Stat(dataFilePath)
	newModificationTime := stats.ModTime()
	utils.AssertEqual(t, newModificationTime == modificationTime, true)

}

func TestReadDataFile(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath)
	// attempt to read file and parse it
	reminderData := model.ReadDataFile(dataFilePath)
	utils.AssertEqual(t, reminderData.UpdatedAt > 0, true)
}

func TestUpdateDataFile(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath)
	reminderData := model.ReadDataFile(dataFilePath)
	// old_updated_at := reminderData.UpdatedAt
	testUser := model.User{Name: "Test User", EmailId: "user@test.com"}
	reminderData.User = &testUser
	_ = reminderData.UpdateDataFile("")
	remiderDataRe := model.ReadDataFile(dataFilePath)
	// utils.AssertEqual(t, remiderDataRe.UpdatedAt > old_updated_at, true)
	utils.AssertEqual(t, remiderDataRe.User.EmailId == testUser.EmailId, true)
}

func TestSortedTagsSlug(t *testing.T) {
	reminderData := model.ReminderData{
		User:  &model.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*model.Note{},
		Tags:  model.Tags{},
	}
	// creating tags
	var tags model.Tags
	// case 1 (no tags)
	utils.AssertEqual(t, reminderData.SortedTagSlugs(), []string{})
	// case 2 (has couple of tags)
	tags = append(tags, &model.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	reminderData.Tags = tags
	gotSlugs := reminderData.SortedTagSlugs()
	wantSlugs := []string{"a", "b", "c", "z"}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
}

func TestTagsFromIds(t *testing.T) {
	reminderData := model.ReminderData{
		User:  &model.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*model.Note{},
		Tags:  model.Tags{},
	}
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "z", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "c", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// case 1
	tagIDs := []int{1, 3}
	gotSlugs := reminderData.TagsFromIds(tagIDs)
	wantSlugs := model.Tags{&tag1, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 2
	tagIDs = []int{}
	gotSlugs = reminderData.TagsFromIds(tagIDs)
	wantSlugs = model.Tags{}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 3
	tagIDs = []int{1, 4, 2, 3}
	gotSlugs = reminderData.TagsFromIds(tagIDs)
	wantSlugs = model.Tags{&tag1, &tag4, &tag2, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
}

func TestTagFromSlug(t *testing.T) {
	reminderData := model.ReminderData{
		User:  &model.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*model.Note{},
		Tags:  model.Tags{},
	}
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// case 1
	utils.AssertEqual(t, reminderData.TagFromSlug("a"), &tag1)
	// case 2
	utils.AssertEqual(t, reminderData.TagFromSlug("a1"), &tag2)
	// case 3
	utils.AssertEqual(t, reminderData.TagFromSlug("no_slug"), nil)
}

func TestTagIdsForGroup(t *testing.T) {
	reminderData := model.ReminderData{
		User:  &model.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*model.Note{},
		Tags:  model.Tags{},
	}
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// case 1
	utils.AssertEqual(t, reminderData.TagIdsForGroup("tag_group1"), []int{1, 2, 3})
	// case 2
	utils.AssertEqual(t, reminderData.TagIdsForGroup("tag_group2"), []int{4})
	// case 3
	utils.AssertEqual(t, reminderData.TagIdsForGroup("tag_group_NO"), []int{})
}

/*
func TestNextPossibleTagId(t *testing.T) {
	reminderData := model.ReminderData{
		User:  &model.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*model.Note{},
		Tags:  model.Tags{},
	}
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	utils.AssertEqual(t, reminderData.nextPossibleTagId(), 4)
}
*/

func TestFindNotesByTagId(t *testing.T) {
	reminderData := model.ReminderData{
		User:  &model.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*model.Note{},
		Tags:  model.Tags{},
	}
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// create notes
	var notes model.Notes
	note1 := model.Note{Text: "1", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	notes = append(notes, &note1)
	note2 := model.Note{Text: "2", Status: "pending", TagIds: []int{2, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}}
	notes = append(notes, &note2)
	note3 := model.Note{Text: "3", Status: "done", TagIds: []int{2}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}}
	notes = append(notes, &note3)
	note4 := model.Note{Text: "4", Status: "done", TagIds: []int{}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}}
	notes = append(notes, &note4)
	note5 := model.Note{Text: "5", Status: "pending", BaseStruct: model.BaseStruct{UpdatedAt: 1600000005}}
	notes = append(notes, &note5)
	reminderData.Notes = notes
	// searching notes
	// case 1
	utils.AssertEqual(t, reminderData.FindNotesByTagId(2, "pending"), []*model.Note{&note2})
	// case 2
	utils.AssertEqual(t, reminderData.FindNotesByTagId(2, "done"), []*model.Note{&note3})
	// case 3
	utils.AssertEqual(t, reminderData.FindNotesByTagId(4, "pending"), []*model.Note{&note1, &note2})
	// case 4
	utils.AssertEqual(t, reminderData.FindNotesByTagId(1, "done"), []*model.Note{})
}

func TestFindNotesByTagSlug(t *testing.T) {
	reminderData := model.ReminderData{
		User:  &model.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*model.Note{},
		Tags:  model.Tags{},
	}
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// create notes
	var notes model.Notes
	note1 := model.Note{Text: "1", Status: "pending", TagIds: []int{1, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}}
	notes = append(notes, &note1)
	note2 := model.Note{Text: "2", Status: "pending", TagIds: []int{2, 4}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}}
	notes = append(notes, &note2)
	note3 := model.Note{Text: "3", Status: "done", TagIds: []int{2}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}}
	notes = append(notes, &note3)
	note4 := model.Note{Text: "4", Status: "done", TagIds: []int{}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}}
	notes = append(notes, &note4)
	note5 := model.Note{Text: "5", Status: "pending", BaseStruct: model.BaseStruct{UpdatedAt: 1600000005}}
	notes = append(notes, &note5)
	reminderData.Notes = notes
	// searching notes
	// case 1
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("a1", "pending"), []*model.Note{&note2})
	// case 2
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("a1", "done"), []*model.Note{&note3})
	// case 3
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("b", "pending"), []*model.Note{&note1, &note2})
	// case 4
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("a", "done"), []*model.Note{})
}

func TestRegisterBasicTags(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath)
	reminderData := model.ReadDataFile(dataFilePath)
	// register basic tags
	_ = reminderData.RegisterBasicTags()
	utils.AssertEqual(t, len(reminderData.Tags), 7)
}

func TestNotesApprachingDueDate(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath)
	reminderData := model.ReadDataFile(dataFilePath)
	// register basic tags
	_ = reminderData.RegisterBasicTags()
	// get current time
	currentTime := utils.CurrentUnixTimestamp()
	// register notes
	// for reference: here is the list of basic tags
	// {"slug": "current", "group": ""},
	// {"slug": "priority-urgent", "group": "priority"},
	// {"slug": "priority-medium", "group": "priority"},
	// {"slug": "priority-low", "group": "priority"},
	// {"slug": "repeat-annually", "group": "repeat"},
	// {"slug": "repeat-monthly", "group": "repeat"},
	// {"slug": "tips", "group": "tips"}}
	currentTagId := reminderData.TagFromSlug("current").Id
	repeatAnnuallyTagId := reminderData.TagFromSlug("repeat-annually").Id
	repeatMonthlyTagId := reminderData.TagFromSlug("repeat-monthly").Id
	var notes model.Notes
	// non-repeat done notes
	notes = append(notes, &model.Note{Text: "NRD01a", Status: "done", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 8*24*3600})
	notes = append(notes, &model.Note{Text: "NRD02a", Status: "done", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})
	notes = append(notes, &model.Note{Text: "NRD03a", Status: "done", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600})
	// repeat-annually done notes
	notes = append(notes, &model.Note{Text: "RAD01a", Status: "done", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 8*24*3600})
	notes = append(notes, &model.Note{Text: "RAD02a", Status: "done", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})
	notes = append(notes, &model.Note{Text: "RAD03a", Status: "done", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600})
	// repeat-monthally done notes
	notes = append(notes, &model.Note{Text: "RMD01a", Status: "done", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 8*24*3600})
	notes = append(notes, &model.Note{Text: "RMD02a", Status: "done", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})
	notes = append(notes, &model.Note{Text: "RMD03a", Status: "done", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600})
	// non-repeat pending notes
	notes = append(notes, &model.Note{Text: "NRP01a", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600})        // expected
	notes = append(notes, &model.Note{Text: "NRP02a", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 2*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP02b", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 2*24*3600 + 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP03a", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})        // expected
	notes = append(notes, &model.Note{Text: "NRP04a", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP04b", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 + 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP05a", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 3*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP05b", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 3*24*3600 + 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP06a", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 7*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP06b", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 7*24*3600 + 3600})
	notes = append(notes, &model.Note{Text: "NRP07a", Status: "pending", TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600})
	// repeat-annually pending notes
	notes = append(notes, &model.Note{Text: "RAP01", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600})
	notes = append(notes, &model.Note{Text: "RAP02", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RAP03", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime})             // expected
	notes = append(notes, &model.Note{Text: "RAP04", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RAP05", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RAP06", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600})
	notes = append(notes, &model.Note{Text: "RAP07", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600 - 2*365})
	notes = append(notes, &model.Note{Text: "RAP08", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600 - 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP09", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 2*365})             // expected
	notes = append(notes, &model.Note{Text: "RAP10", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600 - 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP11", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 - 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP12", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600 - 2*365})
	notes = append(notes, &model.Note{Text: "RAP13", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600 + 2*365})
	notes = append(notes, &model.Note{Text: "RAP14", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600 + 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP15", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*365})             // expected
	notes = append(notes, &model.Note{Text: "RAP16", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600 + 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP17", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 + 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP18", Status: "pending", TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600 + 2*365})
	// repeat-monthly pending notes
	notes = append(notes, &model.Note{Text: "RMP01", Status: "pending", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600})
	notes = append(notes, &model.Note{Text: "RMP02", Status: "pending", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RMP03", Status: "pending", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime})             // expected
	notes = append(notes, &model.Note{Text: "RMP04", Status: "pending", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RMP05", Status: "pending", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 3*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RMP06", Status: "pending", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 7*24*3600})
	notes = append(notes, &model.Note{Text: "RMP07", Status: "pending", TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600})
	reminderData.Notes = notes
	// get urgent notes
	urgentNotes := reminderData.NotesApprachingDueDate("default")
	var urgentNotesText []string
	for _, note := range urgentNotes {
		urgentNotesText = append(urgentNotesText, note.Text)
	}
	utils.AssertEqual(t, urgentNotesText, []string{
		"NRP01a", "NRP02a", "NRP02b", "NRP03a", "NRP04a", "NRP04b", "NRP05a", "NRP05b", "NRP06a",
		"RAP02", "RAP03", "RAP04", "RAP05", "RAP08", "RAP09", "RAP10", "RAP11", "RAP14", "RAP15", "RAP16", "RAP17",
		"RMP02", "RMP03",
	})
}

func TestPrintStats(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath)
	reminderData := model.ReadDataFile(dataFilePath)
	// register basic tags
	_ = reminderData.RegisterBasicTags()
	got := reminderData.Stats()
	want := `
Stats of "temp_test_dir/mydata.json"
  - Number of Tags: 7
  - Pending Notes: 0/0
`
	utils.AssertEqual(t, got, want)
}

func TestNewTagRegistration(t *testing.T) {
	dataFilePath := path.Join("..", "..", "test", "test_data_file.json")
	reminderData := model.ReadDataFile(dataFilePath)
	utils.AssertEqual(t, len(reminderData.Tags), 5)
}

func TestNewTag(t *testing.T) {
	dummySlug := "test_tag_slug"
	dummyGroup := "test_tag_group"
	tag, _ := model.NewTag(10, dummySlug, dummyGroup)
	want := &model.Tag{
		Id:    10,
		Slug:  dummySlug,
		Group: dummyGroup,
	}
	utils.AssertEqual(t, tag, want)
}
func TestNewNote(t *testing.T) {
	tagIDs := []int{1, 3, 5}
	dummyText := "a random note text"
	note, _ := model.NewNote(tagIDs, dummyText)
	want := &model.Note{
		Text:       dummyText,
		TagIds:     tagIDs,
		Status:     note.Status,
		BaseStruct: model.BaseStruct{UpdatedAt: note.UpdatedAt, CreatedAt: note.CreatedAt},
	}
	utils.AssertEqual(t, note, want)
}
