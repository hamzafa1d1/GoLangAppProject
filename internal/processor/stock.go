package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"financial-data-fetcher/internal/api"
	"financial-data-fetcher/internal/domain"
	"financial-data-fetcher/internal/storage"
)

type StockProcessor struct {
	apiClient *api.AlphaVantageClient
	storage   storage.Storage
	rateLimit time.Duration
}

func NewStockProcessor(apiClient *api.AlphaVantageClient, storage storage.Storage, rateLimit time.Duration) *StockProcessor {
	return &StockProcessor{
		apiClient: apiClient,
		storage:   storage,
		rateLimit: rateLimit,
	}
}

func (p *StockProcessor) ProcessSymbols(ctx context.Context, symbols []string) ([]*domain.StockData, error) {
	var (
		wg      sync.WaitGroup
		results = make(chan *domain.StockData, len(symbols))
		errors  = make(chan error, len(symbols))
	)

	// Worker pool
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, symbols, results, errors)
	}

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	var allData []*domain.StockData
	for result := range results {
		allData = append(allData, result)
	}

	// Handle errors
	for err := range errors {
		return nil, err
	}

	// Save data
	if err := p.storage.Save(allData); err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	return allData, nil
}

func (p *StockProcessor) worker(ctx context.Context, wg *sync.WaitGroup, symbols []string, results chan<- *domain.StockData, errors chan<- error) {
	defer wg.Done()

	for _, symbol := range symbols {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := p.apiClient.FetchDailyStockData(symbol)
			if err != nil {
				errors <- fmt.Errorf("failed to process %s: %w", symbol, err)
				continue
			}
			results <- data
			time.Sleep(p.rateLimit)
		}
	}
}
