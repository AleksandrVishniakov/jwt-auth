package e

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound = errors.New("not found")
	ErrForbiddenAction = errors.New("this action is forbidden")
)