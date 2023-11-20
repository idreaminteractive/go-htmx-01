package helpers

import "strings"

func IsEmptyOrWhiteSpace(s string) bool {
	return len(strings.TrimSpace(s)) == 0

}
