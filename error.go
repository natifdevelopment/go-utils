package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/natifdevelopment/go-observability/logging/logger"
	"gorm.io/gorm"
)

func InvalidEnumError[T any](members []T, value string) string {
	var strMembers []string

	for _, member := range members {
		switch v := any(member).(type) {
		case string:
			strMembers = append(strMembers, v)
		case int:
			strMembers = append(strMembers, fmt.Sprintf("%d", v))
		default:
			strMembers = append(strMembers, fmt.Sprintf("%v", member))
		}
	}

	choices := strings.Join(strMembers, ", ")
	return fmt.Sprintf("'%s' is not a valid choice. Valid choices are: %s", value, choices)
}

func HandleGormError(c *gin.Context, err error, operation string, refId, targetId *uuid.UUID) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		SendInvalidInput(c, err)
		return true
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		SendNotFoundError(c)
		return true
	} else if errors.Is(err, gorm.ErrInvalidTransaction) || errors.Is(err, gorm.ErrNotImplemented) ||
		errors.Is(err, gorm.ErrMissingWhereClause) || errors.Is(err, gorm.ErrUnsupportedRelation) ||
		errors.Is(err, gorm.ErrPrimaryKeyRequired) || errors.Is(err, gorm.ErrModelValueRequired) ||
		errors.Is(err, gorm.ErrModelAccessibleFieldsRequired) || errors.Is(err, gorm.ErrSubQueryRequired) ||
		errors.Is(err, gorm.ErrInvalidData) || errors.Is(err, gorm.ErrUnsupportedDriver) ||
		errors.Is(err, gorm.ErrRegistered) || errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrEmptySlice) || errors.Is(err, gorm.ErrDryRunModeUnsupported) ||
		errors.Is(err, gorm.ErrInvalidDB) || errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrInvalidValueOfLength) || errors.Is(err, gorm.ErrPreloadNotAllowed) ||
		errors.Is(err, gorm.ErrForeignKeyViolated) {
		log.Println(errorPrefix, "foreign key constraint violation:", err)
		SendInternalServerErrorWithMessage(c, ConstForeignKeyViolated)
		return true
	} else if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		SendInternalServerError(c, err)
		return true
	}

	return false
}
