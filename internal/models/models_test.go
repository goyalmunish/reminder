package models_test

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

	models "reminder/internal/models"
	utils "reminder/pkg/utils"
)

func TestDataFile(t *testing.T) {
	defaultDataFilePath := models.FDefaultDataFile()
	utils.AssertEqual(t, strings.HasPrefix(defaultDataFilePath, "/"), true)
	utils.AssertEqual(t, strings.HasSuffix(defaultDataFilePath, ".json"), true)
}

func TestUser(t *testing.T) {
	got := models.User{Name: "Test User", EmailId: "user@test.com"}
	want := "{Name: Test User, EmailId: user@test.com}"
	utils.AssertEqual(t, got, want)
}

func TestTag(t *testing.T) {
	// case 1: general case
	got := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	want := "tag_group1#a#1"
	utils.AssertEqual(t, got, want)
	// case 2: blank group
	got = models.Tag{Id: 1, Slug: "a", Group: ""}
	want = "#a#1"
	utils.AssertEqual(t, got, want)
	// case 3: omitted group
	got = models.Tag{Id: 1, Slug: "a"}
	want = "#a#1"
	utils.AssertEqual(t, got, want)
}

func TestTags(t *testing.T) {
	var tags models.Tags
	tags = append(tags, &models.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	sort.Sort(models.Tags(tags))
	var got []int
	for _, value := range tags {
		got = append(got, value.Id)
	}
	want := []int{1, 4, 3, 2}
	utils.AssertEqual(t, got, want)
}

func TestSlugs(t *testing.T) {
	var tags models.Tags
	utils.AssertEqual(t, tags, "[]")
	tags = append(tags, &models.Tag{Id: 1, Slug: "tag_1", Group: "tag_group"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "tag_2", Group: "tag_group"})
	tags = append(tags, &models.Tag{Id: 3, Slug: "tag_3", Group: "tag_group"})
	got := tags.Slugs()
	want := "[tag_1 tag_2 tag_3]"
	utils.AssertEqual(t, got, want)
}

func TestFromSlug(t *testing.T) {
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1
	utils.AssertEqual(t, tags.FromSlug("a"), &tag1)
	// case 2
	utils.AssertEqual(t, tags.FromSlug("a1"), &tag2)
	// case 3
	utils.AssertEqual(t, tags.FromSlug("no_slug"), nil)
}

func TestFromIds(t *testing.T) {
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "z", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "c", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1
	tagIDs := []int{1, 3}
	gotSlugs := tags.FromIds(tagIDs)
	wantSlugs := models.Tags{&tag1, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 2
	tagIDs = []int{}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = models.Tags{}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 3
	tagIDs = []int{1, 4, 2, 3}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = models.Tags{&tag1, &tag4, &tag2, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
}

func TestIdsForGroup(t *testing.T) {
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1
	utils.AssertEqual(t, tags.IdsForGroup("tag_group1"), []int{1, 2, 3})
	// case 2
	utils.AssertEqual(t, tags.IdsForGroup("tag_group2"), []int{4})
	// case 3
	utils.AssertEqual(t, tags.IdsForGroup("tag_group_NO"), []int{})
}

func TestFBasicTags(t *testing.T) {
	basicTags := models.FBasicTags()
	slugs := basicTags.Slugs()
	want := "[current priority-urgent priority-medium priority-low repeat-annually repeat-monthly tips]"
	utils.AssertEqual(t, slugs, want)
}

func TestNotes(t *testing.T) {
	var notes []*models.Note
	notes = append(notes, &models.Note{Text: "1", Status: "pending", UpdatedAt: 1600000001})
	notes = append(notes, &models.Note{Text: "2", Status: "pending", UpdatedAt: 1600000004})
	notes = append(notes, &models.Note{Text: "3", Status: "done", UpdatedAt: 1600000003})
	notes = append(notes, &models.Note{Text: "4", Status: "done", UpdatedAt: 1600000002})
	sort.Sort(models.Notes(notes))
	var gotTexts []string
	for _, value := range notes {
		gotTexts = append(gotTexts, value.Text)
	}
	wantTexts := []string{"2", "3", "4", "1"}
	utils.AssertEqual(t, gotTexts, wantTexts)
}

func TestNoteString(t *testing.T) {
	note := &models.Note{Text: "dummy text", Comments: []string{"c1", "c2", "c3"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669235}
	want := `[  |          Text:  dummy text
   |              :  c1
  |              :  c2
  |              :  c3
   |        Status:  pending
   |          Tags:  [1 2]
   |    CompleteBy:  Sunday, 03-Jan-21 10:20:35 UTC
   |     CreatedAt:  nil
   |     UpdatedAt:  nil
]`
	utils.AssertEqual(t, note.String(), want)
}

func TestExternalText(t *testing.T) {
	note := &models.Note{Text: "dummy text", Comments: []string{"c1", "c2", "c3"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669235}
	var tags models.Tags
	tags = append(tags, &models.Tag{Id: 0, Slug: "tag_0", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 1, Slug: "tag_1", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "tag_2", Group: "tag_group2"})
	reminderData := &models.ReminderData{Tags: tags}
	want := `Note Details: -------------------------------------------------
  |          Text:  dummy text
  |              :  c1
  |              :  c2
  |              :  c3
  |        Status:  pending
  |              :  tag_1
  |              :  tag_2
  |    CompleteBy:  Sunday, 03-Jan-21 10:20:35 UTC
  |     CreatedAt:  nil
  |     UpdatedAt:  nil
`
	utils.AssertEqual(t, note.ExternalText(reminderData), want)
}

func TestSearchableText(t *testing.T) {
	// case 1
	note := models.Note{Text: "a beautiful cat", Comments: []string{"c1"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231}
	got := note.SearchableText()
	utils.AssertEqual(t, got, "a beautiful cat [c1]")
	// case 2
	note = models.Note{Text: "a cute dog", Comments: []string{"c1", "foo bar", "c3"}, Status: "done", TagIds: []int{1, 2}, CompleteBy: 1609669232}
	got = note.SearchableText()
	utils.AssertEqual(t, got, "a cute dog [c1, foo bar, c3]")
	// case 3
	note = models.Note{Text: "a cute dog", Comments: []string{}}
	got = note.SearchableText()
	utils.AssertEqual(t, got, "a cute dog [no-comments]")
}

func TestExternalTexts(t *testing.T) {
	var notes models.Notes
	notes = append(notes, &models.Note{Text: "beautiful little cat", Comments: []string{"c1"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231})
	notes = append(notes, &models.Note{Text: "cute brown dog", Comments: []string{"c1", "foo bar", "c3", "baz"}, Status: "done", TagIds: []int{1, 2}, CompleteBy: 1609669232})
	// case 1
	got := notes.ExternalTexts(0)
	want := "[beautiful little cat {C:01, S:P, D:03-Jan-21} cute brown dog {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 2
	got = notes.ExternalTexts(5)
	want = "[be... {C:01, S:P, D:03-Jan-21} cu... {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 3
	got = notes.ExternalTexts(15)
	want = "[beautiful li... {C:01, S:P, D:03-Jan-21} cute brown dog  {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 4
	got = notes.ExternalTexts(25)
	want = "[beautiful little cat      {C:01, S:P, D:03-Jan-21} cute brown dog            {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
}

func TestWithStatus(t *testing.T) {
	var notes models.Notes
	note1 := models.Note{Text: "big fat cat", Comments: []string{"c1"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231}
	notes = append(notes, &note1)
	note2 := models.Note{Text: "cute brown dog", Comments: []string{"c1", "foo bar"}, Status: "done", TagIds: []int{1, 3}, CompleteBy: 1609669232}
	notes = append(notes, &note2)
	note3 := models.Note{Text: "little hamster", Comments: []string{"foo bar", "c3"}, Status: "pending", TagIds: []int{1}, CompleteBy: 1609669233}
	notes = append(notes, &note3)
	// case 1
	got := notes.WithStatus("pending")
	want := []*models.Note{&note1, &note3}
	utils.AssertEqual(t, got, want)
	// case 2
	got = notes.WithStatus("done")
	want = []*models.Note{&note2}
	utils.AssertEqual(t, got, want)
}

func TestFMakeSureFileExists(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))

	// make sure file doesn't exists already
	_, err := os.Stat(dataFilePath)
	utils.AssertEqual(t, err != nil, true)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), true)
	// attempt to create the file and required dirs, when the file doesn't exist already
	models.FMakeSureFileExists(dataFilePath)
	// prove that the file was created
	stats, err := os.Stat(dataFilePath)
	utils.AssertEqual(t, err != nil, false)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), false)

	// make sure that the existing file is not replaced
	modificationTime := stats.ModTime()
	// attempt to create the file and required dirs, when the file does exist already
	time.Sleep(10 * time.Millisecond)
	models.FMakeSureFileExists(dataFilePath)
	utils.AssertEqual(t, err != nil, false)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), false)
	stats, err = os.Stat(dataFilePath)
	newModificationTime := stats.ModTime()
	utils.AssertEqual(t, newModificationTime == modificationTime, true)

}

func TestFReadDataFile(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	models.FMakeSureFileExists(dataFilePath)
	// attempt to read file and parse it
	reminderData := models.FReadDataFile(dataFilePath)
	utils.AssertEqual(t, reminderData.UpdatedAt > 0, true)
}

func TestUpdateDataFile(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	models.FMakeSureFileExists(dataFilePath)
	reminderData := models.FReadDataFile(dataFilePath)
	// old_updated_at := reminderData.UpdatedAt
	testUser := models.User{Name: "Test User", EmailId: "user@test.com"}
	reminderData.User = &testUser
	reminderData.UpdateDataFile()
	remiderDataRe := models.FReadDataFile(dataFilePath)
	// utils.AssertEqual(t, remiderDataRe.UpdatedAt > old_updated_at, true)
	utils.AssertEqual(t, remiderDataRe.User.EmailId == testUser.EmailId, true)
}

func TestTagsSlug(t *testing.T) {
	reminderData := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  models.Tags{},
	}
	// creating tags
	var tags models.Tags
	tags = append(tags, &models.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	reminderData.Tags = tags
	gotSlugs := reminderData.TagsSlugs()
	wantSlugs := []string{"a", "b", "c", "z"}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
}

func TestTagsFromIds(t *testing.T) {
	reminderData := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  models.Tags{},
	}
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "z", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "c", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// case 1
	tagIDs := []int{1, 3}
	gotSlugs := reminderData.TagsFromIds(tagIDs)
	wantSlugs := models.Tags{&tag1, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 2
	tagIDs = []int{}
	gotSlugs = reminderData.TagsFromIds(tagIDs)
	wantSlugs = models.Tags{}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 3
	tagIDs = []int{1, 4, 2, 3}
	gotSlugs = reminderData.TagsFromIds(tagIDs)
	wantSlugs = models.Tags{&tag1, &tag4, &tag2, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
}

func TestTagFromSlug(t *testing.T) {
	reminderData := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  models.Tags{},
	}
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
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
	reminderData := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  models.Tags{},
	}
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// case 1
	utils.AssertEqual(t, reminderData.TagIdsForGroup("tag_group1"), []int{1, 2, 3})
	// case 2
	utils.AssertEqual(t, reminderData.TagIdsForGroup("tag_group2"), []int{4})
	// case 3
	utils.AssertEqual(t, reminderData.TagIdsForGroup("tag_group_NO"), []int{})
}

func TestNextPossibleTagId(t *testing.T) {
	reminderData := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  models.Tags{},
	}
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	utils.AssertEqual(t, reminderData.NextPossibleTagId(), 4)
}

func TestNotesWithTagId(t *testing.T) {
	reminderData := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  models.Tags{},
	}
	// creating tags
	var tags models.Tags
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminderData.Tags = tags
	// create notes
	var notes []*models.Note
	note1 := models.Note{Text: "1", Status: "pending", TagIds: []int{1, 4}, UpdatedAt: 1600000001}
	notes = append(notes, &note1)
	note2 := models.Note{Text: "2", Status: "pending", TagIds: []int{2, 4}, UpdatedAt: 1600000004}
	notes = append(notes, &note2)
	note3 := models.Note{Text: "3", Status: "done", TagIds: []int{2}, UpdatedAt: 1600000003}
	notes = append(notes, &note3)
	note4 := models.Note{Text: "4", Status: "done", TagIds: []int{}, UpdatedAt: 1600000002}
	notes = append(notes, &note4)
	note5 := models.Note{Text: "5", Status: "pending", UpdatedAt: 1600000005}
	notes = append(notes, &note5)
	reminderData.Notes = notes
	// searching notes
	// case 1
	utils.AssertEqual(t, reminderData.NotesWithTagId(2, "pending"), []*models.Note{&note2})
	// case 2
	utils.AssertEqual(t, reminderData.NotesWithTagId(2, "done"), []*models.Note{&note3})
	// case 3
	utils.AssertEqual(t, reminderData.NotesWithTagId(4, "pending"), []*models.Note{&note1, &note2})
	// case 4
	utils.AssertEqual(t, reminderData.NotesWithTagId(1, "done"), []*models.Note{})
}

func TestRegisterBasicTags(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	models.FMakeSureFileExists(dataFilePath)
	reminderData := models.FReadDataFile(dataFilePath)
	// register basic tags
	reminderData.RegisterBasicTags()
	utils.AssertEqual(t, len(reminderData.Tags), 7)
}

func TestPrintStats(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	models.FMakeSureFileExists(dataFilePath)
	reminderData := models.FReadDataFile(dataFilePath)
	// register basic tags
	reminderData.RegisterBasicTags()
	got := reminderData.Stats()
	want := `
Stats of "temp_test_dir/mydata.json"
  - Number of Tags: 7
  - Pending Notes: 0/0
`
	utils.AssertEqual(t, got, want)
}
