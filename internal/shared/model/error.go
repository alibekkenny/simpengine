package model

import (
	"errors"
)

var ErrNoRecord = errors.New("no matching record found")
var ErrEmailOrLoginExists = errors.New("user with such email or login already exists")
