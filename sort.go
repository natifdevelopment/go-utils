package utils

import (
	"regexp"
	"strings"
)

var sanitizeOrderByRegex = regexp.MustCompile(`[^a-zA-Z0-9_\.]`)

func SanitizeOrderBy(column string) string {
	return sanitizeOrderByRegex.ReplaceAllString(column, "")
}

func ValidateSortDirection(sort string) string {
	s := strings.ToLower(sort)
	if s != "asc" && s != "desc" {
		return "desc"
	}
	return s
}

var isValidOrderByClauseRegex = regexp.MustCompile(`^[a-zA-Z0-9_\.]+\s+(asc|desc)$`)

func IsValidOrderByClause(clause string) bool {
	return isValidOrderByClauseRegex.MatchString(strings.TrimSpace(strings.ToLower(clause)))
}
