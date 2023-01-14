package model_test

import (
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func TestTagString(t *testing.T) {
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
