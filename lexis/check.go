package lexis

import (
	"strconv"
	"strings"
	"unicode"
)

func checkValueInArray(value string, array []string) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}

	return false
}

func checkValueIsString(value string) bool {
	if strings.Contains(value, "\"") {
		return true
	}

	return false
}

func checkValueIsNumber(value string) bool {
	if checkValueIsString(value) {
		return false
	}

	if _, err := strconv.Atoi(string(value)); err != nil {
		return false
	}

	return true
}

func checkValueIsLetter(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func checkValueIsVariable(value string) bool {
	if checkValueIsString(value) {
		return false
	}

	if _, err := strconv.Atoi(string(value)); err == nil {
		return false
	}

	if !checkValueIsLetter(value) {
		return false
	}

	return true
}
