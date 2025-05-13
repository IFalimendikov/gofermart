package service

import (
	"context"
	"strconv"
	"github.com/ShiraazMoollatjie/goluhn"
)

func (s *Gofermart) PostOrders(ctx context.Context, userID string, orderNum int) error {
	numStr := strconv.Itoa(orderNum)
	err := goluhn.Validate(numStr)
	if err != nil {
		return ErrWrongFormat
	}

	err = s.Storage.PostOrders(ctx, userID, numStr)
	if err != nil {
		return err
	}
	return nil
}