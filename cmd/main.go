package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/sudesh856/ecom-go-api-project/internal/env"
)

func main() {

	ctx := context.Background()

	cfg := config{
		addr: "8080",
		db:   dbConfig{

			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=2582 dbname=ecom sslmode=disable"),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)


	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to database", "dsn", cfg.db.dsn)



	api := application{
		config: cfg,
		db: conn,
	}

	
	if err := api.run(api.mount()); err != nil {

		slog.Error("server failed to start", "error", err)		
		os.Exit(1)
	}

}