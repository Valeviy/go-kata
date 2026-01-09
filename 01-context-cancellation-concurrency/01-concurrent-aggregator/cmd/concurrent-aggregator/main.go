package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"concurrent-aggregator/internal/client/orderservice"
	"concurrent-aggregator/internal/client/profileservice"
	"concurrent-aggregator/internal/service/useraggregator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// TODO: move to tests
	for i, userAggregator := range []*useraggregator.UserAggregator{
		// "Slow Poke" case
		useraggregator.New(
			orderservice.New(
				orderservice.WithTimeout(time.Second),
				orderservice.WithDelay(time.Second*10),
			),
			profileservice.New(
				profileservice.WithTimeout(time.Second),
			),
			useraggregator.WithTimeout(time.Second),
			useraggregator.WithLogger(logger),
		),
		// "Domino Effect" case
		useraggregator.New(
			orderservice.New(
				orderservice.WithTimeout(time.Second),
				orderservice.WithDelay(time.Second*10),
			),
			profileservice.New(
				profileservice.WithTimeout(time.Second),
				profileservice.WithReturnError(),
			),
			useraggregator.WithTimeout(time.Second),
			useraggregator.WithLogger(logger),
		),
		// Valid case
		useraggregator.New(
			orderservice.New(
				orderservice.WithTimeout(time.Second),
				orderservice.WithDelay(time.Second),
			),
			profileservice.New(
				profileservice.WithTimeout(time.Second),
			),
			useraggregator.WithTimeout(time.Second*2),
			useraggregator.WithLogger(logger),
		),
	} {
		info, err := userAggregator.Aggregate(ctx, 15)
		if err != nil {
			logger.Error("cannot aggregate info", "case", i, "error", err)

			continue
		}

		logger.Info("aggregated info", "case", i, "info", info)
	}
}
