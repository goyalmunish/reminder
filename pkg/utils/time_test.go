package utils_test

import (
	"errors"
	"testing"
	"time"

	utils "github.com/goyalmunish/reminder/pkg/utils"
)

func TestCurrentUnixTimestamp(t *testing.T) {
	utils.CurrentTime = func() time.Time { return time.Now() }
	want := time.Now().Unix()
	output := utils.CurrentUnixTimestamp()
	utils.AssertEqual(t, output, want)
}

func TestUnixTimestampToTime(t *testing.T) {
	currentTime := time.Now()
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

func TestYearForDueDateDDMM(t *testing.T) {
	var tests = []struct {
		name        string
		input       string
		currentTime string // RFC3339 string
		want        int
		wantedErr   bool
	}{
		{
			name:        "DDMM before currentTime",
			input:       "01-03",
			currentTime: "2022-04-02T00:18:18.929Z",
			want:        2023,
			wantedErr:   false,
		},
		{
			name:        "DDMM before currentTime",
			input:       "25-11",
			currentTime: "2022-11-27T11:18:18.929+11:00",
			want:        2023,
			wantedErr:   false,
		},
		{
			name:        "DDMM after currentTime",
			input:       "05-07",
			currentTime: "2022-04-02T00:18:18.929Z",
			want:        2022,
			wantedErr:   false,
		},
		{
			name:        "DDMM after currentTime",
			input:       "28-12",
			currentTime: "2022-11-27T11:18:18.929+11:00",
			want:        2022,
			wantedErr:   false,
		},
		{
			name:        "simple DDMM",
			input:       "1-3",
			currentTime: "2022-04-02T00:18:18.929Z",
			want:        2023,
			wantedErr:   false,
		},
		{
			name:        "another simple DDMM",
			input:       "3-4",
			currentTime: "2022-04-02T00:18:18.929Z",
			want:        2022,
			wantedErr:   false,
		},
		{
			name:        "invalid input",
			input:       "4",
			currentTime: "2022-04-02T00:18:18.929Z",
			wantedErr:   true,
		},
		{
			name:        "invalid input",
			input:       "99-99",
			currentTime: "2022-04-02T00:18:18.929Z",
			wantedErr:   true,
		},
	}
	for position, subtest := range tests {
		tm, err := time.Parse(time.RFC3339, subtest.currentTime)
		if err != nil {
			t.Fatalf("Test input %q is incorrect", subtest.input)
		}
		utils.CurrentTime = func() time.Time {
			return tm
		}
		got, err := utils.YearForDueDateDDMM(subtest.input)
		if (err != nil) != subtest.wantedErr {
			t.Fatalf("YearForDueDateDDMM case %q (position=%d) with input <%+v> returns error <%v>", subtest.name, position, subtest.input, err)
		}
		if got != subtest.want {
			t.Errorf("YearForDueDateDDMM case %q (position=%d) failed for input %q; returns <%+v>; wants <%+v>", subtest.name, position, tm, got, subtest.want)
		}
	}
}

func TestStrToTime(t *testing.T) {
	// note: refer this format to quickly write table based tests
	var tests = []struct {
		name      string
		input     string // RFC3339 string
		timezone  string
		want      string
		wantedErr bool // whether an error was expected
	}{
		{name: "time in GMT", input: "2022-12-28T00:18:18.929Z", want: "2022-12-28T00:18:18Z"},
		{name: "time in Melbourne/Australia", input: "2022-12-28T11:18:18.929+11:00", want: "2022-12-28T11:18:18+11:00"},
		{name: "time in Melbourne/Australia to UTC", input: "2022-12-27T11:18:18.929+11:00", want: "2022-12-27T00:18:18Z", timezone: "UTC"},
		{name: "with invalid timeze", input: "2022-12-27T11:18:18.929+11:00", want: "2022-12-27T11:18:18+11:00", timezone: "INVALID", wantedErr: true},
		{name: "with invalid time string", input: "SomeInvalidValue", want: "0001-01-01T00:00:00Z", timezone: "UTC", wantedErr: true},
	}
	for position, subtest := range tests {
		got, err := utils.StrToTime(subtest.input, subtest.timezone)
		if (err != nil) != subtest.wantedErr {
			t.Fatalf("StrToTime case %q (position=%d) failed for input %q with error %q", subtest.name, position, subtest.input, err)
		}
		gotStr := got.Format(time.RFC3339)
		if gotStr != subtest.want {
			t.Errorf("StrToTime case %q (position=%d) failed for input %q; returns <%+v>; wants <%+v>", subtest.name, position, subtest.input, gotStr, subtest.want)
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
func TestGetLocalZone(t *testing.T) {
	var tests = []struct {
		name        string
		currentTime string // RFC3339 string
		wantAbbr    string
		wantDur     time.Duration
	}{
		{
			name:        "time in UTC",
			currentTime: "2022-04-02T00:18:18.929Z",
			wantAbbr:    "UTC",
			wantDur:     0,
		},
		{
			name:        "time in IST",
			currentTime: "2022-11-27T11:18:18.929+05:30",
			wantAbbr:    "",
			wantDur:     time.Duration(5*time.Hour + 30*time.Minute),
		},
		{
			name:        "time in AEDT",
			currentTime: "2022-11-27T11:18:18.929+11:00",
			wantAbbr:    "AEDT",
			wantDur:     time.Duration(11 * time.Hour),
		},
	}
	for position, subtest := range tests {
		tm, err := time.Parse(time.RFC3339, subtest.currentTime)
		if err != nil {
			t.Fatalf("Test input %q is incorrect", subtest.currentTime)
		}
		utils.CurrentTime = func() time.Time {
			return tm
		}
		abbr, dur := utils.GetLocalZone()
		if (abbr != subtest.wantAbbr) || (dur != subtest.wantDur) {
			t.Errorf("GetLocalZone case %q (position=%d) failed for input %q; returns <%+v>, <%+v>; wants <%+v>, <%+v>", subtest.name, position, tm, abbr, dur, subtest.wantAbbr, subtest.wantDur)
		}
	}
}

func TestGetZoneFromLocation(t *testing.T) {
	var tests = []struct {
		name      string
		input     string
		want      time.Duration
		wantedErr bool
	}{
		{
			name:      "GMT",
			input:     "GMT",
			want:      time.Duration(0 * time.Hour),
			wantedErr: false,
		},
		{
			name:      "UTC",
			input:     "UTC",
			want:      time.Duration(0 * time.Hour),
			wantedErr: false,
		},
		{
			name:      "Asia/Singapore",
			input:     "Asia/Singapore",
			want:      time.Duration(8 * time.Hour),
			wantedErr: false,
		},
		{
			name:      "Asia/Calcutta",
			input:     "Asia/Calcutta",
			want:      time.Duration(5*time.Hour + 30*time.Minute),
			wantedErr: false,
		},
		{
			name:      "Australia/Melbourne",
			input:     "Australia/Melbourne",
			want:      time.Duration(11 * time.Hour),
			wantedErr: false,
		},
		{
			name:      "Invalid Location",
			input:     "SomeInvalidValue",
			want:      0,
			wantedErr: true,
		},
	}
	for position, subtest := range tests {
		got, err := utils.GetZoneFromLocation(subtest.input)
		if (err != nil) != subtest.wantedErr {
			t.Fatalf("GetZoneFromLocation case %q (position=%d) with input <%+v> returns error <%v>", subtest.name, position, subtest.input, err)
		}
		if got != subtest.want {
			t.Errorf("GetZoneFromLocation case %q (position=%d) failed for input %q; returns <%+v>; wants <%+v>", subtest.name, position, subtest.input, got, subtest.want)
		}
	}
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
	utils.AssertEqual(t, utils.ValidateDateString()(2020), errors.New("Invalid type int"))
}
