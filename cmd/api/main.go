package main

import (
	"context"
	"log"
	"net/http"
	"os"
	product "products-api/internal/application/product"
	httpadapter "products-api/internal/infrastructure/http"
	"products-api/internal/infrastructure/persistence/config"
	"products-api/internal/infrastructure/persistence/migration"
	postgres "products-api/internal/infrastructure/persistence/postgres"

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

	if err := migration.Up(dbURL); err != nil {
		log.Fatal(err)
	}

	pool, err := config.NewPool(context.Background(), config.DefaultConfig(dbURL, "products-api"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	productService := product.NewProductService(postgres.NewPostgresProductRepository(pool))
	productHandler := httpadapter.NewProductHandler(productService)
	log.Fatal(http.ListenAndServe(":8080", httpadapter.NewRouter(productHandler)))
}
