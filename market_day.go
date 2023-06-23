package main

import (
	"time"
)

const trading_days uint8 = (1<<time.Monday |
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

var half_days = map[time.Time]bool{
	Day(2023, 7, 3):   true,
	Day(2023, 11, 24): true,
}

func Day(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 15, 0, 0, 0, time.UTC)
}

func IsMarketDay(date time.Time) bool {
	return isTradingDay(date) && !isHoliday(date)
}

func IsFullMarketDay(date time.Time) bool {
	return IsMarketDay(date) && !isHalfDay(date)
}

func isTradingDay(date time.Time) bool {
	return trading_days&(1<<date.Weekday()) != 0
}

func isHoliday(date time.Time) bool {
	return holidays[date]
}

func isHalfDay(date time.Time) bool {
	return half_days[date]
}
