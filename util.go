package main

import (
	"fmt"
	"strings"
)

func CenterLine(s string, width int) string {
	return LeftRightCenter("", "", s, width)
}

func RightJustifyLine(s string, width int) string {
	return LeftRightCenter("", s, "", width)
}

func LeftRightCenter(l string, r string, c string, width int) string {
	nSpaces := width - len(l) - len(r) - len(c)
	lSpaces := (width / 2) - len(l) - (len(c) / 2)
	rSpaces := nSpaces - lSpaces

	return fmt.Sprint(l, strings.Repeat(" ", lSpaces), c, strings.Repeat(" ", rSpaces), r)

}

func CleanBody(s *string) string {
	if s == nil {
		return ""
	}
	str := *s
	str = strings.Replace(str, "\r\n", "\n", -1)
	return str
}
