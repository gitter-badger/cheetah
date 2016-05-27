package stringutil

import (
	"testing"
	"strings"
)

func TestSubString(t *testing.T) {
	s := "ABCDEFG"
	subString := SubString(s, 1, 5)

	trueSubStr := "BCDEF"
	if !strings.EqualFold(trueSubStr, subString) {
		t.Errorf("SubString(\"%s\", 1, 5) != \"%s\".\nthe wrong result: \"%s\"", s, trueSubStr, subString)
	}
}

func TestUpperFirst(t *testing.T) {
	s := "abcdefg"
	upperStr := UpperFirst(s)

	trueStr := "Abcdefg"
	if !strings.EqualFold(upperStr, trueStr) {
		t.Errorf("UpperFirst(\"%s\") != \"%s\".\nthe wrong result: \"%s\"", s, trueStr, upperStr)
	}
}

func TestLowerFirst(t *testing.T) {
	s := "ABCDEFG"
	lowerStr := LowerFirst(s)

	trueStr := "aBCDEFG"
	if !strings.EqualFold(lowerStr, trueStr) {
		t.Errorf("LowerFirst(\"%s\") != \"%s\".\nthe wrong result: \"%s\"", s, trueStr, lowerStr)
	}
}