package array

import (
	"log"
	"regexp"
	"strings"
)

func ConsistStringInArrayByRegex(haystack []string, needle string) bool {
	for _, v := range haystack {
		pattern := "^"+v+"$"
		match, err := regexp.Match(pattern, []byte(needle))
		if err != nil{
			log.Println(err)
		}
		if match {
			return true
		}
	}
	return false
}

func InStringArrayExists(needle string, haystack []string) bool{
	for _, v := range haystack{
		if v == needle{
			return true
		}
	}

	return false
}


func InStringArrayExistsIgnoreCase(needle string, haystack []string) bool{
	for _, v := range haystack{
		if strings.ToLower(v) == strings.ToLower(needle) {
			return true
		}
	}

	return false
}
