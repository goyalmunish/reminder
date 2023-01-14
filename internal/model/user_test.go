package model_test

import (
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func TestUserString(t *testing.T) {
	user := model.User{Name: "Foo Bar", EmailId: "user@test.com"}
	want := "{Name: Foo Bar, EmailId: user@test.com}"
	utils.AssertEqual(t, user.String(), want)
}
