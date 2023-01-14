package model_test

import (
	"sort"
	"testing"

	model "github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func TestTagsSort(t *testing.T) {
	var tags model.Tags
	tags = append(tags, &model.Tag{Id: 1, Slug: "a", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "z", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 3, Slug: "c", Group: "tag_group1"})
	tags = append(tags, &model.Tag{Id: 4, Slug: "b", Group: "tag_group2"})
	sort.Sort(model.Tags(tags))
	var got []int
	for _, value := range tags {
		got = append(got, value.Id)
	}
	want := []int{1, 4, 3, 2}
	utils.AssertEqual(t, got, want)
}

func TestTagsSlugs(t *testing.T) {
	var tags model.Tags
	utils.AssertEqual(t, tags, "[]")
	// case 1 (no tags)
	utils.AssertEqual(t, tags.Slugs(), "[]")
	// case 2 (non-empty tags)
	tags = append(tags, &model.Tag{Id: 1, Slug: "tag_1", Group: "tag_group"})
	tags = append(tags, &model.Tag{Id: 2, Slug: "tag_2", Group: "tag_group"})
	tags = append(tags, &model.Tag{Id: 3, Slug: "tag_3", Group: "tag_group"})
	got := tags.Slugs()
	want := "[tag_1 tag_2 tag_3]"
	utils.AssertEqual(t, got, want)
}

func TestTagsFromSlug(t *testing.T) {
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1 (passing non-existing slug)
	utils.AssertEqual(t, tags.FromSlug("no_such_slug"), nil)
	// case 2 (passing tag which is part of another tag as well)
	utils.AssertEqual(t, tags.FromSlug("a"), &tag1)
	// case 3
	utils.AssertEqual(t, tags.FromSlug("a1"), &tag2)
}

func TestTagsFromIds(t *testing.T) {
	var tags model.Tags
	// creating tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "z", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "c", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1 (passing blank tagIDs)
	tagIDs := []int{}
	gotSlugs := tags.FromIds(tagIDs)
	wantSlugs := model.Tags{}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 2 (no matching tagIDs)
	tagIDs = []int{100, 101}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = model.Tags{}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 3 (two matching tagIDs)
	tagIDs = []int{1, 3}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = model.Tags{&tag1, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
	// case 4
	tagIDs = []int{1, 4, 2, 3}
	gotSlugs = tags.FromIds(tagIDs)
	wantSlugs = model.Tags{&tag1, &tag4, &tag2, &tag3}
	utils.AssertEqual(t, gotSlugs, wantSlugs)
}

func TestTagsIdsForGroup(t *testing.T) {
	// creating tags
	var tags model.Tags
	tag1 := model.Tag{Id: 1, Slug: "a", Group: "tag_group1"}
	tags = append(tags, &tag1)
	tag2 := model.Tag{Id: 2, Slug: "a1", Group: "tag_group1"}
	tags = append(tags, &tag2)
	tag3 := model.Tag{Id: 3, Slug: "a2", Group: "tag_group1"}
	tags = append(tags, &tag3)
	tag4 := model.Tag{Id: 4, Slug: "b", Group: "tag_group2"}
	tags = append(tags, &tag4)
	// case 1 (group with no such name)
	utils.AssertEqual(t, tags.IdsForGroup("tag_group_NO"), []int{})
	// case 1 (group with multiple tags)
	utils.AssertEqual(t, tags.IdsForGroup("tag_group1"), []int{1, 2, 3})
	// case 2 (group with single tag)
	utils.AssertEqual(t, tags.IdsForGroup("tag_group2"), []int{4})
}
