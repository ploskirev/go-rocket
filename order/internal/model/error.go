package model

import "errors"

var ErrPartNotFound = errors.New("part not found")
var ErrOrderNotFound = errors.New("order not found")
var ErrGenUUID = errors.New("error while generating UUID")
var ErrConflict = errors.New("error conflict")
var ErrBadRequest = errors.New("error bad request")
