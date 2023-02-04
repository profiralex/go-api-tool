/*Generated code do not modify it*/
package server

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Response the struct to hold the common response field
type Response struct {
	Data   interface{} `json:"data"`
	Errors []APIError  `json:"errors"`
	Status int         `json:"status"`
}

// APIError holds response error information
type APIError struct {
	Message   string `json:"message"`
	Field     string `json:"field"`
	Reference string `json:"ref"`
}

func (response *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.Status)
	if response.Status == 500 {
		for _, e := range response.Errors {
			log.Warnf("API ERROR: %s %s %s", e.Reference, e.Field, e.Message)
		}
	}
	return nil
}

func CreateSuccessResponse(data interface{}, status ...int) Response {
	finalStatus := http.StatusOK
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   data,
		Errors: nil,
		Status: finalStatus,
	}
}

func CreateAPIErrorsResponse(errors []APIError, status ...int) Response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   nil,
		Errors: errors,
		Status: finalStatus,
	}
}

func CreateAPIErrorResponse(err APIError, status ...int) Response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   nil,
		Errors: []APIError{err},
		Status: finalStatus,
	}
}

func CreateErrorResponse(err error, status ...int) Response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return Response{
		Data:   nil,
		Errors: []APIError{{Message: err.Error()}},
		Status: finalStatus,
	}
}

func RespondSuccess(w http.ResponseWriter, r *http.Request, data interface{}, status ...int) {
	rsp := CreateSuccessResponse(data, status...)
	_ = render.Render(w, r, &rsp)
}

func RespondAPIError(w http.ResponseWriter, r *http.Request, err APIError, status ...int) {
	rsp := CreateAPIErrorResponse(err, status...)
	_ = render.Render(w, r, &rsp)
}

func RespondError(w http.ResponseWriter, r *http.Request, err error, status ...int) {
	rsp := CreateErrorResponse(err, status...)
	_ = render.Render(w, r, &rsp)
}

func RespondValidationError(w http.ResponseWriter, r *http.Request, err error, status ...int) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		RespondError(w, r, err, status...)
		return
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		RespondError(w, r, err, status...)
		return
	}

	var apiErrors []APIError
	for _, fieldError := range validationErrors {

		apiError := APIError{
			Field:   fieldError.Field(),
			Message: fieldError.Error(),
		}
		apiErrors = append(apiErrors, apiError)
	}

	rsp := CreateAPIErrorsResponse(apiErrors, status...)
	_ = render.Render(w, r, &rsp)
}
