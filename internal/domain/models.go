package domain

type StockData struct {
	Symbol       string             `json:"symbol"`
	LastRefreshed string           `json:"last_refreshed"`
	TimeSeries   map[string]DailyData `json:"time_series"`
}

type DailyData struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int     `json:"volume"`
}