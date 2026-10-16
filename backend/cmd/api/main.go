package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	httpapi "spendtrack/internal/adapters/http"
	"spendtrack/internal/adapters/postgres"
	"spendtrack/internal/application"
	"spendtrack/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
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
