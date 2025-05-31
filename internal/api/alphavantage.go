package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"financial-data-fetcher/internal/domain"
)

const (
	baseURL = "https://www.alphavantage.co/query"
)

type AlphaVantageClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *AlphaVantageClient {
	return &AlphaVantageClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *AlphaVantageClient) FetchDailyStockData(symbol string) (*domain.StockData, error) {
	url := fmt.Sprintf("%s?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s",
		baseURL, symbol, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if errMsg, ok := apiResponse["Error Message"].(string); ok {
		return nil, fmt.Errorf("API error: %s", errMsg)
	}

	metaData, ok := apiResponse["Meta Data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid meta data format")
	}

	timeSeries, ok := apiResponse["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid time series format")
	}

	data := &domain.StockData{
		Symbol:        symbol,
		LastRefreshed: metaData["3. Last Refreshed"].(string),
		TimeSeries:    make(map[string]domain.DailyData),
	}

	for date, values := range timeSeries {
		valueMap, ok := values.(map[string]interface{})
		if !ok {
			continue
		}

		open, _ := strconv.ParseFloat(valueMap["1. open"].(string), 64)
		high, _ := strconv.ParseFloat(valueMap["2. high"].(string), 64)
		low, _ := strconv.ParseFloat(valueMap["3. low"].(string), 64)
		closeVal, _ := strconv.ParseFloat(valueMap["4. close"].(string), 64)
		volume, _ := strconv.Atoi(valueMap["5. volume"].(string))

		data.TimeSeries[date] = domain.DailyData{
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closeVal,
			Volume: volume,
		}
	}

	return data, nil
}
