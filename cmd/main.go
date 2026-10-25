package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/DesirArman007/ecom/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {

	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "user=postgres password=postgres host=localhost dbname=ecom sslmode=disable"),
		},
	}

	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	defer conn.Close(ctx)
	logger.Info("Connected to the database successfully", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	// api.run(api.mount())
	handler := api.mount()
	if err := api.run(handler); err != nil {
		slog.Error("Server failed to start", "Error", err)
		os.Exit(1)
	}
}
