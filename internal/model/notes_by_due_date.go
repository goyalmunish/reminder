package model

/*
A NotesByDueDate is a slice of Note objects.

By default it is sorted by its CreatedAt field.
*/
type NotesByDueDate []*Note

func (c NotesByDueDate) Len() int           { return len(c) }
func (c NotesByDueDate) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c NotesByDueDate) Less(i, j int) bool { return c[i].CompleteBy < c[j].CompleteBy }
