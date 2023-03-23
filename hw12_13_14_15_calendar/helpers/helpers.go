package helpers

import (
	"strings"
	"time"
)

func OnlyDate(InputDateTime time.Time) time.Time {

	year := InputDateTime.Year()
	month := InputDateTime.Month()
	day := InputDateTime.Day()
	location := InputDateTime.Location()

	OnlyDate := time.Date(year, month, day, 0, 0, 0, 0, location)

	return OnlyDate
}

func DateStartTime(InputDateTime time.Time) time.Time {
	return OnlyDate(InputDateTime)
}

func DateEndTime(InputDateTime time.Time) time.Time {

	year := InputDateTime.Year()
	month := InputDateTime.Month()
	day := InputDateTime.Day()
	location := InputDateTime.Location()

	DateEndTime := time.Date(year, month, day, 23, 59, 59, 0, location)

	return DateEndTime
}

func DateBetweenInclude(InputDateTime time.Time, StartTime time.Time, EndTime time.Time) bool {
	aft := InputDateTime.After(StartTime) || InputDateTime.Equal(StartTime)
	bf := InputDateTime.Before(EndTime) || InputDateTime.Equal(EndTime)
	return aft && bf
}

func StringBuild(input ...string) string {

	var sb strings.Builder

	for _, str := range input {
		sb.WriteString(str)
	}

	return sb.String()
}
