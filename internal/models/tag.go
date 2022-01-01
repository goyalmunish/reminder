package models

import (
	"errors"
	"fmt"
	"strings"

	"reminder/pkg/utils"
)

/*
Tag represents classification of a note

A note can have multiple tags
A tag can be associated with multiple notes
*/
type Tag struct {
	Id    int    `json:"id"`    // internal int-based id of the tag
	Slug  string `json:"slug"`  // client-facing string-based id for tag
	Group string `json:"group"` // a note can be part of only one tag within a group
	BaseStruct
}

// provide basic string representation of a tag
func (t Tag) String() string {
	return fmt.Sprintf("%v#%v#%v", t.Group, t.Slug, t.Id)
}

// collection of tags with a defined default way of sorting
type Tags []*Tag

func (c Tags) Len() int           { return len(c) }
func (c Tags) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Tags) Less(i, j int) bool { return c[i].Slug < c[j].Slug }

// get slugs of given tags
func (tags Tags) Slugs() []string {
	// assuming there are at least 20 tags (on average)
	allSlugs := make([]string, 0, 20)
	for _, tag := range tags {
		allSlugs = append(allSlugs, tag.Slug)
	}
	return allSlugs
}

// fetch tag with given slug
// return nil if given tag is not found
func (tags Tags) FromSlug(slug string) *Tag {
	for _, tag := range tags {
		if tag.Slug == slug {
			return tag
		}
	}
	return nil
}

// get tags from tagIDs
// return empty Tags if non of tagIDs match
func (tags Tags) FromIds(tagIDs []int) Tags {
	var filteredTags Tags
	for _, tagID := range tagIDs {
		for _, tag := range tags {
			if tagID == tag.Id {
				filteredTags = append(filteredTags, tag)
			}
		}
	}
	return filteredTags
}

// get tag ids of given group
// return empty []int if group with given group name doesn't exist
func (tags Tags) IdsForGroup(group string) []int {
	var tagIDs []int
	for _, tag := range tags {
		if tag.Group == group {
			tagIDs = append(tagIDs, tag.Id)
		}
	}
	return tagIDs
}

// functions

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
			Id:    index,
			Slug:  tagMap["slug"],
			Group: tagMap["group"],
			BaseStruct: BaseStruct{
				CreatedAt: utils.CurrentUnixTimestamp(),
				UpdatedAt: utils.CurrentUnixTimestamp()},
		}
		basicTags = append(basicTags, &tag)
	}
	return basicTags
}

// prompt for new Tag
func FNewTag(tagID int, promptTagSlug Prompter, promptTagGroup Prompter) (*Tag, error) {
	tag := &Tag{
		Id: tagID,
		BaseStruct: BaseStruct{
			CreatedAt: utils.CurrentUnixTimestamp(),
			UpdatedAt: utils.CurrentUnixTimestamp()},
		// Slug:      tagSlug,
		// Group:     tagGroup,
	}
	// ask for tag slug
	tagSlug, err := promptTagSlug.Run()
	tag.Slug = utils.TrimString(tagSlug)
	tag.Slug = strings.ToLower(tag.Slug)
	// in case of error or Ctrl-c as input, don't create the tag
	if err != nil || strings.Contains(tag.Slug, "^C") {
		return tag, err
	}
	if len(utils.TrimString(tag.Slug)) == 0 {
		// this should never be encountered because of validation in earlier step
		fmt.Printf("%v Skipping adding tag with empty slug\n", utils.Symbols["warning"])
		err := errors.New("Tag's slug is empty")
		return tag, err
	}
	// ask for tag's group
	tagGroup, err := promptTagGroup.Run()
	if err != nil {
		return tag, err
	}
	tag.Group = strings.ToLower(tagGroup)
	// return successful tag
	return tag, nil
}
