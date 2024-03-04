package errors

type error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func (e *error) Error() string {
	return e.Message
}

func NewError(code int, message string, details interface{}) *error {
	return &error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e *error) SetDetails(details map[string]interface{}) *error {
	e.Details = details
	return e
}

var NotFoundGoodErr = NewError(3, "errors.good.notFound", map[string]struct{}{})
var NotUpdatedGoodErr = NewError(4, "errors.good.notUpdated", map[string]struct{}{})

var TryAgainErr = NewError(99, "errors.tryAgain", map[string]struct{}{})
