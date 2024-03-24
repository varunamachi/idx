package core

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidState = errors.New("invalid state")
	ErrInvalidRole  = errors.New("invalid role")
)
