package service

import (
	"context"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
)

func (s *Gofermart) PostOrders(ctx context.Context, login string, orderNum int) error {
	numStr := strconv.Itoa(orderNum)
	err := goluhn.Validate(numStr)
	if err != nil {
		return ErrWrongFormat
	}

	if login == "" {
		return ErrMalformedRequest
	}

	err = s.Storage.PostOrders(ctx, login, numStr)
	if err != nil {
		return err
	}
	return nil
}
