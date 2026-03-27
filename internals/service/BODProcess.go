package service

import (
	"bollinger-bands-script/internals/config"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type SuspendedInstrument struct {
	TradingSymbol string `json:"trading_symbol"`
}

func LoadSuspendedInstrument(cfg config.InstrumentConfig) ([]SuspendedInstrument, error) {
	rawSuspendedTradingSymbol, err := os.ReadFile(cfg.Sus_Instrument_File_Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read suspended instrument file: %w", err)
	}

	var suspendedInstrument []SuspendedInstrument
	if err := json.Unmarshal(rawSuspendedTradingSymbol, &suspendedInstrument); err != nil {
		return nil, fmt.Errorf("failed to parse suspended instrument JSON: %w", err)
	}
	for i := range suspendedInstrument {
		suspendedInstrument[i].TradingSymbol =
			FormateTradingSymbol(suspendedInstrument[i].TradingSymbol)
	}

	return suspendedInstrument, nil
}

type ActiveInstrument struct {
	TradingSymbol string `json:"trading_symbol"`
	InstrumentKey string `json:"instrument_key"`
}

func LoadActiveInstrument(cfg config.InstrumentConfig) ([]ActiveInstrument, error) {
	rawSuspendedTradingSymbol, err := os.ReadFile(cfg.Act_Instrument_File_Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read active instrument file: %w", err)
	}

	var activeInstrument []ActiveInstrument
	if err := json.Unmarshal(rawSuspendedTradingSymbol, &activeInstrument); err != nil {
		return nil, fmt.Errorf("failed to parse active instrument JSON: %w", err)
	}
	for i := range activeInstrument {
		activeInstrument[i].TradingSymbol =
			FormateTradingSymbol(activeInstrument[i].TradingSymbol)
	}
	return activeInstrument, nil
}

// Replace spaces with underscores and convert the value to uppercase to ensure consistency,
// since the trading symbol will be entered by end users and may vary in casing or format.
func FormateTradingSymbol(symbol string) string {
	fields := strings.Fields(strings.TrimSpace(symbol))
	return strings.ToUpper(strings.Join(fields, "_"))
}
