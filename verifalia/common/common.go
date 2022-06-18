package common

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Listing

type Direction int

const (
	Forward Direction = iota
	Backward
)

type ListingCursor struct {
	Direction Direction
	Cursor    string
}

type ListingSegment[T any] struct {
	Meta *struct {
		Cursor      string `json:"cursor"`
		IsTruncated bool   `json:"isTruncated"`
	} `json:"meta"`
	Data *[]T `json:"data"`
}

func TimeSpanStringToDuration(timeSpan string) time.Duration {
	timeSpanRe := regexp.MustCompile(`^(?P<Days>\d*\.)?(?P<Hours>\d{1,2})\:(?P<Minutes>\d{1,2}):(?P<Seconds>\d{1,2})$`)
	daysRe := regexp.MustCompile(`^(?P<Days>\d*)\.$`)

	match := timeSpanRe.FindStringSubmatch(timeSpan)

	days := 0
	hours, _ := strconv.Atoi(match[2])
	minutes, _ := strconv.Atoi(match[3])
	seconds, _ := strconv.Atoi(match[4])

	if match[1] != "" {
		daysMatch := daysRe.FindStringSubmatch(match[1])
		days, _ = strconv.Atoi(daysMatch[1])
	}

	return time.Hour*time.Duration(24*days) +
		time.Hour*time.Duration(hours) +
		time.Minute*time.Duration(minutes) +
		time.Second*time.Duration(seconds)
}

// DurationToTimeSpanString generates a time span string from a time.Duration, in the format dd.hh:mm:ss (where dd: days, hh: hours, mm: minutes,
// ss: seconds); the initial dd. part is added only for periods of more than 24 hours.
func DurationToTimeSpanString(duration time.Duration) (result string) {
	const oneMinute = 60
	const oneHour = oneMinute * 60
	const oneDay = oneHour * 60

	result = ""
	totalSeconds := int(duration.Seconds())

	// Days

	if totalSeconds > oneDay {
		var days = totalSeconds / oneDay
		totalSeconds = totalSeconds % oneDay

		result = fmt.Sprintf("%v.", days)
	}

	// Hours

	var hours = totalSeconds / oneHour
	totalSeconds = totalSeconds % oneHour

	result += fmt.Sprintf("%v:", hours)

	// Minutes

	var minutes = totalSeconds / oneMinute
	totalSeconds = totalSeconds % oneMinute

	result += fmt.Sprintf("%v:", minutes)

	// Seconds

	result += fmt.Sprintf("%v", totalSeconds)

	return result
}
