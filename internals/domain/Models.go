package domain

import (
	"encoding/json"
	"fmt"
)

type CandleAPIResponse struct {
	Status string     `json:"status"`
	Data   CandleData `json:"data"`
}

type CandleData struct {
	Candles []RawCandleData `json:"candles"`
}

type RawCandleData struct {
	Timestamp    string
	Open         float64
	High         float64
	Low          float64
	Close        float64
	Volume       int64
	OpenInterest float64
}

// just convert the candle data into slice inorder to use the raw candle data
func (c *RawCandleData) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("candle: expected JSON array, got: %w", err)
	}
	if len(raw) < 6 {
		return fmt.Errorf("candle: expected at least 6 elements, got %d", len(raw))
	}
	if err := json.Unmarshal(raw[0], &c.Timestamp); err != nil {
		return fmt.Errorf("candle: parsing timestamp: %w", err)
	}
	if err := json.Unmarshal(raw[1], &c.Open); err != nil {
		return fmt.Errorf("candle: parsing open: %w", err)
	}
	if err := json.Unmarshal(raw[2], &c.High); err != nil {
		return fmt.Errorf("candle: parsing high: %w", err)
	}
	if err := json.Unmarshal(raw[3], &c.Low); err != nil {
		return fmt.Errorf("candle: parsing low: %w", err)
	}
	if err := json.Unmarshal(raw[4], &c.Close); err != nil {
		return fmt.Errorf("candle: parsing close: %w", err)
	}
	if err := json.Unmarshal(raw[5], &c.Volume); err != nil {
		return fmt.Errorf("candle: parsing volume: %w", err)
	}
	if len(raw) > 6 {
		_ = json.Unmarshal(raw[6], &c.OpenInterest)
	}

	return nil
}
