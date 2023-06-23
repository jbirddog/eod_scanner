package main

import (
"time"
)

func day(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 14, 0, 0, 0, time.UTC)
}
