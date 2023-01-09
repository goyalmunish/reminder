package utils

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/goyalmunish/reminder/pkg/logger"
)

// AskOption function asks option to the user.
// It print error, if encountered any (so that they don't have to printed by calling function).
// It return a tuple (chosen index, chosen string, err if any).
func AskOption(options []string, label string) (int, string, error) {
	if len(options) == 0 {
		err := errors.New("Empty List")
		fmt.Printf("%v Prompt failed %v\n", Symbols["warning"], err)
		return -1, "", err
	}
	// note: any item in options should not have \n character
	// otherwise such item is observed to not getting appear
	// in the rendered list
	var selectedIndex int
	prompt := &survey.Select{
		Message:  label,
		Options:  options,
		PageSize: 25,
		VimMode:  true,
	}
	err := survey.AskOne(prompt, &selectedIndex)
	if err != nil {
		// error can happen if user raises an interrupt (such as Ctrl-c, SIGINT)
		fmt.Printf("%v Prompt failed %v\n", Symbols["warning"], err)
		return -1, "", err
	}
	logger.Info(fmt.Sprintf("You chose %d:%q\n", selectedIndex, options[selectedIndex]))
	return selectedIndex, options[selectedIndex], nil
}

// GeneratePrompt function generates survey.Input.
func GeneratePrompt(promptName string, defaultText string) (string, error) {
	var validator survey.Validator
	var answer string
	var err error

	switch promptName {
	case "user_name":
		prompt := &survey.Input{
			Message: "User Name: ",
			Default: defaultText,
		}
		validator = survey.Required
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "user_email":
		prompt := &survey.Input{
			Message: "User Email: ",
			Default: defaultText,
		}
		validator = survey.MinLength(0)
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "tag_slug":
		prompt := &survey.Input{
			Message: "Tag Slug: ",
			Default: defaultText,
		}
		validator = survey.MinLength(1)
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "tag_group":
		prompt := &survey.Input{
			Message: "Tag Group: ",
			Default: defaultText,
		}
		validator = survey.MinLength(1)
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "tag_another":
		prompt := &survey.Input{
			Message: "Add another tag: yes/no (default: no): ",
			Default: defaultText,
		}
		validator = survey.MinLength(1)
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "note_text":
		prompt := &survey.Input{
			Message: "Note Text: ",
			Default: defaultText,
		}
		validator = survey.MinLength(1)
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "note_summary":
		prompt := &survey.Multiline{
			Message: "Note Summary: ",
			Default: defaultText,
		}
		validator = survey.MinLength(1)
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "note_comment":
		prompt := &survey.Multiline{
			Message: "New Comment: ",
			Default: defaultText,
		}
		validator = survey.MinLength(1)
		err = survey.AskOne(prompt, &answer, survey.WithValidator(validator))
	case "note_completed_by":
		prompt := &survey.Input{
			Message: "Due Date (format: DD-MM-YYYY or DD-MM), or enter nil to clear existing value: ",
			Default: defaultText,
		}
		err = survey.AskOne(prompt, &answer, survey.WithValidator(ValidateDateString()))
	}
	return answer, err
}

// GenerateNoteSearchSelect function generates survey.Select and return index of selected option.
func GenerateNoteSearchSelect(items []string, searchFunc func(filter string, value string, index int) bool) (int, error) {
	var selectedIndex int
	prompt := &survey.Select{
		Message:  "Search: ",
		Options:  items,
		PageSize: 25,
		Filter:   searchFunc,
		VimMode:  true,
	}
	err := survey.AskOne(prompt, &selectedIndex)
	return selectedIndex, err
}
