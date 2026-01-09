package profileservice

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultTimeout = time.Second
)

type ProfileInfo struct {
	ID     uint64
	UserID uint64
	Name   string
}

func (p ProfileInfo) String() string {
	return fmt.Sprintf("Name: %v", p.Name)
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

// GetProfileByUserID - mock function
func (c *Client) GetProfileByUserID(ctx context.Context, userID uint64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if c.isReturnError {
		return "", errors.New("internal profile error")
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(c.delay):
		return ProfileInfo{UserID: userID, Name: "Alice"}.String(), nil
	}
}
