package history

import (
	"time"

	"github.com/yagossc/short-url/app"
)

// GetEntriesInInvertval counts the ocurrences of a req given a time
// interval in hours. For example, to get the entries for the past day:
// ocurrences := GetEntriesInInvertval(entries, 24)
// To Get the entries for the past week:
// ocurrences := GetEntriesInInvertval(entries, 24*7)
func GetEntriesInInvertval(entries []app.ReqHistory, interval int) int {
	ocurrences := 0
	currTime := time.Now()

	for _, entry := range entries {
		t := time.Unix(entry.ReqTime, 0)
		diff := int64(currTime.Sub(t).Hours()) / int64(interval)

		if diff < 1 {
			ocurrences++
		}
	}

	return ocurrences
}
