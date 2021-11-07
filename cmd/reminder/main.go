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

var DataFile = path.Join(os.Getenv("HOME"), "reminder", "data.json")

// recursive function for overall flow
func flow() {
	// make sure DataFile exists
	models.FMakeSureFileExists(DataFile)
	// open the file
	reminderData := *models.FReadDataFile(DataFile)
	// print data
	if len(reminderData.Tags) > 0 {
		fmt.Printf("\nStats of %q\n", DataFile)
		fmt.Printf("%4vNumber of Tags: %v\n", "- ", len(reminderData.Tags))
		fmt.Printf("%4vPending Notes: %v/%v\n", "- ", len(models.FNotesWithStatus(reminderData.Notes, "pending")), len(reminderData.Notes))
	}
	// TEST CODE STARTS ---->
	// TEST CODE ENDS ------>
	// ask the main menu
	fmt.Println("| =========================== MAIN MENU =========================== |")
	fmt.Println("|    Use 'Ctrl-c' to jump from any nested level to the main menu    |")
	fmt.Println("| ----------------------------------------------------------------- |")
	_, result := utils.AskOption([]string{fmt.Sprintf("%v %v", utils.Symbols["spark"], "List Stuff"),
		fmt.Sprintf("%v %v %v", utils.Symbols["checkerd_flag"], "Exit", utils.Symbols["red_flag"]),
		fmt.Sprintf("%v %v", utils.Symbols["clock"], "Pending Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["done"], "Done Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["search"], "Search Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Note"),
		fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Tag"),
		fmt.Sprintf("%v %v", utils.Symbols["clip"], "Register Basic Tags"),
		fmt.Sprintf("%v %v", utils.Symbols["pad"], "Create Backup"),
		fmt.Sprintf("%v %v", utils.Symbols["pad"], "Display Data File")}, "Select Option")
	// operate on main options
	switch result {
	case fmt.Sprintf("%v %v", utils.Symbols["spark"], "List Stuff"):
		var all_tags_slugs_with_emoji []string
		for _, tag_slug := range reminderData.TagsSlugs() {
			all_tags_slugs_with_emoji = append(all_tags_slugs_with_emoji, fmt.Sprintf("%v %v", utils.Symbols["tag"], tag_slug))
		}
		tag_index, _ := utils.AskOption(append(all_tags_slugs_with_emoji, fmt.Sprintf("%v %v", utils.Symbols["add"], "Add Tag")), "Select Tag")
		if tag_index != -1 {
			if tag_index == len(reminderData.TagsSlugs()) {
				// add new tag
				reminderData.NewTagRegistration()
			} else {
				tag := reminderData.Tags[tag_index]
				notes := reminderData.NotesWithTagId(tag.Id, "pending")
				err := reminderData.PrintNotesAndAskOptions(notes, tag.Id)
				utils.PrintErrorIfPresent(err)
			}
		}
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Note"):
		tag_ids := reminderData.AskTagIds([]int{})
		note := models.FNewNote(tag_ids)
		reminderData.NewNoteAppend(note)
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Add Tag"):
		reminderData.NewTagRegistration()
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Register Basic Tags"):
		reminderData.RegisterBasicTags()
	case fmt.Sprintf("%v %v", utils.Symbols["clock"], "Pending Notes"):
		all_notes := reminderData.Notes
		pending_notes := models.FNotesWithStatus(all_notes, "pending")
		var current_notes []*models.Note
		repeat_tag_ids := reminderData.TagIdsForGroup("repeat")
		// populating current_notes
		for _, note := range pending_notes {
			note_ids_with_repeat := utils.GetCommonIntMembers(note.TagIds, repeat_tag_ids)
			// first process notes without tag with group "repeat"
			// start showing such notes 7 days in advance from their due date, and until they are marked done
			min_day := note.CompleteBy - 7*24*60*60
			current_timestamp := utils.CurrentUnixTimestamp()
			if (len(note_ids_with_repeat) == 0) && (note.CompleteBy != 0) && (current_timestamp >= min_day) {
				current_notes = append(current_notes, note)
			}
			// check notes with tag with group "repeat"
			// start showing notes with "repeat-annually" 7 days in advance
			// start showing notes with "repeat-monthly" 3 days in advance
			// don't show such notes after their due date is past by 2 day
			if (len(note_ids_with_repeat) > 0) && (note.CompleteBy != 0) {
				// check for repeat-annually tag
				repeat_annually_tag := reminderData.TagFromSlug("repeat-annually")
				repeat_monthly_tag := reminderData.TagFromSlug("repeat-monthly")
				if (repeat_annually_tag != nil) && utils.IntInSlice(repeat_annually_tag.Id, note.TagIds) {
					_, note_month, note_day := utils.UnixTimestampToTime(note.CompleteBy).Date()
					note_timestamp_current := utils.UnixTimestampForCorrespondingCurrentYear(int(note_month), note_day)
					note_timestamp_previous := note_timestamp_current - 365*24*60*60
					note_timestamp_next := note_timestamp_current + 365*24*60*60
					days_before := int64(7)
					days_after := int64(2)
					if utils.IsTimeForRepeatNote(note_timestamp_current, note_timestamp_previous, note_timestamp_next, days_before, days_after) {
						current_notes = append(current_notes, note)
					}
				}
				// check for repeat-monthly tag
				if (repeat_monthly_tag != nil) && utils.IntInSlice(repeat_monthly_tag.Id, note.TagIds) {
					_, _, note_day := utils.UnixTimestampToTime(note.CompleteBy).Date()
					note_timestamp_current := utils.UnixTimestampForCorrespondingCurrentYearMonth(note_day)
					note_timestamp_previous := note_timestamp_current - 30*24*60*60
					note_timestamp_next := note_timestamp_current + 30*24*60*60
					days_before := int64(3)
					days_after := int64(2)
					if utils.IsTimeForRepeatNote(note_timestamp_current, note_timestamp_previous, note_timestamp_next, days_before, days_after) {
						current_notes = append(current_notes, note)
					}
				}
			}
		}
		fmt.Println("Note: Following are the pending notes with due date:")
		fmt.Println("  - within a week or already crossed (for non repeat-annually or repeat-monthly)")
		fmt.Println("  - within a week for repeat-annually and 2 days post due date (ignoring its year)")
		fmt.Println("  - within 3 days for repeat-monthly and 2 days post due date (ignoring its year and month)")
		err := reminderData.PrintNotesAndAskOptions(current_notes, -1)
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v", utils.Symbols["done"], "Done Notes"):
		all_notes := reminderData.Notes
		done_notes := models.FNotesWithStatus(all_notes, "done")
		fmt.Printf("A total of %v notes marked as 'done':\n", len(done_notes))
		err := reminderData.PrintNotesAndAskOptions(done_notes, -1)
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v", utils.Symbols["search"], "Search Notes"):
		// get texts of all notes
		all_notes := reminderData.Notes
		var all_texts []string
		for _, note := range all_notes {
			all_texts = append(all_texts, note.SearchText())
		}
		// function to search across notes
		search_notes := func(input string, idx int) bool {
			input = strings.ToLower(input)
			note_text := all_texts[idx]
			if strings.Contains(strings.ToLower(note_text), input) {
				return true
			}
			return false
		}
		// display prompt
		prompt := promptui.Select{
			Label:             "Notes",
			Items:             all_texts,
			Size:              25,
			StartInSearchMode: true,
			Searcher:          search_notes,
		}
		fmt.Printf("Searching through a total of %v notes:\n", len(all_texts))
		index, _, err := prompt.Run()
		utils.PrintErrorIfPresent(err)
		if index >= 0 {
			note := all_notes[index]
			reminderData.PrintNoteAndAskOptions(note)
		}
	case fmt.Sprintf("%v %v", utils.Symbols["pad"], "Create Backup"):
		// get backup file name
		ext := path.Ext(DataFile)
		dstFile := DataFile[:len(DataFile)-len(ext)] + "_backup_" + strconv.Itoa(int(utils.CurrentUnixTimestamp())) + ext
		lnFile := DataFile[:len(DataFile)-len(ext)] + "_backup_latest" + ext
		fmt.Printf("Creating backup at %q\n", dstFile)
		fmt.Println(dstFile)
		// create backup
		byteValue, err := ioutil.ReadFile(DataFile)
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
		fmt.Printf("Printing contents (and if possible, its difference since last backup) of %q:\n", DataFile)
		ext := path.Ext(DataFile)
		lnFile := DataFile[:len(DataFile)-len(ext)] + "_backup_latest" + ext
		err := utils.PerformWhich("wdiff")
		if err == nil {
			err = utils.PerformFilePresence(lnFile)
			if err == nil {
				err = utils.FPerformCwdiff(lnFile, DataFile)
			} else {
				fmt.Printf("Warning: `%v` file is not available yet\n", lnFile)
				err = utils.PerformCat(DataFile)
			}
		} else {
			fmt.Printf("%v Warning: `wdiff` command is not available\n", utils.Symbols["error"])
			err = utils.PerformCat(DataFile)
		}
		utils.PrintErrorIfPresent(err)
	case fmt.Sprintf("%v %v %v", utils.Symbols["checkerd_flag"], "Exit", utils.Symbols["red_flag"]):
		fmt.Println("Exiting...")
		return
	}
	flow()
}

func main() {
	go utils.Spinner(100 * time.Millisecond)
	flow()
}
