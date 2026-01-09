package orderservice

import "time"

type OptFunc func(u *Client)

func WithTimeout(timeout time.Duration) OptFunc {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithDelay(delay time.Duration) OptFunc {
	return func(c *Client) {
		c.delay = delay
	}
}

func WithReturnError() OptFunc {
	return func(c *Client) {
		c.isReturnError = true
	}
}
