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
	"fmt"
	"html/template"
	"testing"
	"time"

	// "github.com/golang/mock/gomock"

	utils "github.com/goyalmunish/reminder/pkg/utils"
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
	utils.Location = utils.UTCLocation()
	output := utils.UnixTimestampToTimeStr(int64(1608575176), "02-Jan-06")
	utils.AssertEqual(t, output, "21-Dec-20")
	output = utils.UnixTimestampToTimeStr(int64(1608575176), time.RFC850)
	utils.AssertEqual(t, output, "Monday, 21-Dec-20 18:26:16 UTC")
	output = utils.UnixTimestampToTimeStr(int64(-1), "02-Jan-06")
	utils.AssertEqual(t, output, "nil")
	output = utils.UnixTimestampToTimeStr(int64(1698710400), time.RFC850)
	utils.AssertEqual(t, output, "Tuesday, 31-Oct-23 00:00:00 UTC")
}

func TestUnixTimestampToLongTimeStr(t *testing.T) {
	utils.Location = utils.UTCLocation()
	output := utils.UnixTimestampToLongTimeStr(int64(1608575176))
	utils.AssertEqual(t, output, "Monday, 21-Dec-20 18:26:16 UTC")
}

func TestUnixTimestampToMediumTimeStr(t *testing.T) {
	utils.Location = utils.UTCLocation()
	output := utils.UnixTimestampToMediumTimeStr(int64(1608575176))
	utils.AssertEqual(t, output, "21-Dec-20 18:26:16")
}

func TestUnixTimestampToShortTimeStr(t *testing.T) {
	utils.Location = utils.UTCLocation()
	output := utils.UnixTimestampToShortTimeStr(int64(1608575176))
	utils.AssertEqual(t, output, "21-Dec-20")
}

func TestUnixTimestampForCorrespondingCurrentYear(t *testing.T) {
	utils.Location = utils.UTCLocation()
	got := utils.UnixTimestampForCorrespondingCurrentYear(9, 30) - utils.UnixTimestampForCorrespondingCurrentYear(6, 30)
	utils.AssertEqual(t, got, 7948800)
	got = utils.UnixTimestampForCorrespondingCurrentYear(10, 1) - utils.UnixTimestampForCorrespondingCurrentYear(7, 1)
	utils.AssertEqual(t, got, 7948800)
}

func TestUnixTimestampForCorrespondingCurrentYearMonth(t *testing.T) {
	utils.Location = utils.UTCLocation()
	got := utils.UnixTimestampForCorrespondingCurrentYearMonth(9) - utils.UnixTimestampForCorrespondingCurrentYearMonth(1)
	utils.AssertEqual(t, got, 691200)
	got = utils.UnixTimestampForCorrespondingCurrentYearMonth(28) - utils.UnixTimestampForCorrespondingCurrentYearMonth(1)
	utils.AssertEqual(t, got, 2332800)
}

func TestStrToTime(t *testing.T) {
	// note: refer this format to quickly write table based tests
	var tests = []struct {
		name      string
		input     string // RFC3339 string
		want      int64  // Unix64 time value
		wantedErr bool   // whether an error was expected
	}{
		{name: "time in GMT", input: "2022-12-28T00:18:18.929Z", want: 1672186698},
		{name: "time in Melbourne/Australia", input: "2022-12-28T11:18:18.929+11:00", want: 1672186698},
	}
	for position, subtest := range tests {
		got, err := utils.StrToTime(subtest.input, "")
		if (err != nil) != subtest.wantedErr {
			t.Fatalf("StrToTime case %q (position=%d) failed for input %q with error %q", subtest.name, position, subtest.input, err)
		}
		if got.Unix() != subtest.want {
			t.Errorf("StrToTime case %q (position=%d) failed for input %q; returns <%+v>; wants <%+v>", subtest.name, position, subtest.input, got.Unix(), subtest.want)
		}
	}
}

func TestTimeToStr(t *testing.T) {
	var tests = []struct {
		name  string
		input string // RFC3339 string
		want  string // RFC3339 string
	}{
		{name: "time in GMT", input: "2022-12-28T00:18:18.929Z", want: "2022-12-28T00:18:18Z"},
		{name: "time in Melbourne/Australia", input: "2022-12-28T11:18:18.929+11:00", want: "2022-12-28T11:18:18+11:00"},
	}
	for position, subtest := range tests {
		tm, err := time.Parse(time.RFC3339, subtest.input)
		if err != nil {
			t.Fatalf("Test input %q is incorrect", subtest.input)
		}
		got := utils.TimeToStr(tm)
		if got != subtest.want {
			t.Errorf("TimeToStr case %q (position=%d) failed for input %q; returns <%+v>; wants <%+v>", subtest.name, position, tm, got, subtest.want)
		}
	}
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

func TestLogError(t *testing.T) {
	err := errors.New("dummy error")
	utils.LogError(err)
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

func TestValidateDateString(t *testing.T) {
	errorMsg := "The input must be in the format DD-MM-YYYY or DD-MM."
	utils.AssertEqual(t, utils.ValidateDateString()("31-12-2020"), nil)
	utils.AssertEqual(t, utils.ValidateDateString()("nil"), nil)
	utils.AssertEqual(t, utils.ValidateDateString()("12-31-2020"), errors.New(errorMsg))
	utils.AssertEqual(t, utils.ValidateDateString()("2020-12-31"), errors.New(errorMsg))
	utils.AssertEqual(t, utils.ValidateDateString()("2020-31-"), errors.New(errorMsg))
	utils.AssertEqual(t, utils.ValidateDateString()("2020-31"), errors.New(errorMsg))
	utils.AssertEqual(t, utils.ValidateDateString()("2020-"), errors.New(errorMsg))
	utils.AssertEqual(t, utils.ValidateDateString()("2020"), errors.New(errorMsg))
}

func TestTemplateResult(t *testing.T) {
	type TestData struct {
		DataFile string
		Tags     []string
		Notes    []string
	}
	testData := TestData{"random/file/path", []string{"a", "b", "c"}, []string{"foo", "bar"}}
	reportTemplate := `
Stats of "{{.DataFile}}"
  - Number of valid Tags: {{.Tags | numValidTags}}
  - Number of Notes: {{.Notes | len}}
`
	funcMap := template.FuncMap{
		"numValidTags": func(tags []string) int {
			var validTags []string
			for _, elem := range tags {
				if elem > "a" {
					validTags = append(validTags, elem)
				}
			}
			fmt.Println(validTags)
			return len(validTags)
		},
	}
	result, _ := utils.TemplateResult(reportTemplate, funcMap, testData)
	want := `
Stats of "random/file/path"
  - Number of valid Tags: 2
  - Number of Notes: 2
`
	utils.AssertEqual(t, result, want)
}

func TestTerminalSize(t *testing.T) {
	// perhaps stty command doesn't work in tests
	// height, width := utils.TerminalSize()
	// utils.AssertEqual(t, height > 0, true)
	// utils.AssertEqual(t, width > 0, true)
}

func TestPerformShellOperation(t *testing.T) {
	dummyFile := "dummyFile"
	defer func() {
		_, _ = utils.PerformShellOperation("rm -f", dummyFile)
	}()
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
	utils.AssertEqual(t, err, errors.New("exec: no command"))
	_, err = utils.PerformShellOperation("command_do_not_exist", "arg1", "arg2")
	utils.AssertEqual(t, err, errors.New("exec: no command"))
}

func TestHomeDir(t *testing.T) {
	got := utils.HomeDir()
	if len(got) == 0 {
		t.Errorf("HomeDir function returns blank path")
	}
}
