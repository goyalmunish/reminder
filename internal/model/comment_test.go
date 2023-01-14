package model_test

import (
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func TestCommentString(t *testing.T) {
	utils.Location = utils.UTCLocation()
	c := model.Comment{Text: "c1:\n- line 1\n\n- line 2\n- line 3 with \" and < characters", BaseStruct: model.BaseStruct{CreatedAt: 1600000004, UpdatedAt: 1600001004}}
	want := "13-Sep-20 12:26:44 | c1:\n- line 1\n\n- line 2\n- line 3 with \" and < characters"
	utils.AssertEqual(t, c.String(), want)
}
