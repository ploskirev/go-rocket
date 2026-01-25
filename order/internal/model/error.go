package model

import "errors"

var (
	ErrPartNotFound  = errors.New("part not found")
	ErrOrderNotFound = errors.New("order not found")
	ErrGenUUID       = errors.New("error while generating UUID")
	ErrConflict      = errors.New("error conflict")
	ErrBadRequest    = errors.New("error bad request")
)
