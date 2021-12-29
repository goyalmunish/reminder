/*
Tool `reminder` is a command-line (terminal) based interactive app for organizing tasks with minimal efforts.

Just run it as `go run cmd/reminder/main.go`
*/
package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"reminder/internal/models"
	"reminder/pkg/utils"
)

// recursive function for overall flow
func flow() {
	// make sure DataFile exists
	defaultDataFilePath := models.FDefaultDataFile()
	models.FMakeSureFileExists(defaultDataFilePath)
	// read and parse the existing data
	reminderData := *models.FReadDataFile(defaultDataFilePath)
	// print data stats
	fmt.Println(reminderData.Stats())
	// try automatic backup
	reminderData.AutoBackup(24 * 60 * 60)
	// TEST CODE STARTS ---->
	// TEST CODE ENDS ------>
	// ask the main menu
	fmt.Println("| =========================== MAIN MENU =========================== |")
	fmt.Println("|     Use 'Ctrl-c' to jump one level up (towards the Main Menu)     |")
	fmt.Println("| ----------------------------------------------------------------- |")
	_, result := utils.AskOption([]string{fmt.Sprintf("%v %v", utils.Symbols["spark"], "List Stuff"),
		fmt.Sprintf("%v %v %v", utils.Symbols["checkerdFlag"], "Exit", utils.Symbols["redFlag"]),
		fmt.Sprintf("%v %v", utils.Symbols["clock"], "Urgent Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["done"], "Done Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["search"], "Search Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Note"),
		fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Tag"),
		fmt.Sprintf("%v %v", utils.Symbols["clip"], "Register Basic Tags"),
		fmt.Sprintf("%v %v", utils.Symbols["backup"], "Create Backup"),
		fmt.Sprintf("%v %v", utils.Symbols["pad"], "Display Data File")}, "Select Option")
	// operate on main options
	switch result {
	case fmt.Sprintf("%v %v", utils.Symbols["spark"], "List Stuff"):
		tagSymbol := func(tagSlug string) string {
			hasPendingNote := len(reminderData.FindNotesByTagSlug(tagSlug, "pending")) > 0
			if hasPendingNote {
				return utils.Symbols["tag"]
			} else {
				return utils.Symbols["zzz"]
			}
		}
		// assuming there are at least 20 tags (on average)
		allTagSlugsWithEmoji := make([]string, 0, 20)
		for _, tagSlug := range reminderData.SortedTagSlugs() {
			allTagSlugsWithEmoji = append(allTagSlugsWithEmoji, fmt.Sprintf("%v %v", tagSymbol(tagSlug), tagSlug))
		}
		tagIndex, _ := utils.AskOption(append(allTagSlugsWithEmoji, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag")
		if tagIndex != -1 {
			if tagIndex == len(reminderData.SortedTagSlugs()) {
				// add new tag
				_, _ = reminderData.NewTagRegistration()
			} else {
				tag := reminderData.Tags[tagIndex]
				err := reminderData.PrintNotesAndAskOptions(models.Notes{}, tag.Id, "pending")
				utils.PrintErrorIfPresent(err)
			}
		}
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Note"):
		tagIDs := reminderData.AskTagIds([]int{})
		_, _ = reminderData.NewNoteRegistration(tagIDs)
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Tag"):
		_, _ = reminderData.NewTagRegistration()
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Register Basic Tags"):
		reminderData.RegisterBasicTags()
	case fmt.Sprintf("%v %v", utils.Symbols["clock"], "Urgent Notes"):
		err := reminderData.PrintNotesAndAskOptions(models.Notes{}, -1, "pending")
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v", utils.Symbols["done"], "Done Notes"):
		err := reminderData.PrintNotesAndAskOptions(models.Notes{}, -1, "done")
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v", utils.Symbols["search"], "Search Notes"):
		// get texts of all notes
		sort.Sort(reminderData.Notes)
		allNotes := reminderData.Notes
		// assuming the search shows 25 items in general
		allTexts := make([]string, 0, 25)
		for _, note := range allNotes {
			allTexts = append(allTexts, note.SearchableText())
		}
		// function to search across notes
		searchNotes := func(input string, idx int) bool {
			input = strings.ToLower(input)
			noteText := allTexts[idx]
			if strings.Contains(strings.ToLower(noteText), input) {
				return true
			}
			return false
		}
		// display prompt
		promptNoteSelection := utils.GenerateNoteSearchSelect(utils.ChopStrings(allTexts, utils.TerminalWidth()-10), searchNotes)
		fmt.Printf("Searching through a total of %v notes:\n", len(allTexts))
		index, _, err := promptNoteSelection.Run()
		utils.PrintErrorIfPresent(err)
		if index >= 0 {
			note := allNotes[index]
			reminderData.PrintNoteAndAskOptions(note)
		}
	case fmt.Sprintf("%v %v", utils.Symbols["backup"], "Create Backup"):
		reminderData.CreateBackup()
	case fmt.Sprintf("%v %v", utils.Symbols["pad"], "Display Data File"):
		fmt.Printf("Printing contents (and if possible, its difference since last backup) of %q:\n", reminderData.DataFile)
		ext := path.Ext(reminderData.DataFile)
		lnFile := reminderData.DataFile[:len(reminderData.DataFile)-len(ext)] + "_backup_latest" + ext
		err := utils.PerformWhich("wdiff")
		if err != nil {
			fmt.Printf("%v Warning: `wdiff` command is not available\n", utils.Symbols["error"])
			err = utils.PerformCat(reminderData.DataFile)
		} else {
			err = utils.PerformFilePresence(lnFile)
			if err != nil {
				fmt.Printf("Warning: `%v` file is not available yet\n", lnFile)
				err = utils.PerformCat(reminderData.DataFile)
			} else {
				err = utils.FPerformCwdiff(lnFile, reminderData.DataFile)
			}
		}
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v %v", utils.Symbols["checkerdFlag"], "Exit", utils.Symbols["redFlag"]):
		fmt.Println("Exiting...")
		return
	}
	flow()
}

func main() {
	go utils.Spinner(100 * time.Millisecond)
	flow()
}
