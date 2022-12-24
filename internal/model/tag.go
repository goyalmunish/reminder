package model

import (
	"context"
	"fmt"
)

/*
A Tag represents classification of a note.

A note can have multiple tags, and a tag can be associated with multiple notes.
*/
type Tag struct {
	context context.Context
	Id      int    `json:"id"`    // internal int-based id of the tag
	Slug    string `json:"slug"`  // client-facing string-based id for tag
	Group   string `json:"group"` // a note can be part of only one tag within a group
	BaseStruct
}

// String provides basic string representation of a tag.
func (t Tag) String() string {
	return fmt.Sprintf("%v#%v#%v", t.Group, t.Slug, t.Id)
}

func (tag *Tag) SetContext(ctx context.Context) {
	tag.context = ctx
}
