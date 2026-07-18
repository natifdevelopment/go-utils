package utils

import (
	"regexp"
	"strings"
)

func SanitizeOrderBy(column string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_\.]`)
	return re.ReplaceAllString(column, "")
}

func ValidateSortDirection(sort string) string {
	s := strings.ToLower(sort)
	if s != "asc" && s != "desc" {
		return "desc"
	}
	return s
}

func IsValidOrderByClause(clause string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_\.]+\s+(asc|desc)$`)
	return re.MatchString(strings.TrimSpace(strings.ToLower(clause)))
}
