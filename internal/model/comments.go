package model

/*
A Comments is a slice of Comment objects.

By default it is sorted by its CreatedAt field
*/
type Comments []*Comment

func (c Comments) Len() int           { return len(c) }
func (c Comments) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Comments) Less(i, j int) bool { return c[i].CreatedAt > c[j].CreatedAt }

// Strings provides representation of Commments in terms of slice of strings.
func (comments Comments) Strings() []string {
	// assuming each note will have 10 comments on average
	strs := make([]string, 0, 10)
	for _, comment := range comments {
		s := comment.String()
		strs = append(strs, s)
	}
	return strs
}
