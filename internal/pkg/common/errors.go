package common

import (
	"errors"
)

// SignalDone is the signal notification sent to the server handler at the end of process
//
//	 CAUTION !!!!
//		Do not change if you don't know them all!!!
var (
	SignalDone = errors.New("context is over")
)

// Common errors
var (
	ErrInterfaceNonSlice    = errors.New("InterfaceSlice() given a non-slice type")
	ErrNotImplemented       = errors.New("Server.Protocol() is not implemented")
	ErrNonSlice             = errors.New("data response must be a slice")
	ErrInvalidAuthorization = errors.New("invalid authorization")
	ErrInvalidDataType      = errors.New("invalid data type")
	ErrUnAuthorization      = errors.New("unAuthorization")
)

// Common error codes
const (
	ErrorCodeNonSlice            = "NON_SLICE_INPUT"
	ErrorCodeInvalidBson         = "INVALID_BSON_OBJECT"
	ErrorCodeDBTransactionFailed = "DB_TRANSACTION_FAILED"
	ErrorCodeInvalidRequest      = "INVALID"
	ErrorCodeInternalError       = "INTERNAL_ERROR"
)

type Error struct {
	Type    string
	Code    string
	Message string
	Data    any
}

func (e *Error) Error() string {
	return e.Type + " : " + e.Message
}
