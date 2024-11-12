package datetime

import (
	"fmt"
	"log/slog"
	"time"
)

// Parse a string representing a datetime
func Parse(datetime string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	parsedTime, err := time.Parse(layout, datetime)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Error occurred when parsing time: %v", err),
		)
	}
	return parsedTime
}
