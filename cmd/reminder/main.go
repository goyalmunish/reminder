/*
Tool `reminder` is a command-line (terminal) based interactive app for organizing tasks with minimal efforts.

Just run it as `go run ./cmd/reminder`
*/
package main

import (
	"fmt"

	"github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/pkg/utils"
)

// flow is recursive function for overall flow of interactivity
func flow() {
	var err error
	// make sure DataFile exists
	defaultDataFilePath := model.DefaultDataFile()
	err = model.MakeSureFileExists(defaultDataFilePath)
	utils.PrintError(err)
	// read and parse the existing data
	reminderData := *model.ReadDataFile(defaultDataFilePath)
	// print data stats
	fmt.Println(reminderData.Stats())
	// try automatic backup
	_, err = reminderData.AutoBackup(24 * 60 * 60)
	utils.PrintError(err)
	// ask the main menu
	fmt.Println("| =========================== MAIN MENU =========================== |")
	fmt.Println("|     Use 'Ctrl-c' to jump one level up (towards the Main Menu)     |")
	fmt.Println("| ----------------------------------------------------------------- |")
	/*
		How 'Ctrl-c' works?
		Hitting 'Ctrl-c' in golang raises as SIGINT signal.
		By default SIGINT signal (https://pkg.go.dev/os/signal) is converted to run-time panic,
		and eventually causes the program to exit.
		But, if you are inside PromptUI's `Run()`, then it cancels the input and moves to next
		statement in the code.
	*/
	_, result, _ := utils.AskOption([]string{
		fmt.Sprintf("%s %s", utils.Symbols["spark"], "List Stuff"),
		fmt.Sprintf("%s %s %s", utils.Symbols["checkerdFlag"], "Exit", utils.Symbols["redFlag"]),
		fmt.Sprintf("%s %s", utils.Symbols["clock"], "Approaching Due Date"),
		fmt.Sprintf("%s %s", utils.Symbols["hat"], "Main Notes"),
		fmt.Sprintf("%s %s", utils.Symbols["search"], "Search Notes"),
		fmt.Sprintf("%s %s", utils.Symbols["backup"], "Create Backup"),
		fmt.Sprintf("%s %s", utils.Symbols["zzz"], "Suspended Notes"),
		fmt.Sprintf("%s %s", utils.Symbols["telescope"], "Look Ahead"),
		fmt.Sprintf("%s %s", utils.Symbols["pad"], "Display Data File"),
		fmt.Sprintf("%s %s", utils.Symbols["clip"], "Register Basic Tags")}, "Select Option")
	// operate on main options
	switch result {
	case fmt.Sprintf("%s %s", utils.Symbols["spark"], "List Stuff"):
		err = reminderData.ListTags()
	case fmt.Sprintf("%s %s", utils.Symbols["clip"], "Register Basic Tags"):
		err = reminderData.RegisterBasicTags()
	case fmt.Sprintf("%s %s", utils.Symbols["clock"], "Approaching Due Date"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, -1, "pending", false, "default")
	case fmt.Sprintf("%s %s", utils.Symbols["hat"], "Main Notes"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, -1, "pending", true, "default")
	case fmt.Sprintf("%s %s", utils.Symbols["search"], "Search Notes"):
		err = reminderData.SearchNotes()
	case fmt.Sprintf("%s %s", utils.Symbols["backup"], "Create Backup"):
		_, err = reminderData.CreateBackup()
	case fmt.Sprintf("%s %s", utils.Symbols["zzz"], "Suspended Notes"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, -1, "suspended", true, "default")
	case fmt.Sprintf("%s %s", utils.Symbols["telescope"], "Look Ahead"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, -1, "pending", false, "due-date")
	case fmt.Sprintf("%s %s", utils.Symbols["pad"], "Display Data File"):
		err = reminderData.DisplayDataFile()
	case fmt.Sprintf("%s %s %s", utils.Symbols["checkerdFlag"], "Exit", utils.Symbols["redFlag"]):
		fmt.Println("Exiting...")
		return
	}
	utils.PrintError(err)
	flow()
}

func main() {
	// go utils.Spinner(100 * time.Millisecond)
	flow()
}
