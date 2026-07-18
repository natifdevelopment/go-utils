package utils

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var filterColumnRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`)

func ParseFilters(c *gin.Context) map[string]map[string]string {
	filters := make(map[string]map[string]string)

	for key, values := range c.Request.URL.Query() {
		if !strings.HasPrefix(key, "filters[") {
			continue
		}

		key = strings.TrimPrefix(key, "filters[")
		key = strings.TrimSuffix(key, "]")

		parts := strings.Split(key, "][")
		if len(parts) != 2 {
			continue
		}

		column := parts[0]
		operator := parts[1]
		value := values[0]

		if !filterColumnRegex.MatchString(column) {
			continue
		}

		if _, exists := filters[column]; !exists {
			filters[column] = make(map[string]string)
		}
		filters[column][operator] = value
	}

	return filters
}
