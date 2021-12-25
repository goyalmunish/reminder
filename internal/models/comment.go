package models

import (
	"fmt"
	"reminder/pkg/utils"
	"strings"
)

type Comment struct {
	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at"`
}

type Comments []*Comment

func (c Comments) Len() int           { return len(c) }
func (c Comments) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Comments) Less(i, j int) bool { return c[i].CreatedAt > c[j].CreatedAt }

// provide basic string representation of a commment
func (comment *Comment) String() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("[%v]", utils.UnixTimestampToMediumTimeStr(comment.CreatedAt)))
	strs = append(strs, comment.Text)
	return strings.Join(strs, " ")
}

// provide basic string representation of commments
func (comments Comments) ToStrings() []string {
	var strs []string
	for _, comment := range comments {
		strs = append(strs, comment.String())
	}
	return strs
}
