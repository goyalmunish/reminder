package models

/*
A Prompter representing Prompt which can be Run.
*/
type Prompter interface {
	Run() (string, error)
}
