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
type BySlug []*Tag

func (c BySlug) Len() int           { return len(c) }
func (c BySlug) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c BySlug) Less(i, j int) bool { return c[i].Slug < c[j].Slug }

// get slugs of given tags
func FTagsSlugs(tags []*Tag) []string {
	var all_slugs []string
	for _, tag := range tags {
		all_slugs = append(all_slugs, tag.Slug)
	}
	return all_slugs
}

// return a array of basic tags
func FBasicTags() []*Tag {
	basic_tags_map := []map[string]string{{"slug": "current", "group": ""},
		{"slug": "priority-urgent", "group": "priority"},
		{"slug": "priority-medium", "group": "priority"},
		{"slug": "priority-low", "group": "priority"},
		{"slug": "repeat-annually", "group": "repeat"},
		{"slug": "repeat-monthly", "group": "repeat"},
		{"slug": "tips", "group": "tips"}}
	var basic_tags []*Tag
	for index, tag_map := range basic_tags_map {
		tag := Tag{
			Id:        index,
			Slug:      tag_map["slug"],
			Group:     tag_map["group"],
			CreatedAt: utils.CurrentUnixTimestamp(),
			UpdatedAt: utils.CurrentUnixTimestamp(),
		}
		fmt.Println(tag)
		basic_tags = append(basic_tags, &tag)
	}
	return basic_tags
}

// prompt for new Tag
func FNewTag(tag_id int) *Tag {
	prompt := promptui.Prompt{
		Label:    "Tag Slug",
		Validate: utils.ValidateNonEmptyString,
	}
	tag_slug, err := prompt.Run()
	utils.PrintErrorIfPresent(err)
	tag_slug = strings.ToLower(tag_slug)
	prompt = promptui.Prompt{
		Label:    "Tag Group",
		Validate: utils.ValidateString,
	}
	tag_group, err := prompt.Run()
	tag_group = strings.ToLower(tag_group)
	utils.PrintErrorIfPresent(err)
	return &Tag{
		Id:        tag_id,
		Slug:      tag_slug,
		Group:     tag_group,
		CreatedAt: utils.CurrentUnixTimestamp(),
		UpdatedAt: utils.CurrentUnixTimestamp(),
	}
}
