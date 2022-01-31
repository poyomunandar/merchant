package common

import (
	"fmt"
	"regexp"
	"strings"
)

func CreateErrorMessageWithCode(errorCode int, params ...string) ErrorMessage {
	var message = ErrorMessageMap[errorCode]
	if len(params) > 0 {
		message = fmt.Sprintf(message, strings.Join(params, ", "))
	}
	return ErrorMessage{
		ErrorCode: errorCode,
		Message:   message,
	}
}

func CreateErrorMessage(errorCode int, message string) ErrorMessage {
	return ErrorMessage{
		ErrorCode: errorCode,
		Message:   message,
	}
}

func CreateSuccessMessage() ErrorMessage {
	return ErrorMessage{
		ErrorCode: ErrorCodeSuccess,
		Message:   "OK",
	}
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
