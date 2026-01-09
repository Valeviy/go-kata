package useraggregator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	defaultTimeout = time.Second
)

type OrderService interface {
	GetOrdersByUserID(context.Context, uint64) (string, error)
}

type ProfileService interface {
	GetProfileByUserID(context.Context, uint64) (string, error)
}

type UserAggregator struct {
	orderService   OrderService
	profileService ProfileService

	timeout time.Duration
	logger  *slog.Logger
}

func New(
	orderService OrderService,
	profileService ProfileService,
	opts ...OptFunc,
) *UserAggregator {
	u := &UserAggregator{
		orderService:   orderService,
		profileService: profileService,
		timeout:        defaultTimeout,
		logger:         slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func (a *UserAggregator) Aggregate(ctx context.Context, userID uint64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	var (
		ordersInfo  string
		profileInfo string
	)

	g.Go(func() error {
		var err error
		ordersInfo, err = a.orderService.GetOrdersByUserID(ctx, userID)
		if err != nil {
			a.logger.Error("cannot get orders info", "error", err.Error())

			return errors.Wrap(err, "orderService.GetOrdersByUserID")
		}

		return nil
	})

	g.Go(func() error {
		var err error
		profileInfo, err = a.profileService.GetProfileByUserID(ctx, userID)
		if err != nil {
			a.logger.Error("cannot get profile info", "error", err.Error())

			return errors.Wrap(err, "profileService.GetProfileByUserID")
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return "", errors.Wrap(err, "g.Wait")
	}

	return fmt.Sprintf("%s | %s", profileInfo, ordersInfo), nil
}
