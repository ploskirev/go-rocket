package model

import "errors"

var (
	ErrBadRequest = errors.New("error bad request")
	ErrConflict   = errors.New("error conflict")

	ErrGenUUID = errors.New("error while generating UUID")

	ErrPartNotFound = errors.New("part not found")
	ErrGetListParts = errors.New("error get list parts")

	ErrOrderNotFound = errors.New("order not found")
	ErrPayOrder      = errors.New("error pay order")

	ErrCreateOrder = errors.New("error while creating order")
)
