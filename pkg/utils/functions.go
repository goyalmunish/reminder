package utils

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

// fmt.Sprintf("%v %v", Symbols["tag"], "Update tags"):

// type StandaloneFunctions interface {
// 	CurrentUnixTimestamp() int64
// 	UnixTimestampToTime(unix_timestamp int64) time.Time
// 	UnixTimestampToTimeStr(unix_timestamp int64, time_format string) string
// 	UnixTimestampToLongTimeStr(unix_timestamp int64) string
// 	UnixTimestampToShortTimeStr(unix_timestamp int64) string
// 	UnixTimestampForCorrespondingCurrentYear(month int, day int) int64
// 	UnixTimestampForCorrespondingCurrentYearMonth(day int) int64
// 	IntPresentInSlice(a int, list []int) bool
// 	GetCommonMembersIntSlices(arr1 []int, arr2 []int) []int
// 	PrintErrorIfPresent(err error)
// 	TrimString(str string) string
// 	ValidateString(input string) error
// 	ValidateNonEmptyString(input string) error
// 	ValidateDateString(input string) error
// 	PerformShellOperation(exe string, args ...string) error
// 	PerformFilePresence(file_path string) error
// 	PerformWhich(shell_cmd string) error
// 	PerformCat(file_path string) error
// 	FPerformCwdiff(old_file_path string, new_file_path string) error
// 	AskOption(options []string, label string) (int, string)
// 	IsTimeForRepeatNote(note_timestamp_current, note_timestamp_previous, note_timestamp_next, days_before, days_after int64) bool
// 	Spinner(delay time.Duration)
// 	AssertEqual(t *testing.T, got interface{}, want interface{})
// }

// get current time
// serve as a central place to switch between
// local time and UTC
func CurrentTime() time.Time {
	return time.Now().UTC()
}

// get current unix timestamp
func CurrentUnixTimestamp() int64 {
	return int64(CurrentTime().Unix())
}

// convert unix timestamp to time
func UnixTimestampToTime(unix_timestamp int64) time.Time {
	return time.Unix(unix_timestamp, 0).UTC()
}

// convert unix timestamp to time string
func UnixTimestampToTimeStr(unix_timestamp int64, time_format string) string {
	var time_as_str string
	if unix_timestamp > 0 {
		time_as_str = UnixTimestampToTime(unix_timestamp).Format(time_format)
	} else {
		time_as_str = "nil"
	}
	return time_as_str
}

// convert unix timestamp to long time string
func UnixTimestampToLongTimeStr(unix_timestamp int64) string {
	return UnixTimestampToTimeStr(unix_timestamp, time.RFC850)
}

// convert unix timestamp to short time string
func UnixTimestampToShortTimeStr(unix_timestamp int64) string {
	return UnixTimestampToTimeStr(unix_timestamp, "02-Jan-06")
}

// get unix timestamp for date corresponding to current year
func UnixTimestampForCorrespondingCurrentYear(month int, day int) int64 {
	current_year, _, _ := CurrentTime().Date()
	format := "2006-1-2"
	time_value, _ := time.Parse(format, fmt.Sprintf("%v-%v-%v", current_year, month, day))
	return int64(time_value.Unix())
}

// get unix timestamp for date corresponding to current year and current month
func UnixTimestampForCorrespondingCurrentYearMonth(day int) int64 {
	current_year, current_month, _ := CurrentTime().Date()
	format := "2006-1-2"
	time_value, _ := time.Parse(format, fmt.Sprintf("%v-%v-%v", current_year, int(current_month), day))
	return int64(time_value.Unix())
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

// validate date string (YYYY-MM-DD)
func ValidateDateString(input string) error {
	input = TrimString(input)
	re := regexp.MustCompile("^((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])$")
	if re.MatchString(input) {
		return nil
	} else {
		return errors.New("Invalid input")
	}
}

// display spinner
func Spinner(delay time.Duration) {
	for {
		for _, c := range `-\|/` {
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
func IsTimeForRepeatNote(note_timestamp_current, note_timestamp_previous, note_timestamp_next, days_before, days_after int64) bool {
	current_timestamp := CurrentUnixTimestamp()
	day_secs := int64(24 * 60 * 60)
	return ((current_timestamp >= note_timestamp_current-days_before*day_secs) && (current_timestamp <= note_timestamp_current+days_after*day_secs)) ||
		((current_timestamp >= note_timestamp_previous-days_before*day_secs) && (current_timestamp <= note_timestamp_previous+days_after*day_secs)) ||
		((current_timestamp >= note_timestamp_next-days_before*day_secs) && (current_timestamp <= note_timestamp_next+days_after*day_secs))
}

// ask option to the user
func AskOption(options []string, label string) (int, string) {
	if len(options) == 0 {
		fmt.Println("No results")
		return -1, "error"
	}
	// note: any item in options should not have \n character
	// otherwise such item is observed to not getting appear
	// in the rendered list
	prompt := promptui.Select{
		Label: label,
		Items: options,
		Size:  25,
	}
	index, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v Prompt failed %v\n", Symbols["error"], err)
		return -1, "error"
	}
	fmt.Printf("You chose %d:%q\n", index, result)
	return index, result
}

// perform shell operation
// note: it is better to avoid such functions
func PerformShellOperation(exe string, args ...string) error {
	executable, _ := exec.LookPath(exe)
	cmd := &exec.Cmd{
		Path:   executable,
		Args:   append([]string{executable}, args...),
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
	err := cmd.Run()
	return err
}

// check presence of a file
func PerformFilePresence(file_path string) error {
	return PerformShellOperation("test", "-f", file_path)
}

// check if a shell command is available
func PerformWhich(shell_cmd string) error {
	return PerformShellOperation("which", shell_cmd)
}

// cat a file
func PerformCat(file_path string) error {
	return PerformShellOperation("cat", file_path)
}

// get colored wdiff between two files
func FPerformCwdiff(old_file_path string, new_file_path string) error {
	return PerformShellOperation("wdiff", "-n", "-w", "\033[30;41m", "-x", "\033[0m", "-y", "\033[30;42m", "-z", "\033[0m", old_file_path, new_file_path)
}
