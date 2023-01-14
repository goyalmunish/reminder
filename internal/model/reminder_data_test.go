package model_test

import (
	"fmt"

	// "fmt"

	"os"
	"path"
	"sort"
	"strings"
	"testing"

	// "github.com/golang/mock/gomock"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/internal/settings"
	utils "github.com/goyalmunish/reminder/pkg/utils"
)

// mocks

type MockPromptTagSlug struct {
}
type MockPromptTagGroup struct {
}
type MockPromptNoteText struct {
}
type TestTagger struct{}

func (tagger TestTagger) TagsFromIds(tagIDs []int) []string {
	slugs := []string{}
	for i, id := range tagIDs {
		slug := fmt.Sprintf("%d-%d", i, id)
		slugs = append(slugs, slug)
	}
	return slugs
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
	defaultDataFilePath := settings.DefaultSettings().AppInfo.DataFile
	utils.AssertEqual(t, strings.HasPrefix(defaultDataFilePath, "/") || strings.HasPrefix(defaultDataFilePath, "~/"), true)
	utils.AssertEqual(t, strings.HasSuffix(defaultDataFilePath, ".json"), true)
}

func TestReminderDataTagFromSlug(t *testing.T) {
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

func TestReminderDataTagIdsForGroup(t *testing.T) {
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

func TestReminderDataTagsFromIds(t *testing.T) {
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
	wantSlugs := model.Tags{&tag1, &tag3}.Slugs()
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 2
	tagIDs = []int{}
	gotSlugs = reminderData.TagsFromIds(tagIDs)
	wantSlugs = model.Tags{}.Slugs()
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 3
	tagIDs = []int{1, 4, 2, 3}
	gotSlugs = reminderData.TagsFromIds(tagIDs)
	wantSlugs = model.Tags{&tag1, &tag4, &tag2, &tag3}.Slugs()
	utils.AssertEqual(t, gotSlugs, wantSlugs)
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
	reminderData.Notes = notes
	// searching notes
	// case 1
	utils.AssertEqual(t, reminderData.FindNotesByTagId(2, model.NoteStatus_Pending), []*model.Note{&note2})
	// case 2
	utils.AssertEqual(t, reminderData.FindNotesByTagId(2, model.NoteStatus_Done), []*model.Note{&note3})
	// case 3
	utils.AssertEqual(t, reminderData.FindNotesByTagId(4, model.NoteStatus_Pending), []*model.Note{&note1, &note2})
	// case 4
	utils.AssertEqual(t, reminderData.FindNotesByTagId(1, model.NoteStatus_Done), []*model.Note{})
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
	reminderData.Notes = notes
	// searching notes
	// case 1
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("a1", model.NoteStatus_Pending), []*model.Note{&note2})
	// case 2
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("a1", model.NoteStatus_Done), []*model.Note{&note3})
	// case 3
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("b", model.NoteStatus_Pending), []*model.Note{&note1, &note2})
	// case 4
	utils.AssertEqual(t, reminderData.FindNotesByTagSlug("a", model.NoteStatus_Done), []*model.Note{})
}

func TestNewTagRegistration(t *testing.T) {
	dataFilePath := path.Join("..", "..", "test", "test_data_file.json")
	reminderData, _ := model.ReadDataFile(dataFilePath, false)
	utils.AssertEqual(t, len(reminderData.Tags), 5)
}

func TestNotes(t *testing.T) {
	var notes []*model.Note
	notes = append(notes, &model.Note{Text: "1", Status: model.NoteStatus_Pending, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}})
	notes = append(notes, &model.Note{Text: "2", Status: model.NoteStatus_Pending, BaseStruct: model.BaseStruct{UpdatedAt: 1600000004}})
	notes = append(notes, &model.Note{Text: "3", Status: model.NoteStatus_Done, BaseStruct: model.BaseStruct{UpdatedAt: 1600000003}})
	notes = append(notes, &model.Note{Text: "4", Status: model.NoteStatus_Done, BaseStruct: model.BaseStruct{UpdatedAt: 1600000002}})
	sort.Sort(model.Notes(notes))
	var gotTexts []string
	for _, value := range notes {
		gotTexts = append(gotTexts, value.Text)
	}
	wantTexts := []string{"2", "3", "4", "1"}
	utils.AssertEqual(t, gotTexts, wantTexts)
}

func TestCreateDataFile(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath, false)
	reminderData, _ := model.ReadDataFile(dataFilePath, false)
	// old_updated_at := reminderData.UpdatedAt
	testUser := model.User{Name: "Test User", EmailId: "user@test.com"}
	reminderData.User = &testUser
	_ = reminderData.CreateDataFile("")
	remiderDataRe, _ := model.ReadDataFile(dataFilePath, false)
	// utils.AssertEqual(t, remiderDataRe.UpdatedAt > old_updated_at, true)
	utils.AssertEqual(t, remiderDataRe.User.EmailId == testUser.EmailId, true)
}

func TestUpdateDataFile(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath, false)
	reminderData, _ := model.ReadDataFile(dataFilePath, false)
	// old_updated_at := reminderData.UpdatedAt
	testUser := model.User{Name: "Test User", EmailId: "user@test.com"}
	reminderData.User = &testUser
	_ = reminderData.UpdateDataFile("")
	remiderDataRe, _ := model.ReadDataFile(dataFilePath, false)
	// utils.AssertEqual(t, remiderDataRe.UpdatedAt > old_updated_at, true)
	utils.AssertEqual(t, remiderDataRe.User.EmailId == testUser.EmailId, true)
}

func TestRegisterBasicTags(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath, false)
	reminderData, _ := model.ReadDataFile(dataFilePath, false)
	// register basic tags
	_ = reminderData.RegisterBasicTags()
	utils.AssertEqual(t, len(reminderData.Tags), 7)
}

func TestNotesApproachingDueDate(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath, false)
	reminderData, _ := model.ReadDataFile(dataFilePath, false)
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
	notes = append(notes, &model.Note{Text: "NRD01a", Status: model.NoteStatus_Done, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 8*24*3600})
	notes = append(notes, &model.Note{Text: "NRD02a", Status: model.NoteStatus_Done, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})
	notes = append(notes, &model.Note{Text: "NRD03a", Status: model.NoteStatus_Done, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600})
	// repeat-annually done notes
	notes = append(notes, &model.Note{Text: "RAD01a", Status: model.NoteStatus_Done, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 8*24*3600})
	notes = append(notes, &model.Note{Text: "RAD02a", Status: model.NoteStatus_Done, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})
	notes = append(notes, &model.Note{Text: "RAD03a", Status: model.NoteStatus_Done, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600})
	// repeat-monthally done notes
	notes = append(notes, &model.Note{Text: "RMD01a", Status: model.NoteStatus_Done, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 8*24*3600})
	notes = append(notes, &model.Note{Text: "RMD02a", Status: model.NoteStatus_Done, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})
	notes = append(notes, &model.Note{Text: "RMD03a", Status: model.NoteStatus_Done, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600})
	// non-repeat pending notes
	notes = append(notes, &model.Note{Text: "NRP01a", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600})        // expected
	notes = append(notes, &model.Note{Text: "NRP02a", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 2*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP02b", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 2*24*3600 + 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP03a", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600})        // expected
	notes = append(notes, &model.Note{Text: "NRP04a", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP04b", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 + 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP05a", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 3*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP05b", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 3*24*3600 + 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP06a", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 7*24*3600 - 3600}) // expected
	notes = append(notes, &model.Note{Text: "NRP06b", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 7*24*3600 + 3600})
	notes = append(notes, &model.Note{Text: "NRP07a", Status: model.NoteStatus_Pending, TagIds: []int{currentTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600})
	// repeat-annually pending notes
	notes = append(notes, &model.Note{Text: "RAP01", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600})
	notes = append(notes, &model.Note{Text: "RAP02", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RAP03", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime})             // expected
	notes = append(notes, &model.Note{Text: "RAP04", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RAP05", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RAP06", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600})
	notes = append(notes, &model.Note{Text: "RAP07", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600 - 2*365})
	notes = append(notes, &model.Note{Text: "RAP08", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600 - 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP09", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 2*365})             // expected
	notes = append(notes, &model.Note{Text: "RAP10", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600 - 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP11", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 - 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP12", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600 - 2*365})
	notes = append(notes, &model.Note{Text: "RAP13", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600 + 2*365})
	notes = append(notes, &model.Note{Text: "RAP14", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600 + 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP15", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*365})             // expected
	notes = append(notes, &model.Note{Text: "RAP16", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 1*24*3600 + 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP17", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600 + 2*365}) // expected
	notes = append(notes, &model.Note{Text: "RAP18", Status: model.NoteStatus_Pending, TagIds: []int{repeatAnnuallyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600 + 2*365})
	// repeat-monthly pending notes
	notes = append(notes, &model.Note{Text: "RMP01", Status: model.NoteStatus_Pending, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 9*24*3600})
	notes = append(notes, &model.Note{Text: "RMP02", Status: model.NoteStatus_Pending, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime - 1*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RMP03", Status: model.NoteStatus_Pending, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime})             // expected
	notes = append(notes, &model.Note{Text: "RMP04", Status: model.NoteStatus_Pending, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 2*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RMP05", Status: model.NoteStatus_Pending, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 3*24*3600}) // expected
	notes = append(notes, &model.Note{Text: "RMP06", Status: model.NoteStatus_Pending, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 7*24*3600})
	notes = append(notes, &model.Note{Text: "RMP07", Status: model.NoteStatus_Pending, TagIds: []int{repeatMonthlyTagId}, BaseStruct: model.BaseStruct{UpdatedAt: 1600000001}, CompleteBy: currentTime + 9*24*3600})
	reminderData.Notes = notes
	// get urgent notes
	urgentNotes := reminderData.NotesApprachingDueDate("default")
	var urgentNotesText []string
	for _, note := range urgentNotes {
		urgentNotesText = append(urgentNotesText, note.Text)
	}
	expectNotesText := []string{
		"NRP01a", "NRP02a", "NRP02b", "NRP03a", "NRP04a", "NRP04b", "NRP05a", "NRP05b", "NRP06a",
		"RAP02", "RAP03", "RAP04", "RAP05", "RAP08", "RAP09", "RAP10", "RAP11", "RAP14", "RAP15", "RAP16", "RAP17",
		"RMP02", "RMP03",
	}
	t.Logf("Received Texts: %v", urgentNotesText)
	t.Logf("Expected Texts: %v", expectNotesText)
	utils.SkipCI(t)
	utils.AssertEqual(t, urgentNotesText, expectNotesText)
	// [NRP01a NRP02a NRP02b NRP03a NRP04a NRP04b NRP05a NRP05b NRP06a RAP02 RAP03 RAP04 RAP05 RAP08 RAP09 RAP10 RAP11 RAP14 RAP15 RAP16 RAP17 RMP03]
	// [NRP01a NRP02a NRP02b NRP03a NRP04a NRP04b NRP05a NRP05b NRP06a RAP02 RAP03 RAP04 RAP05 RAP08 RAP09 RAP10 RAP11 RAP14 RAP15 RAP16 RAP17 RMP02 RMP03]}
}

func TestPrintStats(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))
	// create the file and required dirs
	_ = model.MakeSureFileExists(dataFilePath, false)
	reminderData, _ := model.ReadDataFile(dataFilePath, false)
	// register basic tags
	_ = reminderData.RegisterBasicTags()
	got, _ := reminderData.Stats()
	want := `
Stats of "temp_test_dir/mydata.json":
  - Number of Tags:  7
  - Pending Notes:   0/0
  - Suspended Notes: 0
  - Done Notes:      0
`
	utils.AssertEqual(t, got, want)
}
