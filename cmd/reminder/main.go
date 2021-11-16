package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
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
	// open the file
	reminderData := *models.FReadDataFile(defaultDataFilePath)
	// print data
	if len(reminderData.Tags) > 0 {
		fmt.Printf("\nStats of %q\n", reminderData.DataFile)
		fmt.Printf("%4vNumber of Tags: %v\n", "- ", len(reminderData.Tags))
		fmt.Printf("%4vPending Notes: %v/%v\n", "- ", len(reminderData.Notes.WithStatus("pending")), len(reminderData.Notes))
	}
	// TEST CODE STARTS ---->
	// TEST CODE ENDS ------>
	// ask the main menu
	fmt.Println("| =========================== MAIN MENU =========================== |")
	fmt.Println("|    Use 'Ctrl-c' to jump from any nested level to the main menu    |")
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
		var allTagsSlugsWithEmoji []string
		for _, tagSlug := range reminderData.TagsSlugs() {
			allTagsSlugsWithEmoji = append(allTagsSlugsWithEmoji, fmt.Sprintf("%v %v", utils.Symbols["tag"], tagSlug))
		}
		tagIndex, _ := utils.AskOption(append(allTagsSlugsWithEmoji, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag")
		if tagIndex != -1 {
			if tagIndex == len(reminderData.TagsSlugs()) {
				// add new tag
				reminderData.NewTagRegistration()
			} else {
				tag := reminderData.Tags[tagIndex]
				notes := reminderData.NotesWithTagId(tag.Id, "pending")
				err := reminderData.PrintNotesAndAskOptions(notes, tag.Id)
				utils.PrintErrorIfPresent(err)
			}
		}
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Note"):
		tagIDs := reminderData.AskTagIds([]int{})
		note := models.FNewNote(tagIDs)
		reminderData.NewNoteAppend(note)
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Tag"):
		reminderData.NewTagRegistration()
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Register Basic Tags"):
		reminderData.RegisterBasicTags()
	case fmt.Sprintf("%v %v", utils.Symbols["clock"], "Pending Notes"):
		allNotes := reminderData.Notes
		pendingNotes := allNotes.WithStatus("pending")
		var currentNotes []*models.Note
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
		allNotes := reminderData.Notes
		var allTexts []string
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
		prompt := promptui.Select{
			Label:             "Notes",
			Items:             allTexts,
			Size:              25,
			StartInSearchMode: true,
			Searcher:          searchNotes,
		}
		fmt.Printf("Searching through a total of %v notes:\n", len(allTexts))
		index, _, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		if index >= 0 {
			note := allNotes[index]
			reminderData.PrintNoteAndAskOptions(note)
		}
	case fmt.Sprintf("%v %v", utils.Symbols["backup"], "Create Backup"):
		// get backup file name
		ext := path.Ext(reminderData.DataFile)
		dstFile := reminderData.DataFile[:len(reminderData.DataFile)-len(ext)] + "_backup_" + strconv.Itoa(int(utils.CurrentUnixTimestamp())) + ext
		lnFile := reminderData.DataFile[:len(reminderData.DataFile)-len(ext)] + "_backup_latest" + ext
		fmt.Printf("Creating backup at %q\n", dstFile)
		// create backup
		byteValue, err := ioutil.ReadFile(reminderData.DataFile)
		utils.PrintErrorIfPresent(err)
		err = ioutil.WriteFile(dstFile, byteValue, 0644)
		utils.PrintErrorIfPresent(err)
		// create alias of latest backup
		fmt.Printf("Creating synlink at %q\n", lnFile)
		executable, _ := exec.LookPath("ln")
		cmd := &exec.Cmd{
			Path:   executable,
			Args:   []string{executable, "-f", dstFile, lnFile},
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
		}
		err = cmd.Run()
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v", utils.Symbols["pad"], "Display Data File"):
		fmt.Printf("Printing contents (and if possible, its difference since last backup) of %q:\n", reminderData.DataFile)
		ext := path.Ext(reminderData.DataFile)
		lnFile := reminderData.DataFile[:len(reminderData.DataFile)-len(ext)] + "_backup_latest" + ext
		err := utils.PerformWhich("wdiff")
		if err == nil {
			err = utils.PerformFilePresence(lnFile)
			if err == nil {
				err = utils.FPerformCwdiff(lnFile, reminderData.DataFile)
			} else {
				fmt.Printf("Warning: `%v` file is not available yet\n", lnFile)
				err = utils.PerformCat(reminderData.DataFile)
			}
		} else {
			fmt.Printf("%v Warning: `wdiff` command is not available\n", utils.Symbols["error"])
			err = utils.PerformCat(reminderData.DataFile)
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
