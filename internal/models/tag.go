package models

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"

	"reminder/pkg/utils"
)

type Tag struct {
	Id        int    `json:"id"`    // internal int-based id of the tag
	Slug      string `json:"slug"`  // client-facing string-based id for tag
	Group     string `json:"group"` // a note can be part of only one tag within a group
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// method to provide basic string representation of a tag
func (t Tag) String() string {
	return fmt.Sprintf("%v#%v#%v", t.Group, t.Slug, t.Id)
}

// collection of tags with a defined default way of sorting
type Tags []*Tag

func (c Tags) Len() int           { return len(c) }
func (c Tags) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Tags) Less(i, j int) bool { return c[i].Slug < c[j].Slug }

// return an array of basic tags
// which can be used for initial setup of the application
// here some of the tags will have special meaning/functionality
// such as repeat-annually and repeat-monthly
func FBasicTags() Tags {
	basicTagsMap := []map[string]string{{"slug": "current", "group": ""},
		{"slug": "priority-urgent", "group": "priority"},
		{"slug": "priority-medium", "group": "priority"},
		{"slug": "priority-low", "group": "priority"},
		{"slug": "repeat-annually", "group": "repeat"},
		{"slug": "repeat-monthly", "group": "repeat"},
		{"slug": "tips", "group": "tips"}}
	var basicTags Tags
	for index, tagMap := range basicTagsMap {
		tag := Tag{
			Id:        index,
			Slug:      tagMap["slug"],
			Group:     tagMap["group"],
			CreatedAt: utils.CurrentUnixTimestamp(),
			UpdatedAt: utils.CurrentUnixTimestamp(),
		}
		basicTags = append(basicTags, &tag)
	}
	return basicTags
}

// get slugs of given tags
func FTagsSlugs(tags Tags) []string {
	var allSlugs []string
	for _, tag := range tags {
		allSlugs = append(allSlugs, tag.Slug)
	}
	return allSlugs
}

// prompt for new Tag
func FNewTag(tagID int) *Tag {
	prompt := promptui.Prompt{
		Label:    "Tag Slug",
		Validate: utils.ValidateNonEmptyString,
	}
	tagSlug, err := prompt.Run()
	utils.PrintErrorIfPresent(err)
	tagSlug = strings.ToLower(tagSlug)
	prompt = promptui.Prompt{
		Label:    "Tag Group",
		Validate: utils.ValidateString,
	}
	tagGroup, err := prompt.Run()
	tagGroup = strings.ToLower(tagGroup)
	utils.PrintErrorIfPresent(err)
	return &Tag{
		Id:        tagID,
		Slug:      tagSlug,
		Group:     tagGroup,
		CreatedAt: utils.CurrentUnixTimestamp(),
		UpdatedAt: utils.CurrentUnixTimestamp(),
	}
}
