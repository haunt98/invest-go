package invest

import (
	"errors"
	"fmt"
	"time"
)

var dateFormats = []string{
	"2006-01-02",
	"2006/01/02",
}

func dateFromInput(date string, location *time.Location) (string, error) {
	var t time.Time
	for _, dateFormat := range dateFormats {
		var err error
		t, err = time.ParseInLocation(dateFormat, date, location)
		if err != nil {
			continue
		}

		break
	}

	if t.IsZero() {
		return "", errors.New("date is not valid")
	}

	return t.Format(time.RFC3339), nil
}

func dateToOutput(date string, location *time.Location) (string, error) {
	t, err := time.ParseInLocation(time.RFC3339, date, location)
	if err != nil {
		return "", fmt.Errorf("failed to parse time: %w", err)
	}

	return t.Format(dateFormats[0]), nil
}
