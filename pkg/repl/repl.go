package repl

import (
	"strings"
)

func CleanInput(text string) []string {
	lowered := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowered)
	split := strings.Fields(trimmed)

	return split
}
