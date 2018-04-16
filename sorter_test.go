package main

import "testing"

func Test_SortMixed(t *testing.T) {
	tags := []string{"latest", "1.0.1"}

	compareStringNumber := func(str1, str2 string) bool {
		return extractNumberFromString(str1) < extractNumberFromString(str2)
	}
	Compare(compareStringNumber).Sort(tags)

	if tags[0] != "1.0.1" && tags[1] != "latest" {
		t.Errorf("ordering incorrect when checking mixed tags")
	}
}

func Test_SortAllDigits(t *testing.T) {
	tags := []string{"1.2.1", "1.0.1"}

	compareStringNumber := func(str1, str2 string) bool {
		return extractNumberFromString(str1) < extractNumberFromString(str2)
	}
	Compare(compareStringNumber).Sort(tags)

	if tags[0] != "1.0.1" && tags[1] != "1.2.1" {
		t.Errorf("ordering incorrect in all digits tags")
	}
}
