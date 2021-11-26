package models

import (
	"reminder/pkg/utils"

	"github.com/manifoldco/promptui"
)

type PromptInf interface {
	Run() (string, error)
}

func GeneratePrompt(promptName string, defaultText string) *promptui.Prompt {
	var prompt *promptui.Prompt
	switch promptName {
	case "user_name":
		prompt = &promptui.Prompt{
			Label:    "User Name",
			Default:  defaultText,
			Validate: utils.ValidateNonEmptyString,
		}
	case "user_email":
		prompt = &promptui.Prompt{
			Label:    "User Email",
			Default:  defaultText,
			Validate: utils.ValidateNonEmptyString,
		}
	case "tag_slug":
		prompt = &promptui.Prompt{
			Label:    "Tag Slug",
			Default:  defaultText,
			Validate: utils.ValidateNonEmptyString,
		}
	case "tag_group":
		prompt = &promptui.Prompt{
			Label:    "Tag Group",
			Default:  defaultText,
			Validate: utils.ValidateString,
		}
	case "tag_another":
		prompt = &promptui.Prompt{
			Label:    "Add another tag: yes/no (default: no):",
			Default:  defaultText,
			Validate: utils.ValidateString,
		}
	case "note_text":
		prompt = &promptui.Prompt{
			Label:    "Note Text",
			Default:  defaultText,
			Validate: utils.ValidateNonEmptyString,
		}
	case "note_comment":
		prompt = &promptui.Prompt{
			Label:    "New Comment",
			Default:  defaultText,
			Validate: utils.ValidateNonEmptyString,
		}
	case "note_completed_by":
		prompt = &promptui.Prompt{
			Label:    "Due Date (format: YYYY-MM-DD), or enter nil to clear existing value",
			Default:  defaultText,
			Validate: utils.ValidateDateString,
		}
	}
	return prompt
}

func GenerateNoteSearchSelect(items []string, searchFunc func(string, int) bool) *promptui.Select {
	prompt := &promptui.Select{
		Label:             "Notes",
		Items:             items,
		Size:              25,
		StartInSearchMode: true,
		Searcher:          searchFunc,
	}
	return prompt
}
