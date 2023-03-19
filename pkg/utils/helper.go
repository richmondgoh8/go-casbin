package utils

import (
	"errors"
	"regexp"
	"strings"
)

// Contains helps to check in a slice of string if it contains a particular string with case insensitivity
func Contains(s []interface{}, str string) bool {
	for _, v := range s {
		if strings.EqualFold(v.(string), str) {
			return true
		}
	}

	return false
}

func PaginationValues(order *string, first *int, after *int) (string, int, int, error) {
	orderBy := "id desc"
	if order != nil {
		orderBy = *order
	}

	valid := regexp.MustCompile("^[A-Za-z0-9_ ]+$")
	if !valid.MatchString(orderBy) {
		return "", 0, 0, errors.New("invalid order by clause")
	}

	// acts as a limit
	firstRows := 10
	if first != nil {
		firstRows = *first
	}

	afterRow := 0
	if after != nil {
		afterRow = *after
	}

	return orderBy, firstRows, afterRow, nil
}

// SanitizeText remove \t & \n from strings
func SanitizeText(stringToReplace string) string {
	stringToReplace = strings.Replace(stringToReplace, "\t", "", -1)
	stringToReplace = strings.Replace(stringToReplace, "\n", "", -1)
	return stringToReplace
}
