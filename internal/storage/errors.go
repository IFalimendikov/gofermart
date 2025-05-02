package storage

import (
	"errors"
)

var (
	ErrBadConn = errors.New("error connecting to DB")
)