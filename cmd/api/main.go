package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"products-api/internal/platform/database"
	"products-api/internal/product"
	"products-api/internal/product/postgres"
	"products-api/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// .env é opcional: em produção as variáveis vêm do ambiente.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("load .env: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL não definida")
	}

	if err := database.MigrateUp(dbURL); err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPool(context.Background(), database.DefaultConfig(dbURL, "products-api"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	service := product.NewService(postgres.NewRepository(pool))
	handler := product.NewHandler(service)
	log.Fatal(http.ListenAndServe(":8080", server.NewRouter(handler)))
}
