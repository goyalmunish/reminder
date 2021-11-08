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
	current_time := utils.CurrentTime()
	current_timestamp := current_time.Unix()
	output := utils.UnixTimestampToTime(current_timestamp)
	utils.AssertEqual(t, output.Format(time.UnixDate), current_time.Format(time.UnixDate))
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

func TestIntInSlice(t *testing.T) {
	utils.AssertEqual(t, utils.IntInSlice(100, []int{-100, 0, 100}), true)
	utils.AssertEqual(t, utils.IntInSlice(-100, []int{-100, 0, 100}), true)
	utils.AssertEqual(t, utils.IntInSlice(101, []int{-100, 0, 100}), false)
	utils.AssertEqual(t, utils.IntInSlice(-99, []int{-100, 0, 100}), false)
}

func TestGetCommonIntMembers(t *testing.T) {
	utils.AssertEqual(t,
		utils.GetCommonIntMembers([]int{-100, 0, 100, 1, 10, 8, 2, -51, 4},
			[]int{-21, 100, 0, 8, 4}),
		[]int{0, 100, 8, 4})
	utils.AssertEqual(t,
		utils.GetCommonIntMembers([]int{-21, 100, 0, 8, 4},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{100, 0, 8, 4})
	utils.AssertEqual(t,
		utils.GetCommonIntMembers([]int{2},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{2})
	utils.AssertEqual(t,
		utils.GetCommonIntMembers([]int{},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{})
	utils.AssertEqual(t,
		utils.GetCommonIntMembers([]int{-100, 0, 100, 1, 10, 8, 2, -51, 4},
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

func TestValidateString(t *testing.T) {
	utils.AssertEqual(t, utils.ValidateString("str"), nil)
	utils.AssertEqual(t, utils.ValidateString(""), nil)
}

func TestValidateNonEmptyString(t *testing.T) {
	utils.AssertEqual(t, utils.ValidateNonEmptyString("str"), nil)
	utils.AssertEqual(t, utils.ValidateNonEmptyString(""), errors.New("Empty input"))
}

func TestValidateDateString(t *testing.T) {
	utils.AssertEqual(t, utils.ValidateDateString("2020-12-31"), nil)
	utils.AssertEqual(t, utils.ValidateDateString("2020-31-12"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020-31-"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020-31"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020-"), errors.New("Invalid input"))
	utils.AssertEqual(t, utils.ValidateDateString("2020"), errors.New("Invalid input"))
}

func TestPerformShellOperation(t *testing.T) {
	dummy_file := "dummy_file"
	defer utils.PerformShellOperation("rm -f", dummy_file)
	// attempt to delete a non-existing file
	err := utils.PerformShellOperation("rm", dummy_file)
	utils.AssertEqual(t, err, errors.New("exit status 1"))
	// create and delete a file
	err = utils.PerformShellOperation("touch", dummy_file)
	utils.AssertEqual(t, err, nil)
	err = utils.PerformShellOperation("ls", "-lhFa", dummy_file)
	utils.AssertEqual(t, err, nil)
	err = utils.PerformShellOperation("rm", dummy_file)
	utils.AssertEqual(t, err, nil)
	// attempt to invoke a command that do not exist
	err = utils.PerformShellOperation("command_do_not_exist")
	utils.AssertEqual(t, err, errors.New("fork/exec : no such file or directory"))
	err = utils.PerformShellOperation("command_do_not_exist", "arg1", "arg2")
	utils.AssertEqual(t, err, errors.New("fork/exec : no such file or directory"))
}
