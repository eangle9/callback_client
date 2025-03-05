package callbackclient

import "dev.azure.com/2f-capital/go-packages/pagewave.git"

type Response struct {
	// OK is only true if the request was successful.
	OK bool `json:"ok"`
	// MetaData contains additional data like filtering, pagination, etc.
	MetaData *MetaData `json:"meta_data,omitempty"`
	// Data contains the actual data of the response.
	Data interface{} `json:"data,omitempty"`
	// Error contains the error detail if the request was not successful.
	Error *Error `json:"error,omitempty"`
}

type MetaData struct {
	pagewave.FilterParams
	// Total is the total number of data without pagination
	Total int `json:"total"`
	// Extra contains other response specific data
	Extra interface{} `json:"extra,omitempty"`
}

type Error struct {
	// Code is the error code. It is not status code
	Code int `json:"code"`
	// Message is the error message.
	Message string `json:"message,omitempty"`
	// Description is the error description.
	Description string `json:"description,omitempty"`
	// StackTrace is the stack trace of the error.
	// It is only returned for debugging
	StackTrace string `json:"stack_trace,omitempty"`
	// FieldError is the error detail for each field, if available that is.
	FieldError []FieldError `json:"field_error,omitempty"`
}

type FieldError struct {
	// Name is the name of the field that caused the error.
	Name string `json:"name"`
	// Description is the error description for this field.
	Description string `json:"description"`
}

type ErrorResponse struct {
	OK            bool  `json:"ok"`
	CallbackError Error `json:"error,omitempty"` //nolint:tagliatelle // depends on authz service
}

func (e *ErrorResponse) Error() string {
	return e.CallbackError.Message
}
