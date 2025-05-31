package main

import (
	"context"
	"financial-data-fetcher/internal/api"
	"financial-data-fetcher/internal/processor"
	"financial-data-fetcher/internal/storage"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	outputFile = "stock_data.json"
	rateLimit  = 13 * time.Second // 5 requests/minute limit (~13 seconds between requests)
)

func main() {
	start := time.Now()

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get API key
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY must be set in environment variables or .env file")
	}

	// Initialize dependencies
	apiClient := api.NewClient(apiKey)
	storage := storage.NewJSONStorage(outputFile)
	processor := processor.NewStockProcessor(apiClient, storage, rateLimit)

	// Symbols to fetch
	symbols := []string{"IBM", "AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA"}

	// Process symbols
	ctx := context.Background()
	data, err := processor.ProcessSymbols(ctx, symbols)
	if err != nil {
		log.Fatalf("Error processing symbols: %v", err)
	}

	log.Printf("Successfully processed %d stocks in %v", len(data), time.Since(start))
}