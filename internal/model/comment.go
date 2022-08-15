package model

import (
	"html/template"

	"github.com/goyalmunish/reminder/pkg/utils"
)

/*
A Comment is an update to a note.

Consider it a statement representing an action to be taken/done or just an update about the Note.

A comment belongs to a particular note,
whereas a note can have multiple comments
*/
type Comment struct {
	Text string `json:"text"`
	BaseStruct
}

// String provides basic string representation of a commment.
func (comment *Comment) String() string {
	reportTemplate := `[{{.CreatedAt | mediumTimeStr}}] {{.Text}}`
	funcMap := template.FuncMap{
		"mediumTimeStr": utils.UnixTimestampToMediumTimeStr,
	}
	return utils.TemplateResult(reportTemplate, funcMap, comment)
}
