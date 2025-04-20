package entity

import "errors"

var (
	ErrInvalidCEP  = errors.New("invalid zipcode")
	ErrCEPNotFound = errors.New("can not find zipcode")
)
