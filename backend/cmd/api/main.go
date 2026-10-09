package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"

	httpapi "spendtrack/internal/adapters/http"
	"spendtrack/internal/adapters/postgres"
	"spendtrack/internal/application"
)

type config struct {
	Port        int    `env:"PORT" envDefault:"8000"`
	DatabaseURL string `env:"DATABASE_URL,required"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := postgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	svc := application.NewExpenseService(postgres.NewExpenseRepository(pool))
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           httpapi.NewRouter(svc),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%d", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
