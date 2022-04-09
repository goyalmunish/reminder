package model

import (
	"html/template"
	"reminder/pkg/utils"
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

/*
A Comments is a slice of Comment objects.

By default it is sorted by its CreatedAt field
*/
type Comments []*Comment

func (c Comments) Len() int           { return len(c) }
func (c Comments) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Comments) Less(i, j int) bool { return c[i].CreatedAt > c[j].CreatedAt }

// String provides basic string representation of a commment.
func (comment *Comment) String() string {
	reportTemplate := `[{{.CreatedAt | mediumTimeStr}}] {{.Text}}`
	funcMap := template.FuncMap{
		"mediumTimeStr": utils.UnixTimestampToMediumTimeStr,
	}
	return utils.TemplateResult(reportTemplate, funcMap, comment)
}

// Strings provides representation of Commments in terms of slice of strings.
func (comments Comments) Strings() []string {
	// assuming each note will have 10 comments on average
	strs := make([]string, 0, 10)
	for _, comment := range comments {
		strs = append(strs, comment.String())
	}
	return strs
}
