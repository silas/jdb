package jdb

import (
	"strings"
)

func placeholders(count int) string {
	if count < 1 {
		return ""
	}

	return strings.Repeat(", ?", count)[2:]
}
