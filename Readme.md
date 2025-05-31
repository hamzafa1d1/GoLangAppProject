# Financial Data Fetcher

A Go application that fetches daily stock market data from Alpha Vantage API and stores it in JSON format.

## Features

- Fetches daily stock data (open, high, low, close, volume) for multiple symbols
- Implements rate limiting to comply with API restrictions (5 requests/minute)
- Uses worker pool pattern for concurrent processing
- Stores collected data in a structured JSON file
- Environment variable configuration
- Modular architecture with clear separation of concerns

## Prerequisites

- Go 1.16 or higher
- Alpha Vantage API key (free tier available)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/financial-data-fetcher.git
   cd financial-data-fetcher
   ```
2. Create a .env file in the root directory with your API key:
  ``` API_KEY=your_alpha_vantage_api_key
  ```

3. Build and Run the APP 
    ```
    go build -o fetcher
    ./fetcher
    ```
## Architecture 
financial-data-fetcher/
├── internal/
│   ├── api/            # API client implementation
│   ├── domain/         # Data models
│   ├── processor/      # Business logic and processing
│   └── storage/        # Data storage interfaces and implementations
└── main.go             # Application entry point