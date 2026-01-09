package orderservice

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultTimeout = time.Second
)

type OrdersInfo struct {
	UserID uint64
	Amount uint64
}

func (o OrdersInfo) String() string {
	return fmt.Sprintf("Orders: %v", o.Amount)
}

type Client struct {
	timeout       time.Duration
	delay         time.Duration
	isReturnError bool
}

func New(opts ...OptFunc) *Client {
	c := &Client{timeout: defaultTimeout}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetOrdersByUserID - mock function
func (c *Client) GetOrdersByUserID(ctx context.Context, userID uint64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if c.isReturnError {
		return "", errors.New("internal order error")
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(c.delay):
		return OrdersInfo{UserID: userID, Amount: 5}.String(), nil
	}
}
