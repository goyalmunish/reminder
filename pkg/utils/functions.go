/*
Package utils provides common utility functions that are not reminder specific
*/
package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/goyalmunish/reminder/pkg/logger"
)

// HomeDir return the home directory path for current user.
// Note: It is deliberately defined as a variable to make it
// easer to patch in tests.
var HomeDir func() string = func() string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return dir
}

// TerminalSize function gets terminal size.
// Note: It is deliberately defined as a variable to make it
// easer to patch in tests.
var TerminalSize = func() (int, int, error) {
	out, err := PerformShellOperation("stty", "size")
	if err != nil {
		return 0, 0, err
	}
	output := strings.TrimSpace(string(out))
	dims := strings.Split(output, " ")
	height, _ := strconv.Atoi(dims[0])
	width, _ := strconv.Atoi(dims[1])
	return height, width, nil
}

// IsMemberOf function performs membership test.
func IsMemberOfSlice[V comparable](a V, list []V) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetCommonMembersOfSlices function gets common elements of slices.
func GetCommonMembersOfSlices[V comparable](list1 []V, list2 []V) []V {
	var result []V
	for _, e1 := range list1 {
		for _, e2 := range list2 {
			if e1 == e2 {
				result = append(result, e1)
			}
		}
	}
	return result
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

// TemplateResult function runs given go template with given data and function map, and return the result as string.
// It is interesting to note that even though data is recieved as `interface{}`, the template
// is able to access those attributes without even having to perform type assertion to get
// the underneath concrete value; this is contrary to masking behavior of interfaces.
func TemplateResult(reportTemplate string, funcMap template.FuncMap, data interface{}) (string, error) {
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
		return "", err
	} else {
		// return report data as string
		return reportResult.String(), nil
	}
}

// LogError function ignores but prints the error (if present).
func LogError(err error) {
	if err != nil {
		logger.Error(fmt.Sprintf("%v %v\n", Symbols["error"], err))
	}
}

// AssertEqual function makes assertion that `go` and `want` are nearly equal.
func AssertEqual(t *testing.T, got interface{}, want interface{}) {
	if reflect.DeepEqual(got, want) {
		t.Logf("Pass: Matched value (by deep equality): WANT %v GOT %v", want, got)
	} else if reflect.DeepEqual(fmt.Sprintf("%v", got), fmt.Sprintf("%v", want)) {
		t.Logf("Pass: Matched value (by string conversion): WANT %v GOT %v", want, got)
	} else {
		var errorMsg = struct {
			gotType   string
			gotValue  interface{}
			wantType  string
			wantValue interface{}
		}{
			gotValue:  fmt.Sprintf("%v", got),
			gotType:   fmt.Sprintf("%T", got),
			wantValue: fmt.Sprintf("%v", want),
			wantType:  fmt.Sprintf("%T", want),
		}
		t.Errorf("Error: %+v", errorMsg)
	}
}

// TryConvertTildaBasedPath converts tilda (~) based path to complete path.
// For a non-tilda based path, return as it is.
func TryConvertTildaBasedPath(path string) string {
	homeDir := HomeDir()
	if path == "~" {
		path = homeDir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(homeDir, path[2:])
	}
	return path
}

// AskBoolean asks a boolean question to the user.
func AskBoolean(msg string) (bool, error) {
	return askBoolean(msg, os.Stdin)
}

func askBoolean(msg string, in io.Reader) (bool, error) {
	var res string
	fmt.Printf("%s (y/n): ", msg)
	_, err := fmt.Fscanln(in, &res)
	if err != nil {
		return false, err
	}
	logger.Info(fmt.Sprintf("Received response: %q", res))
	res = strings.Trim(res, " \n\t")
	res = strings.ToLower(res)
	return res == "y", nil
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

// TerminalWidth function gets terminal width.
func TerminalWidth() (int, error) {
	_, width, err := TerminalSize()
	if err != nil {
		return 0, err
	}
	return width, nil
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

// MatchedTimestamp function determines if currentTime matches with (within allowed daysBefore and daysAfter) given 3 timestamps (current, previous, and next).
// That is, it checks to see if any of the current timestamp falls in between [TIMESTAMP - DaysBefore, TIMESTAMP + DaysAfter]
func MatchedTimestamp(noteTimestampCurrent, noteTimestampPrevious, noteTimestampNext, daysBefore, daysAfter int64) (bool, int64) {
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

// Spinner function displays spinner.
func Spinner(delay time.Duration) {
	for {
		for _, c := range `–\|/` {
			fmt.Printf("\r%c", c)
			time.Sleep(delay)
		}
	}
}
