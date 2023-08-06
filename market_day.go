package main

import (
	"time"
)

const tradingDays uint8 = (1<<time.Monday |
	1<<time.Tuesday |
	1<<time.Wednesday |
	1<<time.Thursday |
	1<<time.Friday)

var holidays = map[time.Time]bool{
	Day(2023, 1, 2):   true,
	Day(2023, 1, 16):  true,
	Day(2023, 2, 20):  true,
	Day(2023, 4, 7):   true,
	Day(2023, 5, 29):  true,
	Day(2023, 6, 19):  true,
	Day(2023, 7, 4):   true,
	Day(2023, 9, 4):   true,
	Day(2023, 11, 23): true,
	Day(2023, 12, 25): true,
}

var halfDays = map[time.Time]bool{
	Day(2023, 7, 3):   true,
	Day(2023, 11, 24): true,
}

func Day(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 15, 0, 0, 0, time.UTC)
}

func PreviousMarketDay(date time.Time) time.Time {
	date = date.AddDate(0, 0, -1)

	for !IsMarketDay(date) {
		date = date.AddDate(0, 0, -1)
	}

	return date
}

func PreviousMarketDays(date time.Time, count int) []time.Time {
	days := make([]time.Time, 0, count)

	for len(days) < count {
		date = PreviousMarketDay(date)
		days = append(days, date)
	}

	return days
}

func IsMarketDay(date time.Time) bool {
	return isTradingDay(date) && !isHoliday(date)
}

func IsFullMarketDay(date time.Time) bool {
	return IsMarketDay(date) && !isHalfDay(date)
}

func isTradingDay(date time.Time) bool {
	return tradingDays&(1<<date.Weekday()) != 0
}

func isHoliday(date time.Time) bool {
	return holidays[date]
}

func isHalfDay(date time.Time) bool {
	return halfDays[date]
}
