package helpers

import (
	"fmt"
	"time"
)

func ParseToDate(day, month, year string) time.Time {
	if year == "" {
		return time.Time{}
	}

	// Default to January 1 if missing
	if month == "" {
		month = "1"
	}
	if day == "" {
		day = "1"
	}

	// Try parsing month as a number first, then as a name
	var t time.Time
	var err error

	// January 2, 2006 is a fixed date used by the time library to recognize the pattern
	t, err = time.Parse("2 1 2006", fmt.Sprintf("%s %s %s", day, month, year))
	if err != nil {
		// Try month as a full name ("January")
		t, err = time.Parse("2 January 2006", fmt.Sprintf("%s %s %s", day, month, year))
		if err != nil {
			// Try month as a abbreviation ("Jan")
			t, err = time.Parse("2 Jan 2006", fmt.Sprintf("%s %s %s", day, month, year))
			// Try month as season ("Spring")
			if err != nil {
				t, err = time.Parse("2 Spring 2006", fmt.Sprintf("%s %s %s", day, month, year))
				if err != nil {
					return time.Time{}
				}
			}
		}
	}

	return t
}
