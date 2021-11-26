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
			Validate: utils.ValidateNonEmptyString,
		}
	case "user_email":
		prompt = &promptui.Prompt{
			Label:    "User Email",
			Validate: utils.ValidateNonEmptyString,
		}
	case "tag_slug":
		prompt = &promptui.Prompt{
			Label:    "Tag Slug",
			Validate: utils.ValidateNonEmptyString,
		}
	case "tag_group":
		prompt = &promptui.Prompt{
			Label:    "Tag Group",
			Validate: utils.ValidateString,
		}
	case "tag_another":
		prompt = &promptui.Prompt{
			Label:    "Add another tag: yes/no (default: no):",
			Validate: utils.ValidateString,
		}
	case "note_text":
		prompt = &promptui.Prompt{
			Label:    "Note Text",
			Validate: utils.ValidateNonEmptyString,
		}
	case "note_text_with_default":
		prompt = &promptui.Prompt{
			Label:    "New Text",
			Default:  defaultText,
			Validate: utils.ValidateNonEmptyString,
		}
	case "note_comment":
		prompt = &promptui.Prompt{
			Label:    "New Comment",
			Validate: utils.ValidateNonEmptyString,
		}
	case "note_completed_by":
		prompt = &promptui.Prompt{
			Label:    "Due Date (format: YYYY-MM-DD), or enter nil to clear existing value",
			Validate: utils.ValidateDateString,
		}
	}
	return prompt
}
