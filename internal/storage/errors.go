package storage

import (
	"errors"
)

var (
	ErrBadConn = errors.New("error connecting to DB")
	ErrDuplicateLogin = errors.New("login already taken")
	ErrWrongPassword = errors.New("login/password pair is wrong")
)