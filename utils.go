package dgc

import "strings"

// stringHasPrefix checks whether or not the string contains one of the given prefixes and returns the string without the prefix
func stringHasPrefix(str string, prefixes []string, ignoreCase bool) (bool, string) {
	for _, prefix := range prefixes {
		stringToCheck := str
		if ignoreCase {
			stringToCheck = strings.ToLower(stringToCheck)
			prefix = strings.ToLower(prefix)
		}
		if strings.HasPrefix(stringToCheck, prefix) {
			return true, strings.TrimSpace(string(str[len(prefix):]))
		}
	}
	return false, ""
}

// equals provides a simple method to check whether or not 2 strings are equal
func equals(str1, str2 string, ignoreCase bool) bool {
	if !ignoreCase {
		return str1 == str2
	}
	return strings.ToLower(str1) == strings.ToLower(str2)
}
