package storage

import (
	"encoding/json"
	"financial-data-fetcher/internal/domain"
	"os"
)

type Storage interface {
	Save(data []*domain.StockData) error
}

type JSONStorage struct {
	Filename string
}

func NewJSONStorage(filename string) *JSONStorage {
	return &JSONStorage{Filename: filename}
}

func (s *JSONStorage) Save(data []*domain.StockData) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Filename, file, 0644)
}