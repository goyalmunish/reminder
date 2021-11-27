package models

type PromptInf interface {
	Run() (string, error)
}
