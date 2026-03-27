package utils

import (
	"bollinger-bands-script/internals/service"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const dateInputLayout = "02-01-2006"
const dateAPILayout = "2006-01-02"

type InputRequestParams struct {
	Symbol   string
	FromDate string
	ToDate   string
}

func ValidateInputs() (*InputRequestParams, error) {
	symbol := flag.String("symbol", "", "equity symbol, e.g. INFY")
	from := flag.String("from", "", "start date in DD-MM-YYYY format")
	to := flag.String("to", "", "end date in DD-MM-YYYY format")
	flag.Parse()

	if *symbol == "" {
		return nil, fmt.Errorf("required flag --symbol is missing\nUsage: go run main.go --symbol=INFY --from=01-07-2025 --to=31-07-2025")
	}
	if *from == "" {
		return nil, fmt.Errorf("required flag --from is missing\nUsage: go run main.go --symbol=INFY --from=01-07-2025 --to=31-07-2025")
	}
	if *to == "" {
		return nil, fmt.Errorf("required flag --to is missing\nUsage: go run main.go --symbol=INFY --from=01-07-2025 --to=31-07-2025")
	}

	fromTime, err := time.Parse(dateInputLayout, *from)
	if err != nil {
		return nil, fmt.Errorf("invalid --from date %q: expected DD-MM-YYYY format", *from)
	}

	toTime, err := time.Parse(dateInputLayout, *to)
	if err != nil {
		return nil, fmt.Errorf("invalid --to date %q: expected DD-MM-YYYY format", *to)
	}

	if !toTime.After(fromTime) {
		return nil, fmt.Errorf("--to date (%s) must be after --from date (%s)", *to, *from)
	}

	token := os.Getenv("UPSTOX_ACCESS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("environment variable UPSTOX_ACCESS_TOKEN is not set")
	}

	return &InputRequestParams{
		Symbol:   *symbol,
		FromDate: fromTime.Format(dateAPILayout),
		ToDate:   toTime.Format(dateAPILayout),
	}, nil
}

func IsInstrumentSuspended(symbol string, suspendedInstrument []service.SuspendedInstrument) bool {

	for _, v := range suspendedInstrument {
		if strings.EqualFold(v.TradingSymbol, service.FormateTradingSymbol(symbol)) {
			return true
		}
	}
	return false
}

func GetInstrumentKeyFromTradingSymbol(symbol string, activeInstrument []service.ActiveInstrument) string {
	for _, v := range activeInstrument {
		if strings.EqualFold(v.TradingSymbol, service.FormateTradingSymbol(symbol)) {
			return v.InstrumentKey
		}
	}
	return ""
}
