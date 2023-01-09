package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

// Location variable provides location info for `time`.
// It can be set to update behavior of UnixTimestampToTime.
var Location *time.Location

// CurrentTime function gets current time.
// Note: It is deliberately defined as a variable to make it
// easer to patch in tests.
var CurrentTime = func() time.Time {
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
// In either case, the value of the time (in seconds) remains same, the
// use of the Location just changes how time is displayed.
// For example, 5:30 PM in SGT is equivalent to 12 Noon in India.
func UnixTimestampToTime(unixTimestamp int64) time.Time {
	t := time.Unix(unixTimestamp, 0)
	if Location == nil {
		return t
	}
	return t.In(Location)
}

// UnixTimestampToTimeStr function converts unix timestamp to time string.
func UnixTimestampToTimeStr(unixTimestamp int64, timeFormat string) string {
	if unixTimestamp > 0 {
		return UnixTimestampToTime(unixTimestamp).Format(timeFormat)
	}
	return "nil"
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
	// try with current year
	year := currentTime.Year()
	dateString := fmt.Sprintf("%s-%d", dateMonth, year)
	tryTimeValue, err := time.Parse(format, dateString)
	if err != nil {
		return 0, err
	}
	if tryTimeValue.Unix() <= currentTime.Unix() {
		// the due date falls before current date in current year
		// so, select next year instead
		year += 1
	}
	return year, nil
}

// StrToTime converts RFC3339 time sting to time.Time, and sets location to
// given timezone. If location is blank, then it returns the time as it is.
func StrToTime(tString string, timezone string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, tString)
	if err != nil {
		return t, fmt.Errorf("Unable to parse the time %v: %w", tString, err)
	}
	if timezone == "" {
		return t, nil
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return t, fmt.Errorf("Unable to parse the timezone %v: %w", timezone, err)
	}
	fmt.Println(t.In(location).Format(time.RFC3339))
	return t.In(location), nil
}

// TimeToStr converts time.Time to RFC3339 time string.
func TimeToStr(t time.Time) string {
	return t.Format(time.RFC3339)
}

// GetLocalZone returns the local timezone abbreviation and offset (in time.Duration).
func GetLocalZone() (string, time.Duration) {
	abbr, seconds := CurrentTime().Zone()
	dur := time.Duration(seconds * int(time.Second))
	return abbr, dur
}

// GetZoneFromLocation returns zone offset (in time.Duration) for given location string like "Melbourne/Australia".
func GetZoneFromLocation(loc string) (time.Duration, error) {
	location, err := time.LoadLocation(loc)
	if err != nil {
		return time.Duration(0 * time.Second), err
	}
	_, seconds := CurrentTime().In(location).Zone()
	dur := time.Duration(seconds * int(time.Second))

	return dur, nil
}

// ValidateDateString function validates date string (DD-MM-YYYY) or (DD-MM).
// nil is also valid input
func ValidateDateString() survey.Validator {
	// return a validator that checks the length of the string
	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			// if the string is shorter than the given value
			input := strings.TrimSpace(str)
			re := regexp.MustCompile(`^((0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])(-((19|20)\d\d))?|(nil))$`)
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
