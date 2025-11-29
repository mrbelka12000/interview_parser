package interview_parser

import (
	"strings"
)

func FormatText(text string) string {

	var (
		lastDot, l int
		newText    strings.Builder
	)

	for i := 0; i < len(text)-1; i++ {
		if text[i] == '?' {
			newText.WriteString(text[l:lastDot+1] + "\n\n") // add answer to prev question
			newText.WriteString(text[lastDot+1:i+1] + "\n") // add question
			l = i + 1
			lastDot = i + 1
		}
		if text[i] == '.' {
			lastDot = i
		}
	}

	newText.WriteString(text[l:lastDot+1] + "\n") // add answer to prev question

	return newText.String()
}
