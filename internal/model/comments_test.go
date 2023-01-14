package model_test

import (
	"sort"
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func TestCommentsStrings(t *testing.T) {
	utils.Location = utils.UTCLocation()
	comments := model.Comments{&model.Comment{Text: "c1:\n- line 1\n\n- line 2\n- line 3 with \" and < characters"}, &model.Comment{Text: "c2"}, &model.Comment{Text: "c3"}}
	want := []string{"nil | c1:\n- line 1\n\n- line 2\n- line 3 with \" and < characters", "nil | c2", "nil | c3"}
	utils.AssertEqual(t, comments.Strings(), want)
}

func TestCommentsSort(t *testing.T) {
	utils.Location = utils.UTCLocation()
	c1 := &model.Comment{Text: "c1", BaseStruct: model.BaseStruct{CreatedAt: 1600000004}}
	c2 := &model.Comment{Text: "c2", BaseStruct: model.BaseStruct{CreatedAt: 1600000002}}
	c3 := &model.Comment{Text: "c3", BaseStruct: model.BaseStruct{CreatedAt: 1600000003}}
	comments := model.Comments{c1, c2, c3}
	sort.Sort(comments)
	utils.AssertEqual(t, comments, model.Comments{c1, c3, c2})
}
