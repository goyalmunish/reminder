package models_test

import (
	// "fmt"
	"errors"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
	// "time"
	// "github.com/golang/mock/gomock"

	models "reminder/internal/models"
	utils "reminder/pkg/utils"
)

func TestDataFile(t *testing.T) {
	utils.AssertEqual(t, strings.HasPrefix(models.DataFile, "/"), true)
	utils.AssertEqual(t, strings.HasSuffix(models.DataFile, ".json"), true)
}

func TestUser(t *testing.T) {
	got := models.User{Name: "Test User", EmailId: "user@test.com"}
	want := "{Name: Test User, EmailId: user@test.com}"
	utils.AssertEqual(t, got, want)
}

func TestFTagsBySlug(t *testing.T) {
	var tags []*models.Tag
	tags = append(tags, &models.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	sort.Sort(models.FTagsBySlug(tags))
	var got_ids []int
	for _, value := range tags {
		got_ids = append(got_ids, value.Id)
	}
	want_ids := []int{1, 4, 3, 2}
	utils.AssertEqual(t, got_ids, want_ids)
}

func TestFTagsSlugs(t *testing.T) {
	var tags []*models.Tag
	utils.AssertEqual(t, tags, "[]")
	tags = append(tags, &models.Tag{Id: 1, Slug: "tag_1", Group: "tag_group"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "tag_2", Group: "tag_group"})
	tags = append(tags, &models.Tag{Id: 3, Slug: "tag_3", Group: "tag_group"})
	got := models.FTagsSlugs(tags)
	want := "[tag_1 tag_2 tag_3]"
	utils.AssertEqual(t, got, want)
}

func TestFBasicTags(t *testing.T) {
	basic_tags := models.FBasicTags()
	slugs := models.FTagsSlugs(basic_tags)
	want := "[current priority-urgent priority-medium priority-low repeat-annually repeat-monthly tips]"
	utils.AssertEqual(t, slugs, want)
}

func TestFNotesByUpdatedAt(t *testing.T) {
	var notes []*models.Note
	notes = append(notes, &models.Note{Text: "1", Status: "pending", UpdatedAt: 1600000001})
	notes = append(notes, &models.Note{Text: "2", Status: "pending", UpdatedAt: 1600000004})
	notes = append(notes, &models.Note{Text: "3", Status: "done", UpdatedAt: 1600000003})
	notes = append(notes, &models.Note{Text: "4", Status: "done", UpdatedAt: 1600000002})
	sort.Sort(models.FNotesByUpdatedAt(notes))
	var got_texts []string
	for _, value := range notes {
		got_texts = append(got_texts, value.Text)
	}
	want_texts := []string{"2", "3", "4", "1"}
	utils.AssertEqual(t, got_texts, want_texts)
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

func TestStringRepr(t *testing.T) {
	note := &models.Note{Text: "dummy text", Comments: []string{"c1", "c2", "c3"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669235}
	var tags []*models.Tag
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
	utils.AssertEqual(t, note.StringRepr(reminderData), want)
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
}

func TestFNotesTexts(t *testing.T) {
	var notes []*models.Note
	notes = append(notes, &models.Note{Text: "beautiful little cat", Comments: []string{"c1"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231})
	notes = append(notes, &models.Note{Text: "cute brown dog", Comments: []string{"c1", "foo bar", "c3", "baz"}, Status: "done", TagIds: []int{1, 2}, CompleteBy: 1609669232})
	// case 1
	got := models.FNotesTexts(notes, 0)
	want := "[beautiful little cat {C:01, S:P, D:03-Jan-21} cute brown dog {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 2
	got = models.FNotesTexts(notes, 5)
	want = "[be... {C:01, S:P, D:03-Jan-21} cu... {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 3
	got = models.FNotesTexts(notes, 15)
	want = "[beautiful li... {C:01, S:P, D:03-Jan-21} cute brown dog  {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
	// case 4
	got = models.FNotesTexts(notes, 25)
	want = "[beautiful little cat      {C:01, S:P, D:03-Jan-21} cute brown dog            {C:04, S:D, D:03-Jan-21}]"
	utils.AssertEqual(t, got, want)
}

func TestFNotesWithStatus(t *testing.T) {
	var notes []*models.Note
	note1 := models.Note{Text: "big fat cat", Comments: []string{"c1"}, Status: "pending", TagIds: []int{1, 2}, CompleteBy: 1609669231}
	notes = append(notes, &note1)
	note2 := models.Note{Text: "cute brown dog", Comments: []string{"c1", "foo bar"}, Status: "done", TagIds: []int{1, 3}, CompleteBy: 1609669232}
	notes = append(notes, &note2)
	note3 := models.Note{Text: "little hamster", Comments: []string{"foo bar", "c3"}, Status: "pending", TagIds: []int{1}, CompleteBy: 1609669233}
	notes = append(notes, &note3)
	// case 1
	got := models.FNotesWithStatus(notes, "pending")
	want := []*models.Note{&note1, &note3}
	utils.AssertEqual(t, got, want)
	// case 2
	got = models.FNotesWithStatus(notes, "done")
	want = []*models.Note{&note2}
	utils.AssertEqual(t, got, want)
}

func TestFMakeSureFileExists(t *testing.T) {
	var data_file_path = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(data_file_path))
	// make sure file doesn't exists already
	_, err := os.Stat(data_file_path)
	utils.AssertEqual(t, err != nil, true)
	errors.Is(err, fs.ErrNotExist)
	// attempt to create the file and required dirs
	models.FMakeSureFileExists(data_file_path)
	// prove that the file was created
	_, err = os.Stat(data_file_path)
	utils.AssertEqual(t, err == nil, true)
}

func TestFReadDataFile(t *testing.T) {
	var data_file_path = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(data_file_path))
	// create the file and required dirs
	models.FMakeSureFileExists(data_file_path)
	// attempt to read file and parse it
	reminder_data := models.FReadDataFile(data_file_path)
	utils.AssertEqual(t, reminder_data.UpdatedAt > 0, true)
}

func TestUpdateDataFile(t *testing.T) {
	var data_file_path = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(data_file_path))
	// create the file and required dirs
	models.FMakeSureFileExists(data_file_path)
	reminder_data := models.FReadDataFile(data_file_path)
	// old_updated_at := reminder_data.UpdatedAt
	test_user := models.User{Name: "Test User", EmailId: "user@test.com"}
	reminder_data.User = &test_user
	reminder_data.UpdateDataFile(data_file_path)
	reminder_data_re := models.FReadDataFile(data_file_path)
	// utils.AssertEqual(t, reminder_data_re.UpdatedAt > old_updated_at, true)
	utils.AssertEqual(t, reminder_data_re.User.EmailId == test_user.EmailId, true)
}

func TestTagsSlug(t *testing.T) {
	reminder_data := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  []*models.Tag{},
	}
	// creating tags
	var tags []*models.Tag
	tags = append(tags, &models.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	reminder_data.Tags = tags
	got_slugs := reminder_data.TagsSlugs()
	want_slugs := []string{"a", "b", "c", "z"}
	utils.AssertEqual(t, got_slugs, want_slugs)
}

func TestTagsFromIds(t *testing.T) {
	reminder_data := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  []*models.Tag{},
	}
	// creating tags
	var tags []*models.Tag
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "z", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "c", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminder_data.Tags = tags
	// case 1
	tag_ids := []int{1, 3}
	got_slugs := reminder_data.TagsFromIds(tag_ids)
	want_slugs := []*models.Tag{&tag1, &tag3}
	utils.AssertEqual(t, got_slugs, want_slugs)
	// case 2
	tag_ids = []int{}
	got_slugs = reminder_data.TagsFromIds(tag_ids)
	want_slugs = []*models.Tag{}
	utils.AssertEqual(t, got_slugs, want_slugs)
	// case 3
	tag_ids = []int{1, 4, 2, 3}
	got_slugs = reminder_data.TagsFromIds(tag_ids)
	want_slugs = []*models.Tag{&tag1, &tag4, &tag2, &tag3}
	utils.AssertEqual(t, got_slugs, want_slugs)
}

func TestTagFromSlug(t *testing.T) {
	reminder_data := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  []*models.Tag{},
	}
	// creating tags
	var tags []*models.Tag
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminder_data.Tags = tags
	// case 1
	utils.AssertEqual(t, reminder_data.TagFromSlug("a"), &tag1)
	// case 2
	utils.AssertEqual(t, reminder_data.TagFromSlug("a1"), &tag2)
	// case 3
	utils.AssertEqual(t, reminder_data.TagFromSlug("no_slug"), nil)
}

func TestTagIdsForGroup(t *testing.T) {
	reminder_data := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  []*models.Tag{},
	}
	// creating tags
	var tags []*models.Tag
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminder_data.Tags = tags
	// case 1
	utils.AssertEqual(t, reminder_data.TagIdsForGroup("tag_group1"), []int{1, 2, 3})
	// case 2
	utils.AssertEqual(t, reminder_data.TagIdsForGroup("tag_group2"), []int{4})
	// case 3
	utils.AssertEqual(t, reminder_data.TagIdsForGroup("tag_group_NO"), []int{})
}

func TestNextPossibleTagId(t *testing.T) {
	reminder_data := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  []*models.Tag{},
	}
	// creating tags
	var tags []*models.Tag
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminder_data.Tags = tags
	utils.AssertEqual(t, reminder_data.NextPossibleTagId(), 4)
}

func TestNotesWithTagId(t *testing.T) {
	reminder_data := models.ReminderData{
		User:  &models.User{Name: "Test User", EmailId: "user@test.com"},
		Notes: []*models.Note{},
		Tags:  []*models.Tag{},
	}
	// creating tags
	var tags []*models.Tag
	tag1 := models.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := models.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := models.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := models.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	reminder_data.Tags = tags
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
	reminder_data.Notes = notes
	// searching notes
	// case 1
	utils.AssertEqual(t, reminder_data.NotesWithTagId(2, "pending"), []*models.Note{&note2})
	// case 2
	utils.AssertEqual(t, reminder_data.NotesWithTagId(2, "done"), []*models.Note{&note3})
	// case 3
	utils.AssertEqual(t, reminder_data.NotesWithTagId(4, "pending"), []*models.Note{&note1, &note2})
	// case 4
	utils.AssertEqual(t, reminder_data.NotesWithTagId(1, "done"), []*models.Note{})
}

func TestRegisterBasicTags(t *testing.T) {
	var data_file_path = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(data_file_path))
	// create the file and required dirs
	models.FMakeSureFileExists(data_file_path)
	reminder_data := models.FReadDataFile(data_file_path)
	// register basic tags
	reminder_data.RegisterBasicTags()
	utils.AssertEqual(t, len(reminder_data.Tags), 7)
}
