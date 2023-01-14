package model_test

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"
	"time"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

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

func TestBasicTags(t *testing.T) {
	basicTags := model.BasicTags()
	slugs := basicTags.Slugs()
	want := "[current priority-urgent priority-medium priority-low repeat-annually repeat-monthly tips]"
	utils.AssertEqual(t, slugs, want)
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

func TestMakeSureFileExists(t *testing.T) {
	var dataFilePath = "temp_test_dir/mydata.json"
	// make sure temporary files and dirs are removed at the end of the test
	defer os.RemoveAll(path.Dir(dataFilePath))

	// make sure file doesn't exists already
	_, err := os.Stat(dataFilePath)
	utils.AssertEqual(t, err != nil, true)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), true)
	// attempt to create the file and required dirs, when the file doesn't exist already
	_ = model.MakeSureFileExists(dataFilePath, false)
	// prove that the file was created
	stats, err := os.Stat(dataFilePath)
	utils.AssertEqual(t, err != nil, false)
	utils.AssertEqual(t, errors.Is(err, fs.ErrNotExist), false)

	// make sure that the existing file is not replaced
	modificationTime := stats.ModTime()
	// attempt to create the file and required dirs, when the file does exist already
	time.Sleep(10 * time.Millisecond)
	_ = model.MakeSureFileExists(dataFilePath, false)
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
	_ = model.MakeSureFileExists(dataFilePath, false)
	// attempt to read file and parse it
	reminderData, _ := model.ReadDataFile(dataFilePath, false)
	utils.AssertEqual(t, reminderData.UpdatedAt > 0, true)
}
