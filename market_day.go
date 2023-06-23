package main

import (
"time"
)

func IsMarketDay(date time.Time) bool {
if date.Year() != 2023 {
panic("Unsupported year provided.")
}
	return false
}
