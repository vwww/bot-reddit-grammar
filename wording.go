package main

import (
	"fmt"
	"strings"

	"victorz.ca/grammar"
)

// englishJoin joins a slice of strings using an English conjunction.
func englishJoin(conjunction string, a []string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	case 2:
		return a[0] + " " + conjunction + " " + a[1]
	default:
		commaJoined := strings.Join(a[:len(a)-1], ", ")
		return commaJoined + ", " + conjunction + " " + a[len(a)-1]
	}
}

func makeWording(title, text, author string) string {
	correctionsTitle, _ := grammar.Load(title)
	correctionsText, _ := grammar.Load(text)

	corrections := append(correctionsTitle, correctionsText...)

	if len(corrections) == 0 {
		return ""
	}

	errPlural := "an error"
	if len(corrections) != 1 {
		errPlural = "errors"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("/u/%v, the Reddit Grammar Bot has detected %v in your post:", author, errPlural))
	for _, c := range corrections {
		sb.WriteString("\n\n >")
		sb.WriteString(c)
	}

	return sb.String()
}
