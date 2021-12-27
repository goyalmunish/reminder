package models

/*
Prompter is an interface representing Prompt which can be Run
*/
type Prompter interface {
	Run() (string, error)
}
