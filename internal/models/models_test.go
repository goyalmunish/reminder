package models_test

import (
	// "fmt"
	// "errors"
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

func TestTagsBySlug(t *testing.T) {
	var tags []*models.Tag
	tags = append(tags, &models.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &models.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	sort.Sort(models.TagsBySlug(tags))
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

func TestNotesByUpdatedAt(t *testing.T) {
	var notes []*models.Note
	notes = append(notes, &models.Note{Text: "1", Status: "pending", UpdatedAt: 1600000001})
	notes = append(notes, &models.Note{Text: "2", Status: "pending", UpdatedAt: 1600000004})
	notes = append(notes, &models.Note{Text: "3", Status: "done", UpdatedAt: 1600000003})
	notes = append(notes, &models.Note{Text: "4", Status: "done", UpdatedAt: 1600000002})
	sort.Sort(models.NotesByUpdatedAt(notes))
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
