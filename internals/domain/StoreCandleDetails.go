package domain

import (
	"bollinger-bands-script/internals/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const dataDir = "storage/candles"

func Save(candles []RawCandleData, InputRequestParams *utils.InputRequestParams) (string, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("create candle dir inside storage folder: %w", err)
	}
	filename := fmt.Sprintf("%s_%s_%s.json", InputRequestParams.Symbol, InputRequestParams.FromDate, InputRequestParams.ToDate)
	path := filepath.Join(dataDir, filename)

	payload := struct {
		Symbol  string          `json:"symbol"`
		From    string          `json:"from"`
		To      string          `json:"to"`
		Total   int             `json:"total_candles"`
		Candles []RawCandleData `json:"candles"`
	}{
		Symbol:  InputRequestParams.Symbol,
		From:    InputRequestParams.FromDate,
		To:      InputRequestParams.ToDate,
		Total:   len(candles),
		Candles: candles,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error in marshal indentcandle data : %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("writing candle data in file %s: %w", path, err)
	}
	return path, nil
}
