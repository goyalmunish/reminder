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
	"io"
	"os"
	"testing"
	"time"

	// "github.com/golang/mock/gomock"

	utils "github.com/goyalmunish/reminder/pkg/utils"
)

func TestIsMemberOfSlice(t *testing.T) {
	// with int
	if got := utils.IsMemberOfSlice(100, []int{-100, 0, 100}); got != true {
		t.Error("IsMemberOfSlice failed")
	}
	if got := utils.IsMemberOfSlice(-100, []int{-100, 0, 100}); got != true {
		t.Error("IsMemberOfSlice failed")
	}
	if got := utils.IsMemberOfSlice(0, []int{-100, 0, 100}); got != true {
		t.Error("IsMemberOfSlice failed")
	}
	if got := utils.IsMemberOfSlice(101, []int{-100, 0, 100}); got != false {
		t.Error("IsMemberOfSlice failed")
	}
	if got := utils.IsMemberOfSlice(-99, []int{-100, 0, 100}); got != false {
		t.Error("IsMemberOfSlice failed")
	}
	// with int64
	if got := utils.IsMemberOfSlice(100, []int64{-100, 0, 100}); got != true {
		t.Error("IsMemberOfSlice failed")
	}
	// with float32
	if got := utils.IsMemberOfSlice(100.0, []float32{-100.0, 0.0, 100.0}); got != true {
		t.Error("IsMemberOfSlice failed")
	}
	// with float64
	if got := utils.IsMemberOfSlice(100.0, []int64{-100.0, 0.0, 100.0}); got != true {
		t.Error("IsMemberOfSlice failed")
	}
	// with string
	if got := utils.IsMemberOfSlice("100", []string{"-100", "0", "100"}); got != true {
		t.Error("IsMemberOfSlice failed")
	}
}

func TestGetCommonMembersOfSlices(t *testing.T) {
	// with int
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]int{-100, 0, 100, 1, 10, 8, 2, -51, 4},
			[]int{-21, 100, 0, 8, 4}),
		[]int{0, 100, 8, 4})
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]int{-21, 100, 0, 8, 4},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{100, 0, 8, 4})
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]int{2},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{2})
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]int{},
			[]int{-100, 0, 100, 1, 10, 8, 2, -51, 4}),
		[]int{})
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]int{-100, 0, 100, 1, 10, 8, 2, -51, 4},
			[]int{}),
		[]int{})
	// with int64
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]int64{-100, 0, 100, 1, 10, 8, 2, -51, 4},
			[]int64{-21, 100, 0, 8, 4}),
		[]int64{0, 100, 8, 4})
	// with float32
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]float32{-100.0, 0.0, 100.0, 1.0, 10.0, 8.0, 2.0, -51.0, 4.0},
			[]float32{-21.0, 100.0, 0.0, 8.0, 4.0}),
		[]float32{0.0, 100.0, 8.0, 4.0})
	// with float64
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]float64{-100.0, 0.0, 100.0, 1.0, 10.0, 8.0, 2.0, -51.0, 4.0},
			[]float64{-21.0, 100.0, 0.0, 8.0, 4.0}),
		[]float64{0.0, 100.0, 8.0, 4.0})
	// with string
	utils.AssertEqual(t,
		utils.GetCommonMembersOfSlices([]string{"-100", "0", "100", "1", "10", "8", "2", "-51", "4"},
			[]string{"-21", "100", "0", "8", "4"}),
		[]string{"0", "100", "8", "4"})
}

func TestLogError(t *testing.T) {
	err := errors.New("dummy error")
	utils.LogError(err)
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
	// positive case
	result, err := utils.TemplateResult(reportTemplate, funcMap, testData)
	want := `
Stats of "random/file/path"
  - Number of valid Tags: 2
  - Number of Notes: 2
`
	if err != nil {
		t.Fatalf("TemplateResult returns error %v", err)
	}
	utils.AssertEqual(t, result, want)
	// negative case
	_, err = utils.TemplateResult(reportTemplate, funcMap, "")
	utils.AssertEqual(t, err == nil, false)

}

func TestAssertEqual(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		want        interface{}
		wantedError bool
	}{
		{
			name:        "equal string",
			input:       "a string",
			want:        "a string",
			wantedError: false,
		},
		{
			name:        "equal integer",
			input:       64,
			want:        64,
			wantedError: true,
		},
		{
			name:        "equal slices",
			input:       []int{1, 2, 3},
			want:        []int{1, 2, 3},
			wantedError: true,
		},
	}
	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			// defer func() {
			// 	if subtest.wantedError {
			// 		p := recover()
			// 		if p != nil {
			// 			err := fmt.Errorf("assertion error: %v", p)
			// 			t.Log(err)
			// 		} else {
			// 			t.Errorf("AssertEqual failed at position %d", position)
			// 		}
			// 	}
			// }()
			utils.AssertEqual(t, subtest.input, subtest.want)
		})
	}
}

func TestMatchedTimestamp(t *testing.T) {
	var tests = []struct {
		name              string
		timestampCurrent  int64
		timestampPrevious int64
		timestampNext     int64
		daysBefore        int64
		daysAfter         int64
		currentTime       string
		wantFound         bool
		wantTimestamp     int64
	}{
		{
			name:              "matching current via daysBefore",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-11-15T00:09:00.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1668589200,
		},
		{
			name:              "matching current via dayAfter",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-11-19T00:09:00.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1668589200,
		},
		{
			name:              "matching previous via daysBefore",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-10-15T00:09:00.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1665910800,
		},
		{
			name:              "matching previous via dayAfter",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-10-19T00:09:00.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1665910800,
		},
		{
			name:              "matching next via daysBefore",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-12-15T00:09:00.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1671181200,
		},
		{
			name:              "matching next via dayAfter",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-12-19T00:09:00.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1671181200,
		},
		{
			name:              "matching current via daysBefore (at edge)",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-11-14T09:00:01.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1668589200,
		},
		{
			name:              "matching current via dayAfter (at edge)",
			timestampCurrent:  1668589200, // 16 November 2022 09:00:00 GMT
			timestampPrevious: 1665910800, // 16 October 2022 09:00:00 GMT
			timestampNext:     1671181200, // 16 December 2022 09:00:00 GMT
			currentTime:       "2022-11-21T08:59:59.000Z",
			daysBefore:        2,
			daysAfter:         5,
			wantFound:         true,
			wantTimestamp:     1668589200,
		},
	}
	for position, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			tm, err := time.Parse(time.RFC3339, subtest.currentTime)
			if err != nil {
				t.Fatalf("Test input %q is incorrect", subtest.currentTime)
			}
			utils.CurrentTime = func() time.Time {
				return tm
			}
			gotFound, gotTimestamp := utils.MatchedTimestamp(subtest.timestampCurrent, subtest.timestampNext, subtest.timestampPrevious, subtest.daysBefore, subtest.daysAfter)
			if (gotFound != subtest.wantFound) || (gotTimestamp != subtest.wantTimestamp) {
				t.Errorf("MatchedTimestamp case %q (position=%d) with input (<%+v>, <%+v>, <%+v>, <%+v>, <%+v>) returns (<%+v>, <%+v>); want (<%+v>, <%+v>)", subtest.name, position, subtest.timestampCurrent, subtest.timestampNext, subtest.timestampPrevious, subtest.daysBefore, subtest.daysAfter, gotFound, gotTimestamp, subtest.wantFound, subtest.wantTimestamp)
			}
		})
	}
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

func TestTerminalSize(t *testing.T) {
	t.Skip("perhaps stty command doesn't work in tests")
	height, width, err := utils.TerminalSize()
	if err != nil {
		t.Fatalf("TerminalSize didn't work in test: %v", err)
	}
	utils.AssertEqual(t, height > 0, true)
	utils.AssertEqual(t, width > 0, true)
}

func TestTerminalWidth(t *testing.T) {
	var height, width int
	var err error
	utils.TerminalSize = func() (int, int, error) {
		return height, width, err
	}
	// case 1
	height = 5
	width = 10
	err = nil
	got, errR := utils.TerminalWidth()
	utils.AssertEqual(t, got, 10)
	utils.AssertEqual(t, errR, nil)
	// case 1
	height = 0
	width = 0
	err = errors.New("an error")
	got, errR = utils.TerminalWidth()
	utils.AssertEqual(t, got, 0)
	utils.AssertEqual(t, errR, err)
}

func TestPerformFilePresence(t *testing.T) {
	// case 1
	filePath := "./functions.go"
	err := utils.PerformFilePresence(filePath)
	utils.AssertEqual(t, err == nil, true)
	// case 2
	filePath = "./doesnotexist.file"
	err = utils.PerformFilePresence(filePath)
	utils.AssertEqual(t, err == nil, false)
}

func TestPerformWhich(t *testing.T) {
	// case 1
	err := utils.PerformWhich("go")
	utils.AssertEqual(t, err == nil, true)
	// case 2
	err = utils.PerformWhich("unknown_command")
	utils.AssertEqual(t, err == nil, false)
}
func TestPerformCwdiff(t *testing.T) {
	// create temporary files
	file1, err := os.CreateTemp("./", "temp_file")
	defer os.Remove(file1.Name())
	if err != nil {
		t.Fatal("failed to create a file")
	}
	_, err = io.WriteString(file1, "1\n2\n3")
	if err != nil {
		t.Fatal("failed to write to the file")
	}
	file1.Close()
	file2, err := os.CreateTemp("./", "temp_file")
	defer os.Remove(file2.Name())
	if err != nil {
		t.Fatal("failed to create a file")
	}
	_, err = io.WriteString(file2, "1\n2\n3")
	if err != nil {
		t.Fatal("failed to write to the file")
	}
	file3, err := os.CreateTemp("./", "temp_file")
	defer os.Remove(file3.Name())
	if err != nil {
		t.Fatal("failed to create a file")
	}
	_, err = io.WriteString(file3, "1\n4\n3")
	if err != nil {
		t.Fatal("failed to write to the file")
	}
	file3.Close()
	// case 1: same content
	err = utils.PerformCwdiff(file1.Name(), file2.Name())
	fmt.Println(err)
	utils.AssertEqual(t, err == nil, true)
	// case 2: different content
	err = utils.PerformCwdiff(file1.Name(), file3.Name())
	fmt.Println(err)
	utils.AssertEqual(t, err == nil, false)
}

func TestPerformCat(t *testing.T) {
	// case 1
	filePath := "./functions.go"
	err := utils.PerformCat(filePath)
	utils.AssertEqual(t, err == nil, true)
	// case 2
	filePath = "./doesnotexist.file"
	err = utils.PerformFilePresence(filePath)
	utils.AssertEqual(t, err == nil, false)
}

func TestHomeDir(t *testing.T) {
	got := utils.HomeDir()
	if len(got) == 0 {
		t.Errorf("HomeDir function returns blank path")
	}
}

func TestTryConvertTildaBasedPath(t *testing.T) {
	utils.HomeDir = func() string {
		return "/Users/goyalmunish/"
	}
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "path within home directory",
			input: "~/afile.txt",
			want:  "/Users/goyalmunish/afile.txt",
		},
		{
			name:  "path somewhere in nested diretory",
			input: "~/dir1/dir2/afile.txt",
			want:  "/Users/goyalmunish/dir1/dir2/afile.txt",
		},
		{
			name:  "path starting with '/'",
			input: "/dir1/dir2/afile.txt",
			want:  "/dir1/dir2/afile.txt",
		},
		{
			name:  "path starting with '~'",
			input: "~dir1/afile.txt",
			want:  "~dir1/afile.txt",
		},
		{
			name:  "empty path",
			input: "",
			want:  "",
		},
	}
	for position, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			got := utils.TryConvertTildaBasedPath(subtest.input)
			if got != subtest.want { // or reflect.DeepEqual, https://pkg.go.dev/reflect#DeepEqual
				t.Errorf("TryConvertTildaBasedPath case %q (position=%d) with input <%+v> returns <%+v>; want <%+v>", subtest.name, position, subtest.input, got, subtest.want)
			}
		})
	}
}

func TestAskBoolean(t *testing.T) {
	var tests = []struct {
		name      string
		input     string
		buffer    string
		want      bool
		wantedErr bool
	}{
		{
			name:      "y",
			input:     "y",
			want:      true,
			wantedErr: false,
		},
		{
			name:      "y with padding with spaces",
			input:     "  y ",
			want:      true,
			wantedErr: false,
		},
		{
			name:      "word containing y",
			input:     "ayb",
			want:      false,
			wantedErr: false,
		},
		{
			name:      "n",
			input:     "n",
			want:      false,
			wantedErr: false,
		},
		{
			name:      "word containing n",
			input:     "anb",
			want:      false,
			wantedErr: false,
		},
	}
	for position, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			in, err := os.CreateTemp("./", "a_temp_file")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(in.Name())
			defer in.Close()
			_, err = io.WriteString(in, subtest.input)
			if err != nil {
				t.Fatal(err)
			}
			_, err = in.Seek(0, io.SeekStart)
			if err != nil {
				t.Fatal(err)
			}
			if err != nil {
				t.Fatal(err)
			}
			got, err := utils.PrivateAskBoolean("greeting message", in)
			if (err != nil) != subtest.wantedErr {
				t.Fatalf("AskBoolean case %q (position=%d) with input <%+v> returns error <%v>", subtest.name, position, subtest.input, err)
			}
			if got != subtest.want {
				t.Errorf("AskBoolean case %q (position=%d) with input <%+v> returns <%+v>; want <%+v>", subtest.name, position, subtest.input, got, subtest.want)
			}
		})
	}
}
