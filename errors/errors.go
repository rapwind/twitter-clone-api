package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

// ResponseError is error struct
type ResponseError struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Params  interface{} `json:"-"`
}

// Error is error interface method
func (e *ResponseError) Error() string {
	return e.Message
}

// ErrDataNotFound is wrap mgo.ErrNotFound, redis.ErrNil.
var ErrDataNotFound = fmt.Errorf("data not found")

// IsDataNotFound judges errors.ErrDataNotFound, mgo.ErrNotFound or not.
func IsDataNotFound(err error) (yes bool) {
	return err == ErrDataNotFound || err == mgo.ErrNotFound
}

// BadParams is caused when params are not enough or bad
func BadParams(name string, value string) *ResponseError {
	params := map[string]string{}
	params[name] = value
	return &ResponseError{
		Status:  http.StatusBadRequest,
		Code:    "bad_params",
		Message: "bad params",
		Params:  params,
	}
}

// Unauthorized is caused when authorization required
func Unauthorized() *ResponseError {
	return &ResponseError{
		Status:  http.StatusUnauthorized,
		Code:    "unauthorized",
		Message: "headers are not enough",
	}
}

// Send is reply error
// if err is ...
// 	*ResponseError, responds err.Status
// 	mgo.ErrNotFound and errors.DataNotFoundError, responds NotFound
// 	mgo.IsDup(err) == true, responds Conflict
//	else, responds InternalServerError
func Send(c *gin.Context, err error) {
	c.Abort()
	if e, ok := err.(*ResponseError); ok {
		c.JSON(e.Status, e)
	} else if err == ErrDataNotFound || err == mgo.ErrNotFound {
		e := NotFound()
		c.JSON(e.Status, e)
	} else if mgo.IsDup(err) {
		e := DataConflict()
		c.JSON(e.Status, e)
	} else {
		c.JSON(http.StatusInternalServerError, &ResponseError{
			Code:    "internal_server_error",
			Message: err.Error(),
		})
	}
}

// NotFound is caused when data is not found
func NotFound() *ResponseError {
	return &ResponseError{
		Status:  http.StatusNotFound,
		Code:    "not_found",
		Message: "not found",
	}
}

// DataConflict is caused when already registered data
func DataConflict() *ResponseError {
	return &ResponseError{
		Status:  http.StatusConflict,
		Code:    "data_conflict",
		Message: "already registered data",
	}
}

// EmailUnverified is caused when bad precondition
func EmailUnverified() *ResponseError {
	return &ResponseError{
		Status:  http.StatusPreconditionFailed,
		Code:    "email_unverified",
		Message: "email unverified",
	}
}

// ImageSizeTooLarge is caused when data size is too large
func ImageSizeTooLarge() *ResponseError {
	return &ResponseError{
		Status:  http.StatusRequestEntityTooLarge,
		Code:    "image_size_too_large",
		Message: "image size is too large",
	}
}

// UnsupportedMediaType is caused when send unsupported media type
func UnsupportedMediaType() *ResponseError {
	return &ResponseError{
		Status:  http.StatusUnsupportedMediaType,
		Code:    "unsupported_media_type",
		Message: "unsupported media type",
	}
}

// RequestEntityTooLarge is caused when data size is too large
func RequestEntityTooLarge() *ResponseError {
	return &ResponseError{
		Status:  http.StatusRequestEntityTooLarge,
		Code:    "request_entity_too_large",
		Message: "request entity too large",
	}
}

// InvalidSyntax is caused when syntax is invalid
func InvalidSyntax() *ResponseError {
	return &ResponseError{
		Status:  http.StatusBadRequest,
		Code:    "invalid_syntax",
		Message: "Invalid syntax in request",
	}
}
