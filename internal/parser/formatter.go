package parser

import (
	"strings"
)

func (p *Parser) FormatText(text string) string {
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

	newText.WriteString(text[l:] + "\n") // add answer to prev question

	return newText.String()
}

func (p *Parser) BatchTranscript(text string) []string {
	lines := strings.Split(text, "\n")
	var (
		out           []string
		l             int
		mod           = 50
		needToCollect bool
	)

	for i, line := range lines {
		if i%mod == 0 && i > 0 {
			if len(line) > 0 && line[len(line)-1] == '?' {
				needToCollect = true
				continue
			}
			out = append(out, strings.Join(lines[l:i+1], "\n"))
			l = i + 1
		}
		if needToCollect {
			out = append(out, strings.Join(lines[l:i+1], "\n"))
			l = i + 1
			needToCollect = false
		}
	}

	out = append(out, strings.Join(lines[l:], "\n"))

	return out
}
