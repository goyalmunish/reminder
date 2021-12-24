/*
Test Process:
- `go test`: Test without supressing printing to console
- `go test .`: Test while supressing printing to console

// general form of each example
func TestFunctionName(t *testing.T) {
	utils.AssertEqual(t, got, want)
}
*/

package utils_test

import (
	// "fmt"
	"errors"
	"testing"
	"time"

	// "github.com/golang/mock/gomock"

	utils "reminder/pkg/utils"
)

func TestCurrentUnixTimestamp(t *testing.T) {
	want := time.Now().Unix()
	output := utils.CurrentUnixTimestamp()
	utils.AssertEqual(t, output, want)
}

func TestUnixTimestampToTime(t *testing.T) {
	currentTime := utils.CurrentTime()
	currentTimestamp := currentTime.Unix()
	output := utils.UnixTimestampToTime(currentTimestamp)
	utils.AssertEqual(t, output.Format(time.UnixDate), currentTime.Format(time.UnixDate))
}

func TestUnixTimestampToTimeStr(t *testing.T) {
	output := utils.UnixTimestampToTimeStr(int64(1608575176), "02-Jan-06")
	utils.AssertEqual(t, output, "21-Dec-20")
	output = utils.UnixTimestampToTimeStr(int64(1608575176), time.RFC850)
	utils.AssertEqual(t, output, "Monday, 21-Dec-20 18:26:16 UTC")
	output = utils.UnixTimestampToTimeStr(int64(-1), "02-Jan-06")
	utils.AssertEqual(t, output, "nil")
}

func TestUnixTimestampToShortTimeStr(t *testing.T) {
	output := utils.UnixTimestampToShortTimeStr(int64(1608575176))
	utils.AssertEqual(t, output, "21-Dec-20")
}

func TestUnixTimestampToMediumTimeStr(t *testing.T) {
	output := utils.UnixTimestampToMediumTimeStr(int64(1608575176))
	utils.AssertEqual(t, output, "21-Dec-20 18:26:16")
}

func TestUnixTimestampToLongTimeStr(t *testing.T) {
	output := utils.UnixTimestampToLongTimeStr(int64(1608575176))
	utils.AssertEqual(t, output, "Monday, 21-Dec-20 18:26:16 UTC")
}

func TestUnixTimestampForCorrespondingCurrentYear(t *testing.T) {
	got := utils.UnixTimestampForCorrespondingCurrentYear(9, 30) - utils.UnixTimestampForCorrespondingCurrentYear(6, 30)
	utils.AssertEqual(t, got, 7948800)
	got = utils.UnixTimestampForCorrespondingCurrentYear(10, 1) - utils.UnixTimestampForCorrespondingCurrentYear(7, 1)
	utils.AssertEqual(t, got, 7948800)
}

func TestUnixTimestampForCorrespondingCurrentYearMonth(t *testing.T) {
	got := utils.UnixTimestampForCorrespondingCurrentYearMonth(9) - utils.UnixTimestampForCorrespondingCurrentYearMonth(1)
	utils.AssertEqual(t, got, 691200)
	got = utils.UnixTimestampForCorrespondingCurrentYearMonth(28) - utils.UnixTimestampForCorrespondingCurrentYearMonth(1)
	utils.AssertEqual(t, got, 2332800)
}

func TestIntPresentInSlice(t *testing.T) {
	utils.AssertEqual(t, utils.IntPresentInSlice(100, []int{-100, 0, 100}), true)
	utils.AssertEqual(t, utils.IntPresentInSlice(-100, []int{-100, 0, 100}), true)
	utils.AssertEqual(t, utils.IntPresentInSlice(101, []int{-100, 0, 100}), false)
	utils.AssertEqual(t, utils.IntPresentInSlice(-99, []int{-100, 0, 100}), false)
}

func TestGetCommonMembersIntSlices(t *testing.T) {
	utils.AssertEqual(t,
		utils.GetCommonMembersIntSlices([]int{-100, 0, 100, 1, 10, 8, 2, -51, 4},
			[]int{-21, 100, 0, 8, 4}),
		[]int{0, 100, 8, 4})
	utils.AssertEqual(t,
		utils.GetCommonMembersIntSlices([]int{-21, 100, 0, 8, 4},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{100, 0, 8, 4})
	utils.AssertEqual(t,
		utils.GetCommonMembersIntSlices([]int{2},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{2})
	utils.AssertEqual(t,
		utils.GetCommonMembersIntSlices([]int{},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{})
	utils.AssertEqual(t,
		utils.GetCommonMembersIntSlices([]int{-100, 0, 100, 1, 10, 8, 2, -51, 4},
			[]int{}),
		[]int{})
}

func TestPrintErrorIfPresent(t *testing.T) {
	err := errors.New("dummy error")
	utils.PrintErrorIfPresent(err)
}

func TestTrimString(t *testing.T) {
	utils.AssertEqual(t, utils.TrimString("   str"), "str")
	utils.AssertEqual(t, utils.TrimString("str   "), "str")
	utils.AssertEqual(t, utils.TrimString("  str "), "str")
}

func TestChopStrings(t *testing.T) {
	strings := []string{"0123456789", "ABCDEFG", "0123"}
	utils.AssertEqual(t, utils.ChopStrings(strings, 2), strings)
	utils.AssertEqual(t, utils.ChopStrings(strings, 1), strings)
	utils.AssertEqual(t, utils.ChopStrings(strings, 0), strings)
	utils.AssertEqual(t, utils.ChopStrings(strings, -1), strings)
	want := []string{"012..", "ABC..", "0123"}
	utils.AssertEqual(t, utils.ChopStrings(strings, 5), want)
	want = []string{"0123456..", "ABCDEFG", "0123"}
	utils.AssertEqual(t, utils.ChopStrings(strings, 9), want)
}

func TestValidateString(t *testing.T) {
	utils.AssertEqual(t, utils.ValidateString("str"), nil)
	utils.AssertEqual(t, utils.ValidateString(""), nil)
}

func TestValidateNonEmptyString(t *testing.T) {
	utils.AssertEqual(t, utils.ValidateNonEmptyString("str"), nil)
	utils.AssertEqual(t, utils.ValidateNonEmptyString(""), errors.New("Empty input"))
}

func TestValidateDateString(t *testing.T) {
	utils.AssertEqual(t, utils.ValidateDateString("31-12-2020"), nil)
	utils.AssertEqual(t, utils.ValidateDateString("nil"), nil)
	utils.AssertEqual(t, utils.ValidateDateString("12-31-2020"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020-12-31"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020-31-"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020-31"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020-"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020"), errors.New("Invalid input"))
}

func TestTerminalSize(t *testing.T) {
	// perhaps stty command doesn't work in tests
	// height, width := utils.TerminalSize()
	// utils.AssertEqual(t, height > 0, true)
	// utils.AssertEqual(t, width > 0, true)
}

func TestPerformShellOperation(t *testing.T) {
	dummyFile := "dummyFile"
	defer utils.PerformShellOperation("rm -f", dummyFile)
	// attempt to delete a non-existing file
	_, err := utils.PerformShellOperation("rm", dummyFile)
	utils.AssertEqual(t, err, errors.New("exit status 1"))
	// create and delete a file
	_, err = utils.PerformShellOperation("touch", dummyFile)
	utils.AssertEqual(t, err, nil)
	_, err = utils.PerformShellOperation("ls", "-lhFa", dummyFile)
	utils.AssertEqual(t, err, nil)
	_, err = utils.PerformShellOperation("rm", dummyFile)
	utils.AssertEqual(t, err, nil)
	// attempt to invoke a command that do not exist
	_, err = utils.PerformShellOperation("command_do_not_exist")
	utils.AssertEqual(t, err, errors.New("fork/exec : no such file or directory"))
	_, err = utils.PerformShellOperation("command_do_not_exist", "arg1", "arg2")
	utils.AssertEqual(t, err, errors.New("fork/exec : no such file or directory"))
}
