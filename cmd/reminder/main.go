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
		fmt.Sprintf("%v %v", utils.Symbols["clock"], "Pending Notes"),
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
				notes := reminderData.FindNotesByTagId(tag.Id, "pending")
				err := reminderData.PrintNotesAndAskOptions(notes, tag.Id)
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
	case fmt.Sprintf("%v %v", utils.Symbols["clock"], "Pending Notes"):
		allNotes := reminderData.Notes
		pendingNotes := allNotes.WithStatus("pending")
		// assuming there are at least 100 notes (on average)
		currentNotes := make([]*models.Note, 0, 100)
		repeatTagIDs := reminderData.TagIdsForGroup("repeat")
		// populating currentNotes
		for _, note := range pendingNotes {
			noteIDsWithRepeat := utils.GetCommonMembersIntSlices(note.TagIds, repeatTagIDs)
			// first process notes without tag with group "repeat"
			// start showing such notes 7 days in advance from their due date, and until they are marked done
			minDay := note.CompleteBy - 7*24*60*60
			currentTimestamp := utils.CurrentUnixTimestamp()
			if (len(noteIDsWithRepeat) == 0) && (note.CompleteBy != 0) && (currentTimestamp >= minDay) {
				currentNotes = append(currentNotes, note)
			}
			// check notes with tag with group "repeat"
			// start showing notes with "repeat-annually" 7 days in advance
			// start showing notes with "repeat-monthly" 3 days in advance
			// don't show such notes after their due date is past by 2 day
			if (len(noteIDsWithRepeat) > 0) && (note.CompleteBy != 0) {
				// check for repeat-annually tag
				repeatAnnuallyTag := reminderData.TagFromSlug("repeat-annually")
				repeatMonthlyTag := reminderData.TagFromSlug("repeat-monthly")
				if (repeatAnnuallyTag != nil) && utils.IntPresentInSlice(repeatAnnuallyTag.Id, note.TagIds) {
					_, noteMonth, noteDay := utils.UnixTimestampToTime(note.CompleteBy).Date()
					noteTimestampCurrent := utils.UnixTimestampForCorrespondingCurrentYear(int(noteMonth), noteDay)
					noteTimestampPrevious := noteTimestampCurrent - 365*24*60*60
					noteTimestampNext := noteTimestampCurrent + 365*24*60*60
					daysBefore := int64(7)
					daysAfter := int64(2)
					if utils.IsTimeForRepeatNote(noteTimestampCurrent, noteTimestampPrevious, noteTimestampNext, daysBefore, daysAfter) {
						currentNotes = append(currentNotes, note)
					}
				}
				// check for repeat-monthly tag
				if (repeatMonthlyTag != nil) && utils.IntPresentInSlice(repeatMonthlyTag.Id, note.TagIds) {
					_, _, noteDay := utils.UnixTimestampToTime(note.CompleteBy).Date()
					noteTimestampCurrent := utils.UnixTimestampForCorrespondingCurrentYearMonth(noteDay)
					noteTimestampPrevious := noteTimestampCurrent - 30*24*60*60
					noteTimestampNext := noteTimestampCurrent + 30*24*60*60
					daysBefore := int64(3)
					daysAfter := int64(2)
					if utils.IsTimeForRepeatNote(noteTimestampCurrent, noteTimestampPrevious, noteTimestampNext, daysBefore, daysAfter) {
						currentNotes = append(currentNotes, note)
					}
				}
			}
		}
		fmt.Println("Note: Following are the pending notes with due date:")
		fmt.Println("  - within a week or already crossed (for non repeat-annually or repeat-monthly)")
		fmt.Println("  - within a week for repeat-annually and 2 days post due date (ignoring its year)")
		fmt.Println("  - within 3 days for repeat-monthly and 2 days post due date (ignoring its year and month)")
		err := reminderData.PrintNotesAndAskOptions(currentNotes, -1)
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v", utils.Symbols["done"], "Done Notes"):
		allNotes := reminderData.Notes
		doneNotes := allNotes.WithStatus("done")
		fmt.Printf("A total of %v notes marked as 'done':\n", len(doneNotes))
		err := reminderData.PrintNotesAndAskOptions(doneNotes, -1)
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
