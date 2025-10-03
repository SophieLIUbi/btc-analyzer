package dataloader

import (
	"btc-analyzer/internal/timeseries"
	"btc-analyzer/internal/types"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// LoadFromCoinGecko fetches Bitcoin data from CoinGecko API
func LoadFromCoinGecko(days int) (*types.BTCTimeSeries, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=%d", days)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from CoinGecko: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko API returned status %d", resp.StatusCode)
	}
	
	var coinGeckoResp types.CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&coinGeckoResp); err != nil {
		return nil, fmt.Errorf("failed to decode CoinGecko response: %w", err)
	}
	
	bts := timeseries.New("BTC-USD")
	
	// Convert CoinGecko data to our format
	for i, priceData := range coinGeckoResp.Prices {
		if len(priceData) < 2 {
			continue
		}
		
		timestamp := time.UnixMilli(int64(priceData[0]))
		price := priceData[1]
		
		volume := 0.0
		if i < len(coinGeckoResp.TotalVolumes) && len(coinGeckoResp.TotalVolumes[i]) >= 2 {
			volume = coinGeckoResp.TotalVolumes[i][1]
		}
		
		btcPrice := types.BTCPrice{
			Timestamp: timestamp,
			Open:      price, // CoinGecko doesn't provide OHLC, using price for all
			High:      price,
			Low:       price,
			Close:     price,
			Volume:    volume,
		}
		
		timeseries.AddPrice(bts, btcPrice)
	}
	
	return bts, nil
}

// LoadFromCSV loads Bitcoin data from a CSV file
func LoadFromCSV(filename string) (*types.BTCTimeSeries, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	
	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}
	
	// Determine CSV format based on headers
	headers := records[0]
	format := detectCSVFormat(headers)
	
	bts := timeseries.New("BTC-USD")
	
	for i := 1; i < len(records); i++ {
		record := records[i]
		
		btcPrice, err := parseCSVRecord(record, format)
		if err != nil {
			fmt.Printf("Warning: skipping invalid record at line %d: %v\n", i+1, err)
			continue
		}
		
		timeseries.AddPrice(bts, btcPrice)
	}
	
	return bts, nil
}

// CSVFormat represents different CSV formats
type CSVFormat struct {
	TimestampCol int
	OpenCol      int
	HighCol      int
	LowCol       int
	CloseCol     int
	VolumeCol    int
	TimeFormat   string
}

// detectCSVFormat tries to detect the CSV format based on headers
func detectCSVFormat(headers []string) CSVFormat {
	format := CSVFormat{
		TimestampCol: -1,
		OpenCol:      -1,
		HighCol:      -1,
		LowCol:       -1,
		CloseCol:     -1,
		VolumeCol:    -1,
		TimeFormat:   "2006-01-02", // Default format
	}
	
	for i, header := range headers {
		header = strings.ToLower(strings.TrimSpace(header))
		
		switch {
		case strings.Contains(header, "time") || strings.Contains(header, "date"):
			format.TimestampCol = i
			// Try to detect time format
			if strings.Contains(header, "unix") {
				format.TimeFormat = "unix"
			}
		case strings.Contains(header, "open"):
			format.OpenCol = i
		case strings.Contains(header, "high"):
			format.HighCol = i
		case strings.Contains(header, "low"):
			format.LowCol = i
		case strings.Contains(header, "close") || strings.Contains(header, "price"):
			format.CloseCol = i
		case strings.Contains(header, "volume"):
			format.VolumeCol = i
		}
	}
	
	return format
}

// parseCSVRecord parses a single CSV record based on the detected format
func parseCSVRecord(record []string, format CSVFormat) (types.BTCPrice, error) {
	var btcPrice types.BTCPrice
	
	// Parse timestamp
	if format.TimestampCol >= 0 && format.TimestampCol < len(record) {
		timestampStr := record[format.TimestampCol]
		
		var err error
		if format.TimeFormat == "unix" {
			// Parse Unix timestamp
			timestamp, parseErr := strconv.ParseInt(timestampStr, 10, 64)
			if parseErr != nil {
				return btcPrice, fmt.Errorf("invalid unix timestamp: %w", parseErr)
			}
			btcPrice.Timestamp = time.Unix(timestamp, 0)
		} else {
			// Try common date formats
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"01/02/2006",
				"01/02/2006 15:04:05",
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05.000Z",
			}
			
			for _, timeFormat := range formats {
				btcPrice.Timestamp, err = time.Parse(timeFormat, timestampStr)
				if err == nil {
					break
				}
			}
			
			if err != nil {
				return btcPrice, fmt.Errorf("failed to parse timestamp: %w", err)
			}
		}
	} else {
		return btcPrice, fmt.Errorf("timestamp column not found")
	}
	
	// Helper function to parse float from record
	parseFloat := func(colIndex int, defaultValue float64) float64 {
		if colIndex >= 0 && colIndex < len(record) {
			if val, err := strconv.ParseFloat(record[colIndex], 64); err == nil {
				return val
			}
		}
		return defaultValue
	}
	
	// Parse OHLCV data
	btcPrice.Open = parseFloat(format.OpenCol, 0)
	btcPrice.High = parseFloat(format.HighCol, 0)
	btcPrice.Low = parseFloat(format.LowCol, 0)
	btcPrice.Close = parseFloat(format.CloseCol, 0)
	btcPrice.Volume = parseFloat(format.VolumeCol, 0)
	
	// If OHLC values are missing but we have Close, use Close for all
	if btcPrice.Open == 0 && btcPrice.Close != 0 {
		btcPrice.Open = btcPrice.Close
	}
	if btcPrice.High == 0 && btcPrice.Close != 0 {
		btcPrice.High = btcPrice.Close
	}
	if btcPrice.Low == 0 && btcPrice.Close != 0 {
		btcPrice.Low = btcPrice.Close
	}
	
	return btcPrice, nil
}

// SaveToCSV exports Bitcoin time series data to CSV
func SaveToCSV(bts *types.BTCTimeSeries, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write headers
	headers := []string{"Date", "Open", "High", "Low", "Close", "Volume"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}
	
	// Write data
	timeseries.Sort(bts)
	for _, data := range bts.Data {
		record := []string{
			data.Timestamp.Format("2006-01-02"),
			fmt.Sprintf("%.2f", data.Open),
			fmt.Sprintf("%.2f", data.High),
			fmt.Sprintf("%.2f", data.Low),
			fmt.Sprintf("%.2f", data.Close),
			fmt.Sprintf("%.0f", data.Volume),
		}
		
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}
	
	return nil
}

// SaveToJSON exports Bitcoin time series data to JSON
func SaveToJSON(bts *types.BTCTimeSeries, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(bts); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	
	return nil
}

// LoadFromJSON loads Bitcoin data from a JSON file
func LoadFromJSON(filename string) (*types.BTCTimeSeries, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()
	
	var bts types.BTCTimeSeries
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&bts); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	
	return &bts, nil
}

// GenerateSampleData creates sample Bitcoin data for testing
func GenerateSampleData(days int, startPrice float64) *types.BTCTimeSeries {
	bts := timeseries.New("BTC-USD-SAMPLE")
	
	currentPrice := startPrice
	currentTime := time.Now().AddDate(0, 0, -days)
	
	for i := 0; i < days; i++ {
		// Simple random walk for demo purposes
		change := (float64(i%10) - 4.5) / 100.0 // -4.5% to 4.5% daily change
		
		open := currentPrice
		high := open * (1 + math.Abs(change) + 0.01)
		low := open * (1 - math.Abs(change) - 0.01)
		close := open * (1 + change)
		volume := 1000000.0 + float64(i%100)*10000.0
		
		btcPrice := types.BTCPrice{
			Timestamp: currentTime.AddDate(0, 0, i),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}
		
		timeseries.AddPrice(bts, btcPrice)
		currentPrice = close
	}
	
	return bts
}

// ValidateData performs basic validation on the loaded data
func ValidateData(bts *types.BTCTimeSeries) []string {
	var issues []string
	
	if len(bts.Data) == 0 {
		issues = append(issues, "No data points found")
		return issues
	}
	
	for i, data := range bts.Data {
		// Check for invalid prices
		if data.Open <= 0 || data.High <= 0 || data.Low <= 0 || data.Close <= 0 {
			issues = append(issues, fmt.Sprintf("Invalid price data at index %d", i))
		}
		
		// Check OHLC consistency
		if data.High < data.Low {
			issues = append(issues, fmt.Sprintf("High < Low at index %d", i))
		}
		if data.High < data.Open || data.High < data.Close {
			issues = append(issues, fmt.Sprintf("High is not highest at index %d", i))
		}
		if data.Low > data.Open || data.Low > data.Close {
			issues = append(issues, fmt.Sprintf("Low is not lowest at index %d", i))
		}
		
		// Check for negative volume
		if data.Volume < 0 {
			issues = append(issues, fmt.Sprintf("Negative volume at index %d", i))
		}
		
		// Check for future dates
		if data.Timestamp.After(time.Now()) {
			issues = append(issues, fmt.Sprintf("Future date at index %d", i))
		}
	}
	
	// Check for duplicate timestamps
	timestampMap := make(map[int64]bool)
	for i, data := range bts.Data {
		timestamp := data.Timestamp.Unix()
		if timestampMap[timestamp] {
			issues = append(issues, fmt.Sprintf("Duplicate timestamp at index %d", i))
		}
		timestampMap[timestamp] = true
	}
	
	return issues
}