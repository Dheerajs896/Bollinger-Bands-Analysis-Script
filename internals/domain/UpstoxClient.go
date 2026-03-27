package domain

import (
	"bollinger-bands-script/internals/config"
	"bollinger-bands-script/internals/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func FetchCandles(cfg *config.Config, inpurParams *utils.InputRequestParams, instrumentKey string) ([]RawCandleData, error) {

	encodedInstrumentKey := url.PathEscape(instrumentKey)

	// constract the end point in required format
	endPoint := fmt.Sprintf("%s/%s/%s/%s/%s/%s", cfg.HistoricalCandle.Url, encodedInstrumentKey, cfg.HistoricalCandle.TimeUnit, cfg.HistoricalCandle.Interval, inpurParams.ToDate, inpurParams.FromDate)

	// add time out in api call
	client := &http.Client{Timeout: time.Duration(cfg.HistoricalCandle.Timeout) * time.Second}

	req, err := http.NewRequest(http.MethodGet, endPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("candle request body: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.HistoricalCandle.Token)
	req.Header.Set("Accept", "application/json")

	// call the api
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching historical candles data : %w", err)
	}
	defer resp.Body.Close()

	// read the response body of the api
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body: %w", err)
	}

	//check status code of the api response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching historical candles: API returned HTTP %d: %s",
			resp.StatusCode, body)
	}

	// map the response data into candle api response struc
	var apiResp CandleAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("decoding candle response: %w", err)
	}

	//double check that api have given success with required data points
	candles := apiResp.Data.Candles
	if len(candles) == 0 {
		return nil, fmt.Errorf("no candle data returned for the requested date range")
	}

	return candles, nil

}
