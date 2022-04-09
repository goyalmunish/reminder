package model

/*
A Tags is a slice of Tag objects.

By default it is sorted by its Slug field.
*/
type Tags []*Tag

func (c Tags) Len() int           { return len(c) }
func (c Tags) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Tags) Less(i, j int) bool { return c[i].Slug < c[j].Slug }

// Slugs returns slugs of given tags.
func (tags Tags) Slugs() []string {
	// assuming there are at least 20 tags (on average)
	allSlugs := make([]string, 0, 20)
	for _, tag := range tags {
		allSlugs = append(allSlugs, tag.Slug)
	}
	return allSlugs
}

// FromSlug fetches tag with given slug.
// It return nil if given tag is not found.
func (tags Tags) FromSlug(slug string) *Tag {
	for _, tag := range tags {
		if tag.Slug == slug {
			return tag
		}
	}
	return nil
}

// FromIds returns tags from tagIDs.
// It returns empty Tags if non of tagIDs match.
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

// IdsForGroup returns tag ids of given group.
// It returns empty []int if group with given group name doesn't exist.
func (tags Tags) IdsForGroup(group string) []int {
	var tagIDs []int
	for _, tag := range tags {
		if tag.Group == group {
			tagIDs = append(tagIDs, tag.Id)
		}
	}
	return tagIDs
}
