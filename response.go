package utils

import (
	"errors"
	"fmt"

	log "github.com/natifdevelopment/go-observability/logging/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	httpresponse "github.com/natifdevelopment/go-openapi/httpresponse"
	"github.com/natifdevelopment/go-openapi/schemas"
)

const errorPrefix = "Error: "

// LogFunc is a pluggable system log function. Services should set this
// to their own logging implementation. If nil, logging is skipped.
var LogFunc func(c *gin.Context, module, message string, refId, targetId any)

type ErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	FieldName   string `json:"field"`
	RuleMessage string `json:"rule"`
}

func CreateValidationErrorResponse(err error) ErrorResponse {
	var response ErrorResponse

	for _, err := range err.(validator.ValidationErrors) {
		validationError := ValidationError{
			FieldName:   err.Field(),
			RuleMessage: fmt.Sprintf("%s is %s", err.Field(), err.Tag()),
		}
		response.Errors = append(response.Errors, validationError)
	}
	return response
}

func toSchemaValidationErrors(err error) []schemas.ValidationError {
	var result []schemas.ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		result = append(result, schemas.ValidationError{
			Field:       e.Field(),
			Rule:        e.Tag(),
			RuleMessage: fmt.Sprintf("%s is %s", e.Field(), e.Tag()),
		})
	}
	return result
}

func SendInvalidInput(c *gin.Context, err error) {
	if err != nil {
		log.Println(errorPrefix, err)
		if LogFunc != nil {
			LogFunc(c, "base", "Invalid input parameter: "+err.Error(), nil, nil)
		}
	}
	if verr, ok := err.(validator.ValidationErrors); ok {
		httpresponse.SendValidationError(c, "Invalid input parameters", toSchemaValidationErrors(verr))
		return
	}
	httpresponse.SendBadRequest(c, "Invalid input parameters")
}

func SendSuccess(c *gin.Context) {
	httpresponse.SendSuccess(c)
}

func SendSuccessWithData(c *gin.Context, data any) {
	httpresponse.SendSuccessWithData(c, data)
}

func SendInternalServerErrorWithData(c *gin.Context, err error, data any) {
	if errors.As(err, &BboResponseError{}) {
		httpresponse.SendSuccessWithMessage(c, err.Error(), data)
		return
	}

	if err != nil {
		log.Println(errorPrefix, c.FullPath(), err)
		if LogFunc != nil {
			LogFunc(c, "base", "Internal server error: "+err.Error(), nil, nil)
		}
	}
	httpresponse.SendInternalServerError(c, ConstInternalServerError)
}

func SendInternalServerError(c *gin.Context, err error) {
	SendInternalServerErrorWithData(c, err, nil)
}

func SendNotFoundError(c *gin.Context) {
	httpresponse.SendNotFound(c, ConstNotFound)
}

func SendForbiddenError(c *gin.Context, err error) {
	if err != nil {
		log.Println(errorPrefix, err)
		if LogFunc != nil {
			LogFunc(c, "base", "Access forbidden: "+err.Error(), nil, nil)
		}
	}
	httpresponse.SendForbidden(c, "Anda tidak memiliki akses ke halaman ini")
}

func SendInternalServerErrorWithMessage(c *gin.Context, message string) {
	httpresponse.SendInternalServerError(c, message)
}

func GenInternalServerError() error {
	return errors.New("error: " + ConstInternalServerError)
}

func GenInternalServerErrorWithErr(err error) error {
	if err != nil {
		if LogFunc != nil {
			LogFunc(nil, "base", "Internal server error: "+err.Error(), nil, nil)
		}
	}
	return errors.New("error: " + ConstInternalServerError)
}

func LogError(err error) {
	if err != nil {
		if LogFunc != nil {
			LogFunc(nil, "base", errorPrefix+err.Error(), nil, nil)
		}
	}
}
