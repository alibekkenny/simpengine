package model

import (
	"errors"
	"log"
	"net/http"
)

var ErrNoRecord = errors.New("no matching record found")
var ErrEmailOrLoginExists = errors.New("user with such email or login already exists")
var ErrInvalidBody = errors.New("invalid body")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInternal = errors.New("internal server error")
var ErrUniqueViolation = errors.New("unique constraint violation")

func ErrorStatus(err error) (int, string) {
	log.Println(err)
	switch {
	case errors.Is(err, ErrInvalidBody):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, ErrNoRecord):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, ErrUniqueViolation):
		return http.StatusConflict, err.Error()
	case err != nil:
		return http.StatusInternalServerError, ErrInternal.Error()
	default:
		return http.StatusOK, ""
	}
}
