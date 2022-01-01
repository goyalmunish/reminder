package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/manifoldco/promptui"
)

// location info for `time`
// it can be set to update behavior of UnixTimestampToTime
var Location *time.Location

// get current time
func CurrentTime() time.Time {
	return time.Now()
}

// get current unix timestamp
func CurrentUnixTimestamp() int64 {
	return int64(CurrentTime().Unix())
}

// return UTC location
func UTCLocation() *time.Location {
	location, _ := time.LoadLocation("UTC")
	return location
}

// convert unix timestamp to time
// serve as central place to switch between UTC and local time
// by default use local time, but behavior can be changed via `Location`
func UnixTimestampToTime(unixTimestamp int64) time.Time {
	t := time.Unix(unixTimestamp, 0)
	if Location == nil {
		return t
	}
	return t.In(Location)
}

// convert unix timestamp to time string
func UnixTimestampToTimeStr(unixTimestamp int64, timeFormat string) string {
	var timeAsStr string
	if unixTimestamp > 0 {
		timeAsStr = UnixTimestampToTime(unixTimestamp).Format(timeFormat)
	} else {
		timeAsStr = "nil"
	}
	return timeAsStr
}

// convert unix timestamp to long time string
func UnixTimestampToLongTimeStr(unixTimestamp int64) string {
	return UnixTimestampToTimeStr(unixTimestamp, time.RFC850)
}

// convert unix timestamp to medium time string
func UnixTimestampToMediumTimeStr(unixTimestamp int64) string {
	return UnixTimestampToTimeStr(unixTimestamp, "02-Jan-06 15:04:05")
}

// convert unix timestamp to short time string
func UnixTimestampToShortTimeStr(unixTimestamp int64) string {
	return UnixTimestampToTimeStr(unixTimestamp, "02-Jan-06")
}

// get unix timestamp for date corresponding to current year
func UnixTimestampForCorrespondingCurrentYear(month int, day int) int64 {
	currentYear, _, _ := CurrentTime().Date()
	format := "2006-1-2"
	timeValue, _ := time.Parse(format, fmt.Sprintf("%v-%v-%v", currentYear, month, day))
	return int64(timeValue.Unix())
}

// get unix timestamp for date corresponding to current year and current month
func UnixTimestampForCorrespondingCurrentYearMonth(day int) int64 {
	currentYear, currentMonth, _ := CurrentTime().Date()
	format := "2006-1-2"
	timeValue, _ := time.Parse(format, fmt.Sprintf("%v-%v-%v", currentYear, int(currentMonth), day))
	return int64(timeValue.Unix())
}

// membership test for integer based array
func IntPresentInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// get common elements of two integer based slices
func GetCommonMembersIntSlices(arr1 []int, arr2 []int) []int {
	var arr []int
	for _, e1 := range arr1 {
		for _, e2 := range arr2 {
			if e1 == e2 {
				arr = append(arr, e1)
			}
		}
	}
	return arr
}

// print error if present
func PrintErrorIfPresent(err error) {
	if err != nil {
		fmt.Printf("%v %v\n", Symbols["error"], err)
	}
}

// return a trimmed string (with spaces removed from ends)
func TrimString(str string) string {
	return strings.TrimSpace(str)
}

// return a chopped strings (to a desired length)
func ChopStrings(texts []string, length int) []string {
	// return original texts (actually copy of what was passed)
	// if length value is not positive (considert ".." at the end)
	if length <= 2 {
		return texts
	}
	choppedStrings := []string{}
	for _, str := range texts {
		if len(str) <= length {
			choppedStrings = append(choppedStrings, str)
		} else {
			choppedStrings = append(choppedStrings, str[0:length-2]+"..")
		}
	}
	return choppedStrings
}

// validate that input is string
func ValidateString(input string) error {
	return nil
}

// validate that input is non-empty string
func ValidateNonEmptyString(input string) error {
	input = TrimString(input)
	if len(input) > 0 {
		return nil
	} else {
		return errors.New("Empty input")
	}
}

// validate date string (DD-MM-YYYY)
// nil is also valid input
func ValidateDateString(input string) error {
	input = TrimString(input)
	re := regexp.MustCompile("^((0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-((19|20)\\d\\d)|(nil))$")
	if re.MatchString(input) {
		return nil
	} else {
		return errors.New("Invalid input")
	}
}

// run given go template with given data and function map, and return the result as string
// it is interesting to note that even though data is recieved as `interface{}`, the template
// is able to access those attributes without even having to perform type assertion to get
// the underneath concrete value; this is contrary to masking behavior of interfaces
func TemplateResult(reportTemplate string, funcMap template.FuncMap, data interface{}) string {
	/*
		Issue with this function: It uses bytes.Buffer and converts it to string, but at the moment
		bytes (that is, uint8) are not converted back to rune/string properly.
	*/
	// define report result (as bytes)
	var reportResult bytes.Buffer
	// define report
	report := template.Must(template.New("report").Funcs(template.FuncMap(funcMap)).Parse(reportTemplate))
	// execute report to populate `reportResult`
	err := report.Execute(&reportResult, data)
	if err != nil {
		return err.Error()
	} else {
		// return report data as string
		return reportResult.String()
	}
}

// display spinner
func Spinner(delay time.Duration) {
	for {
		for _, c := range `–\|/` {
			fmt.Printf("\r%c", c)
			time.Sleep(delay)
		}
	}
}

// helper function to make assertion that `go` and `want` are nearly equal
func AssertEqual(t *testing.T, got interface{}, want interface{}) {
	if reflect.DeepEqual(got, want) {
		t.Logf("Matched value (by deep equality): %v", want)
	} else if reflect.DeepEqual(fmt.Sprintf("%v", got), fmt.Sprintf("%v", want)) {
		t.Logf("Matched value (by string conversion): %v", want)
	} else {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

// function to determine if it is time to show a repeat-based note/task
// dependency: `CurrentUnixTimestamp`
func IsTimeForRepeatNote(noteTimestampCurrent, noteTimestampPrevious, noteTimestampNext, daysBefore, daysAfter int64) bool {
	// fmt.Printf("Timestamp Curr: %v %v\n", noteTimestampCurrent, UnixTimestampToTime(noteTimestampCurrent))
	// fmt.Printf("Timestamp Prev: %v %v\n", noteTimestampPrevious, UnixTimestampToTime(noteTimestampPrevious))
	// fmt.Printf("Timestamp Next: %v %v\n", noteTimestampNext, UnixTimestampToTime(noteTimestampNext))
	// fmt.Printf("Days before: %v\n", daysBefore)
	// fmt.Printf("Days after: %v\n", daysAfter)
	currentTimestamp := CurrentUnixTimestamp()
	daysSecs := int64(24 * 60 * 60)
	return ((currentTimestamp >= noteTimestampCurrent-daysBefore*daysSecs) && (currentTimestamp <= noteTimestampCurrent+daysAfter*daysSecs)) ||
		((currentTimestamp >= noteTimestampPrevious-daysBefore*daysSecs) && (currentTimestamp <= noteTimestampPrevious+daysAfter*daysSecs)) ||
		((currentTimestamp >= noteTimestampNext-daysBefore*daysSecs) && (currentTimestamp <= noteTimestampNext+daysAfter*daysSecs))
}

// ask option to the user
// print error, if encountered any (so that they don't have to printed by calling function)
// return a tuple (chosen_index, chosen_string, error_if_any)
func AskOption(options []string, label string) (int, string, error) {
	if len(options) == 0 {
		err := errors.New("Empty List")
		fmt.Printf("%v Prompt failed %v\n", Symbols["warning"], err)
		return -1, "", err
	}
	// note: any item in options should not have \n character
	// otherwise such item is observed to not getting appear
	// in the rendered list
	prompt := promptui.Select{
		Label:  label,
		Items:  options,
		Size:   25,
		Stdout: &bellSkipper{},
	}
	index, result, err := prompt.Run()
	if err != nil {
		// error can happen if user raises an interrupt (such as Ctrl-c, SIGINT)
		fmt.Printf("%v Prompt failed %v\n", Symbols["warning"], err)
		return -1, "", err
	}
	fmt.Printf("You chose %d:%q\n", index, result)
	return index, result, nil
}

// perform shell operation and return its output
// note: it is better to avoid such functions
func PerformShellOperation(exe string, args ...string) (string, error) {
	executable, _ := exec.LookPath(exe)
	cmd := &exec.Cmd{
		Path:  executable,
		Args:  append([]string{executable}, args...),
		Stdin: os.Stdin,
	}
	bytes, err := cmd.Output()
	return string(bytes), err
}

// get terminal size
func TerminalSize() (int, int) {
	out, err := PerformShellOperation("stty", "size")
	if err != nil {
		log.Fatal(err)
	}
	output := strings.TrimSpace(string(out))
	dims := strings.Split(output, " ")
	height, _ := strconv.Atoi(dims[0])
	width, _ := strconv.Atoi(dims[1])
	return height, width
}

// get terminal width
func TerminalWidth() int {
	_, width := TerminalSize()
	return width
}

// check presence of a file
func PerformFilePresence(filePath string) error {
	output, err := PerformShellOperation("test", "-f", filePath)
	fmt.Println(output)
	return err
}

// check if a shell command is available
func PerformWhich(shellCmd string) error {
	output, err := PerformShellOperation("which", shellCmd)
	fmt.Println(output)
	return err
}

// cat a file
func PerformCat(filePath string) error {
	output, err := PerformShellOperation("cat", filePath)
	fmt.Println(output)
	return err
}

// get colored wdiff between two files
func FPerformCwdiff(oldFilePath string, newFilePath string) error {
	output, err := PerformShellOperation("wdiff", "-n", "-w", "\033[30;41m", "-x", "\033[0m", "-y", "\033[30;42m", "-z", "\033[0m", oldFilePath, newFilePath)
	fmt.Println(output)
	return err
}

// generate promptui.Prompt
func GeneratePrompt(promptName string, defaultText string) *promptui.Prompt {
	var prompt *promptui.Prompt
	switch promptName {
	case "user_name":
		prompt = &promptui.Prompt{
			Label:    "User Name",
			Default:  defaultText,
			Validate: ValidateString,
		}
	case "user_email":
		prompt = &promptui.Prompt{
			Label:    "User Email",
			Default:  defaultText,
			Validate: ValidateString,
		}
	case "tag_slug":
		prompt = &promptui.Prompt{
			Label:    "Tag Slug",
			Default:  defaultText,
			Validate: ValidateNonEmptyString,
		}
	case "tag_group":
		prompt = &promptui.Prompt{
			Label:    "Tag Group",
			Default:  defaultText,
			Validate: ValidateString,
		}
	case "tag_another":
		prompt = &promptui.Prompt{
			Label:    "Add another tag: yes/no (default: no):",
			Default:  defaultText,
			Validate: ValidateString,
		}
	case "note_text":
		prompt = &promptui.Prompt{
			Label:    "Note Text",
			Default:  defaultText,
			Validate: ValidateNonEmptyString,
		}
	case "note_comment":
		prompt = &promptui.Prompt{
			Label:    "New Comment",
			Default:  defaultText,
			Validate: ValidateNonEmptyString,
		}
	case "note_completed_by":
		prompt = &promptui.Prompt{
			Label:    "Due Date (format: DD-MM-YYYY), or enter nil to clear existing value",
			Default:  defaultText,
			Validate: ValidateDateString,
		}
	}
	return prompt
}

// generate promptui.Select
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
