package models

import (
	"html/template"
	"reminder/pkg/utils"
)

/*
Comment is an update to a note

A comment belongs to a particular note
A note can have multiple comments
*/
type Comment struct {
	Text string `json:"text"`
	BaseStruct
}

type Comments []*Comment

func (c Comments) Len() int           { return len(c) }
func (c Comments) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Comments) Less(i, j int) bool { return c[i].CreatedAt > c[j].CreatedAt }

// provide basic string representation of a commment
func (comment *Comment) String() string {
	reportTemplate := `[{{.CreatedAt | mediumTimeStr}}] {{.Text}}`
	funcMap := template.FuncMap{
		"mediumTimeStr": utils.UnixTimestampToMediumTimeStr,
	}
	return utils.TemplateResult(reportTemplate, funcMap, comment)
}

// provide basic string representation of commments
func (comments Comments) Strings() []string {
	// assuming each note will have 10 comments on average
	strs := make([]string, 0, 10)
	for _, comment := range comments {
		strs = append(strs, comment.String())
	}
	return strs
}
