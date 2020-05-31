package shortener

import "strings"

var valueMapping = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i",
	"j", "k", "l", "m", "n", "o", "p", "q", "r",
	"s", "t", "u", "v", "w", "x", "y", "z", "A",
	"B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S",
	"T", "U", "V", "W", "X", "Y", "Z", "0", "1",
	"2", "3", "4", "5", "6", "7", "8", "9"}

// GetShortURL converts a given decimal
// number to our base62 string
func GetShortURL(num int64) string {
	var converted strings.Builder

	for num > 0 {
		converted.WriteString(valueMapping[num%62])
		num = num / 62
	}

	return converted.String()
}
