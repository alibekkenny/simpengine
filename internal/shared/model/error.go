package model

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var ErrNoRecord = errors.New("no matching record found")
var ErrEmailOrLoginExists = errors.New("user with such email or login already exists")
var ErrInvalidBody = errors.New("invalid body")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrForbidden = errors.New("forbidden")
var ErrInternal = errors.New("internal server error")
var ErrUniqueViolation = errors.New("unique constraint violation")
var ErrInvalidState = errors.New("invalid state")
var ErrValidation = errors.New("validation error")
var ErrInvalidParams = errors.New("invalid params")

func ErrorStatus(err error) (int, string, string) {
	log.Println(err)
	switch {
	case errors.Is(err, ErrInvalidBody):
		return http.StatusBadRequest, "INVALID_BODY", err.Error()
	case errors.Is(err, ErrInvalidParams):
		return http.StatusBadRequest, "INVALID_PARAMETERS", err.Error()
	case errors.Is(err, ErrInvalidCredentials):
		return http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error()
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "FORBIDDEN", err.Error()
	case errors.Is(err, ErrNoRecord):
		return http.StatusNotFound, "NOT_FOUND", err.Error()
	case errors.Is(err, ErrUniqueViolation):
		return http.StatusConflict, "UNIQUE_VIOLATION", err.Error()
	case errors.Is(err, ErrInvalidState):
		return http.StatusConflict, "INVALID_STATE", err.Error()
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest, "VALIDATION_ERROR", err.Error()
	case err != nil:
		return http.StatusInternalServerError, "INTERNAL_ERROR", ErrInternal.Error()
	default:
		return http.StatusOK, "", ""
	}
}

// ErrorResponse represents error details returned by API
type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_BODY"`
	Message string `json:"message" example:"Invalid request"`
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	status, code, msg := ErrorStatus(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: msg,
	})
}

func WriteRawErrorResponse(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: msg,
	})
}
