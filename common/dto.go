package common

import "net/http"

const (
	CodeSuccess        = "SUCCESS"
	CodeError          = "ERROR"
	ParamFormatError   = "PARAM_FORMAT_ERROR"
	ParamNotExistError = "PARAM_NOT_EXIST_ERROR"
	AppNotExistError   = "APP_NOT_EXIST_ERROR"
	AppExistError      = "APP_EXIST_ERROR"

	InternalServiceErr = "INTERNAL_SERVICE_ERROR"

	Unauthorized = "UNAUTHORIZED"
)

type Result struct {
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewOkResult(msg string) *Result {
	return &Result{
		Code:    CodeSuccess,
		Data:    nil,
		Message: msg,
	}
}

func NewDataResult(data interface{}) *Result {
	return &Result{
		Code:    CodeSuccess,
		Data:    data,
		Message: "",
	}
}

func NewErrorResult(code, msg string) *Result {
	return &Result{
		Code:    code,
		Data:    nil,
		Message: msg,
	}
}

type Error struct {
	Status int
	Code   string
	Msg    string
}

func BadRequestErr(code string, msg string) error {
	return &Error{
		Status: http.StatusBadRequest,
		Code:   code,
		Msg:    msg,
	}
}

func NotFoundErr(code string, msg string) error {
	return &Error{
		Status: http.StatusNotFound,
		Code:   code,
		Msg:    msg,
	}
}

func UnauthorizedErr() error {
	return &Error{
		Status: http.StatusUnauthorized,
		Code:   "Unauthorized",
		Msg:    "",
	}
}

func InternalServerErr(code string, msg string) error {
	return &Error{
		Status: http.StatusInternalServerError,
		Code:   code,
		Msg:    msg,
	}
}

func (e *Error) Error() string {
	return e.Msg

}
