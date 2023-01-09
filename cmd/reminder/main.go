/*
Tool `reminder` is a command-line (terminal) based interactive app for organizing tasks with minimal efforts.

Just run it as `go run ./cmd/reminder`
*/
package reminder

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/goyalmunish/reminder/internal/model"
	"github.com/goyalmunish/reminder/internal/settings"
	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"
)

// flow is recursive function for overall flow of interactivity
var config *settings.Settings

func Run() error {
	// initialization
	var err error
	var runID = uuid.New()
	var startInteractiveProcess bool = true
	// note: setting are loaded before logger is being setup; it will assume only default logrus settings
	config, err = settings.LoadConfig()
	if err != nil {
		return err
	}
	logger.SetWithOptions(config.Log)
	logger.SetGlobalFields(map[string]interface{}{
		"app":    "reminder",
		"run_id": runID,
	})

	// make sure DataFile exists
	if err := model.MakeSureFileExists(config.AppInfo.DataFile, true); err != nil {
		return err
	}

	// read and parse the existing data
	reminderData, err := model.ReadDataFile(config.AppInfo.DataFile, false)
	if err != nil {
		return err
	}

	// check if the data file is locked by another session
	if reminderData.MutexLock {
		fmt.Printf("WARNING! %s\n", model.ErrorMutexLockOn.Error())
		reset := utils.AskBoolean("But, do you want to force reset the lock?")
		if !reset {
			// exit now without resetting the lock
			return model.ErrorMutexLockOn
		}
		// proceed forward to reset the lock, but don't start the interactive process
		startInteractiveProcess = false
	}

	// make sure lock is released and any uncommitted data is persisted
	defer func() {
		if !reminderData.MutexLock {
			// prevent multiple runs of cleanup
			return
		}
		// mutex lock is enable
		reminderData.MutexLock = false
		if err := reminderData.UpdateDataFile("Turning OFF the Mutex Lock, persisting the data, and closing the app!"); err != nil {
			utils.LogError(err)
		}
	}()

	// early exit if conditions are not met
	if !startInteractiveProcess {
		return model.ErrorInteractiveProcessSkipped
	}

	reminderData.MutexLock = true
	if err := reminderData.UpdateDataFile("Turning ON the Mutex Lock!"); err != nil {
		return err
	}

	// start the repeating interactive process
	if err := RepeatInteractiveSession(reminderData); err != nil {
		return err
	}
	return nil
}

func RepeatInteractiveSession(reminderData *model.ReminderData) error {
	var err error
	// print data stats
	stats, err := reminderData.Stats()
	fmt.Println(stats)
	utils.LogError(err)
	// try automatic backup
	_, err = reminderData.AutoBackup(24 * 60 * 60)
	utils.LogError(err)
	// ask the main menu
	fmt.Println("| =========================== MAIN MENU =========================== |")
	fmt.Println("|     Use 'Ctrl-c' to jump one level up (towards the Main Menu)     |")
	fmt.Println("| ----------------------------------------------------------------- |")
	/*
		How 'Ctrl-c' works?
		Hitting 'Ctrl-c' in golang raises as SIGINT signal.
		By default SIGINT signal (https://pkg.go.dev/os/signal) is converted to run-time panic,
		and eventually causes the program to exit.
		But, if you are inside PromptUI's `RepeatInteractiveSession()`, then it cancels the input and moves to next
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
		fmt.Sprintf("%s %s", utils.Symbols["refresh"], "Google Cloud Sync"),
		fmt.Sprintf("%s %s", utils.Symbols["pad"], "Display Data File")}, "Select Option")
	// operate on main options
	switch result {
	case fmt.Sprintf("%s %s", utils.Symbols["spark"], "List Stuff"):
		err = reminderData.ListTags()
	case fmt.Sprintf("%s %s", utils.Symbols["clock"], "Approaching Due Date"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, "pending_approaching_notes", -1, "due-date")
	case fmt.Sprintf("%s %s", utils.Symbols["hat"], "Main Notes"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, "pending_only_main_notes", -1, "default")
	case fmt.Sprintf("%s %s", utils.Symbols["search"], "Search Notes"):
		err = reminderData.SearchNotes()
	case fmt.Sprintf("%s %s", utils.Symbols["backup"], "Create Backup"):
		_, err = reminderData.CreateBackup()
	case fmt.Sprintf("%s %s", utils.Symbols["zzz"], "Suspended Notes"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, "suspended_notes", -1, "default")
	case fmt.Sprintf("%s %s", utils.Symbols["telescope"], "Look Ahead"):
		err = reminderData.PrintNotesAndAskOptions(model.Notes{}, "pending_long_view_notes", -1, "due-date")
	case fmt.Sprintf("%s %s", utils.Symbols["refresh"], "Google Cloud Sync"):
		err = reminderData.SyncCalendar(config.Calendar)
	case fmt.Sprintf("%s %s", utils.Symbols["pad"], "Display Data File"):
		err = reminderData.DisplayDataFile()
	case fmt.Sprintf("%s %s %s", utils.Symbols["checkerdFlag"], "Exit", utils.Symbols["redFlag"]):
		fmt.Println("Exiting...")
		return nil
	}
	// ignore the error, but log it
	utils.LogError(err)
	// keep the interactive session running
	return RepeatInteractiveSession(reminderData)
}
