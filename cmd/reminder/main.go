/*
Tool `reminder` is a command-line (terminal) based interactive app for organizing tasks with minimal efforts.

Just run it as `go run ./cmd/reminder`
*/
package main

import (
	"fmt"
	"time"

	"reminder/internal/model"
	"reminder/pkg/utils"
)

// flow is recursive function for overall flow of interactivity
func flow() {
	// make sure DataFile exists
	defaultDataFilePath := model.DefaultDataFile()
	_ = model.MakeSureFileExists(defaultDataFilePath)
	// read and parse the existing data
	reminderData := *model.ReadDataFile(defaultDataFilePath)
	// print data stats
	fmt.Println(reminderData.Stats())
	// try automatic backup
	_, _ = reminderData.AutoBackup(24 * 60 * 60)
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
	var err error
	_, result, _ := utils.AskOption([]string{
		fmt.Sprintf("%v %v", utils.Symbols["spark"], "List Stuff"),
		fmt.Sprintf("%v %v %v", utils.Symbols["checkerdFlag"], "Exit", utils.Symbols["redFlag"]),
		fmt.Sprintf("%v %v", utils.Symbols["clock"], "Approaching Due Date"),
		fmt.Sprintf("%v %v", utils.Symbols["hat"], "Main Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["search"], "Search Notes"),
		fmt.Sprintf("%v %v", utils.Symbols["backup"], "Create Backup"),
		fmt.Sprintf("%v %v", utils.Symbols["pad"], "Display Data File"),
		fmt.Sprintf("%v %v", utils.Symbols["clip"], "Register Basic Tags")}, "Select Option")
	// operate on main options
	switch result {
	case fmt.Sprintf("%v %v", utils.Symbols["spark"], "List Stuff"):
		err = reminderData.ListTags()
	case fmt.Sprintf("%v %v", utils.Symbols["clip"], "Register Basic Tags"):
		err = reminderData.RegisterBasicTags()
	case fmt.Sprintf("%v %v", utils.Symbols["clock"], "Approaching Due Date"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, -1, "pending", false)
	case fmt.Sprintf("%v %v", utils.Symbols["hat"], "Main Notes"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, -1, "pending", true)
	case fmt.Sprintf("%v %v", utils.Symbols["search"], "Search Notes"):
		err = reminderData.SearchNotes()
	case fmt.Sprintf("%v %v", utils.Symbols["backup"], "Create Backup"):
		_, err = reminderData.CreateBackup()
	case fmt.Sprintf("%v %v", utils.Symbols["pad"], "Display Data File"):
		err = reminderData.DisplayDataFile()
	case fmt.Sprintf("%v %v %v", utils.Symbols["checkerdFlag"], "Exit", utils.Symbols["redFlag"]):
		fmt.Println("Exiting...")
		return
	}
	utils.PrintErrorIfPresent(err)
	flow()
}

func main() {
	go utils.Spinner(100 * time.Millisecond)
	flow()
}
