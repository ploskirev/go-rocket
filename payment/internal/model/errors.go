package model

import "errors"

var (
	ErrPartNotFound = errors.New("part not found")
	ErrGenerateUUID = errors.New("error while generating uuid")
)
