package Strings

import "strings"

func FmtQuery(val string) string {
	return strings.ReplaceAll(val, "'", "''")
}

func FmtLength(val string, length int) string {
	if len(val) > length {
		return val[:length]
	} else {
		return val
	}
}
