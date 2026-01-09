package useraggregator

import (
	"log/slog"
	"time"
)

type OptFunc func(u *UserAggregator)

func WithTimeout(timeout time.Duration) OptFunc {
	return func(u *UserAggregator) {
		u.timeout = timeout
	}
}

func WithLogger(logger *slog.Logger) OptFunc {
	return func(u *UserAggregator) {
		u.logger = logger
	}
}
