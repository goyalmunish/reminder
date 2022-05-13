/*
Package utils provides common utility functions that are not reminder specific
*/
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

	"github.com/AlecAivazis/survey/v2"
)

// Location variable provides location info for `time`.
// It can be set to update behavior of UnixTimestampToTime.
var Location *time.Location

// CurrentTime function gets current time.
func CurrentTime() time.Time {
	return time.Now()
}

// CurrentUnixTimestamp function gets current unix timestamp.
func CurrentUnixTimestamp() int64 {
	return int64(CurrentTime().Unix())
}

// UTCLocation function returns UTC location.
func UTCLocation() *time.Location {
	location, _ := time.LoadLocation("UTC")
	return location
}

// UnixTimestampToTime function converts unix timestamp to time.
// It serves as central place to switch between UTC and local time.
// by default use local time, but behavior can be changed via `Location`.
func UnixTimestampToTime(unixTimestamp int64) time.Time {
	t := time.Unix(unixTimestamp, 0)
	if Location == nil {
		return t
	}
	return t.In(Location)
}

// UnixTimestampToTimeStr function converts unix timestamp to time string.
func UnixTimestampToTimeStr(unixTimestamp int64, timeFormat string) string {
	var timeAsStr string
	if unixTimestamp > 0 {
		timeAsStr = UnixTimestampToTime(unixTimestamp).Format(timeFormat)
	} else {
		timeAsStr = "nil"
	}
	return timeAsStr
}

// UnixTimestampToLongTimeStr function converts unix timestamp to long time string.
func UnixTimestampToLongTimeStr(unixTimestamp int64) string {
	return UnixTimestampToTimeStr(unixTimestamp, time.RFC850)
}

// UnixTimestampToMediumTimeStr function converts unix timestamp to medium time string.
func UnixTimestampToMediumTimeStr(unixTimestamp int64) string {
	return UnixTimestampToTimeStr(unixTimestamp, "02-Jan-06 15:04:05")
}

// UnixTimestampToShortTimeStr function converts unix timestamp to short time string.
func UnixTimestampToShortTimeStr(unixTimestamp int64) string {
	return UnixTimestampToTimeStr(unixTimestamp, "02-Jan-06")
}

// UnixTimestampForCorrespondingCurrentYear function gets unix timestamp for date corresponding to current year.
func UnixTimestampForCorrespondingCurrentYear(month int, day int) int64 {
	currentYear, _, _ := CurrentTime().Date()
	format := "2006-1-2"
	timeValue, _ := time.Parse(format, fmt.Sprintf("%v-%v-%v", currentYear, month, day))
	return int64(timeValue.Unix())
}

// UnixTimestampForCorrespondingCurrentYearMonth function gets unix timestamp for date corresponding to current year and current month.
func UnixTimestampForCorrespondingCurrentYearMonth(day int) int64 {
	currentYear, currentMonth, _ := CurrentTime().Date()
	format := "2006-1-2"
	timeValue, _ := time.Parse(format, fmt.Sprintf("%v-%v-%v", currentYear, int(currentMonth), day))
	return int64(timeValue.Unix())
}

// YearForDueDateDDMM return the current year if DD-MM is falling after current date, otherwise returns next year
func YearForDueDateDDMM(dateMonth string) (int, error) {
	format := "2-1-2006"
	currentTime := CurrentTime()
	// set current year as year if year part is missing
	timeSplit := strings.Split(dateMonth, "-")
	if len(timeSplit) != 2 {
		return 0, fmt.Errorf("Provided dateMonth string, %s, is not in DD-MM format", dateMonth)
	}
	// test with current year
	year := currentTime.Year()
	dateString := fmt.Sprintf("%s-%d", dateMonth, year)
	testTimeValue, err := time.Parse(format, dateString)
	if err != nil {
		return 0, err
	}
	if testTimeValue.Unix() <= currentTime.Unix() {
		// the due date falls before current date in current year
		// so, select next year instead
		year += 1
	}
	return year, nil
}

// IntPresentInSlice function performs membership test for integer based array.
func IntPresentInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetCommonMembersIntSlices function gets common elements of two integer based slices.
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

// PrintErrorIfPresent function prints error if present.
func PrintErrorIfPresent(err error) {
	if err != nil {
		fmt.Printf("%v %v\n", Symbols["error"], err)
	}
}

// TrimString function returns a trimmed string (with spaces removed from ends).
func TrimString(str string) string {
	return strings.TrimSpace(str)
}

// ChopStrings function returns a chopped strings (to a desired length).
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

// ValidateDateString function validates date string (DD-MM-YYYY) or (DD-MM).
// nil is also valid input
func ValidateDateString() survey.Validator {
	// return a validator that checks the length of the string
	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			// if the string is shorter than the given value
			input := TrimString(str)
			re := regexp.MustCompile("^((0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])(-((19|20)\\d\\d))?|(nil))$")
			if re.MatchString(input) {
				return nil
			} else {
				return fmt.Errorf("The input must be in the format DD-MM-YYYY or DD-MM.")
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("Invalid type %v", reflect.TypeOf(val).Name())
		}
	}
}

// TemplateResult function runs given go template with given data and function map, and return the result as string.
// It is interesting to note that even though data is recieved as `interface{}`, the template
// is able to access those attributes without even having to perform type assertion to get
// the underneath concrete value; this is contrary to masking behavior of interfaces.
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

// Spinner function displays spinner.
func Spinner(delay time.Duration) {
	for {
		for _, c := range `–\|/` {
			fmt.Printf("\r%c", c)
			time.Sleep(delay)
		}
	}
}

// AssertEqual function makes assertion that `go` and `want` are nearly equal.
func AssertEqual(t *testing.T, got interface{}, want interface{}) {
	if reflect.DeepEqual(got, want) {
		t.Logf("Matched value (by deep equality): %v", want)
	} else if reflect.DeepEqual(fmt.Sprintf("%v", got), fmt.Sprintf("%v", want)) {
		t.Logf("Matched value (by string conversion): %v", want)
	} else {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

// IsTimeForRepeatNote function determines if it is time to show a repeat-based note/task.
// dependency: `CurrentUnixTimestamp`
// It accepts current, previous and next timestamp of a task, and
// checks to see if any of the current timestamp falls in between [TIMESTAMP - DaysBefore, TIMESTAMP + DaysAfter]
func IsTimeForRepeatNote(noteTimestampCurrent, noteTimestampPrevious, noteTimestampNext, daysBefore, daysAfter int64) (bool, int64) {
	// fmt.Printf("Timestamp Curr: %v %v\n", noteTimestampCurrent, UnixTimestampToTime(noteTimestampCurrent))
	// fmt.Printf("Timestamp Prev: %v %v\n", noteTimestampPrevious, UnixTimestampToTime(noteTimestampPrevious))
	// fmt.Printf("Timestamp Next: %v %v\n", noteTimestampNext, UnixTimestampToTime(noteTimestampNext))
	// fmt.Printf("Days before: %v\n", daysBefore)
	// fmt.Printf("Days after: %v\n", daysAfter)
	currentTimestamp := CurrentUnixTimestamp()
	daysSecs := int64(24 * 60 * 60)
	condCurr := ((currentTimestamp >= noteTimestampCurrent-daysBefore*daysSecs) && (currentTimestamp <= noteTimestampCurrent+daysAfter*daysSecs))
	condNext := ((currentTimestamp >= noteTimestampNext-daysBefore*daysSecs) && (currentTimestamp <= noteTimestampNext+daysAfter*daysSecs))
	condPrev := ((currentTimestamp >= noteTimestampPrevious-daysBefore*daysSecs) && (currentTimestamp <= noteTimestampPrevious+daysAfter*daysSecs))
	// find which timestamp matched
	matchingTimestamp := noteTimestampPrevious
	if condCurr {
		matchingTimestamp = noteTimestampCurrent
	} else if condNext {
		matchingTimestamp = noteTimestampNext
	}
	return condCurr || condNext || condPrev, matchingTimestamp
}

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
	fmt.Printf("You chose %d:%q\n", selectedIndex, options[selectedIndex])
	return selectedIndex, options[selectedIndex], nil
}

// PerformShellOperation function performs shell operation and return its output.
// Note: it is better to avoid such functions.
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

// TerminalSize function gets terminal size.
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

// TerminalWidth function gets terminal width.
func TerminalWidth() int {
	_, width := TerminalSize()
	return width
}

// PerformFilePresence function checks presence of a file.
func PerformFilePresence(filePath string) error {
	output, err := PerformShellOperation("test", "-f", filePath)
	fmt.Println(output)
	return err
}

// PerformWhich function checks if a shell command is available.
func PerformWhich(shellCmd string) error {
	output, err := PerformShellOperation("which", shellCmd)
	fmt.Println(output)
	return err
}

// PerformCat function cats a file.
func PerformCat(filePath string) error {
	output, err := PerformShellOperation("cat", filePath)
	fmt.Println(output)
	return err
}

// PerformCwdiff function gets colored wdiff between two files.
func PerformCwdiff(oldFilePath string, newFilePath string) error {
	output, err := PerformShellOperation("wdiff", "-n", "-w", "\033[30;41m", "-x", "\033[0m", "-y", "\033[30;42m", "-z", "\033[0m", oldFilePath, newFilePath)
	fmt.Println(output)
	return err
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
