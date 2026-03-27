package indicators

import (
	"bollinger-bands-script/internals/domain"
	"fmt"
	"math"
)

type BollingerResult struct {
	Timestamp string
	Close     float64
	SMA       float64
	Upper     float64
	Lower     float64
}

type BandAlerts struct {
	Date  string
	Price float64
	Band  string
}

const bollingerBandPeriod = 20
const bollingerBandMultiplier = 2.0

func BollingerBandsCal(candles []domain.RawCandleData) ([]BandAlerts, error) {
	if bollingerBandPeriod <= 0 {
		return nil, fmt.Errorf("bollinger bands period must be a positive integer, got %d", bollingerBandPeriod)
	}
	if len(candles) < bollingerBandPeriod {
		return nil, fmt.Errorf("insuficieant candle data need at least %d candles, have %d", bollingerBandPeriod, len(candles))
	}
	results := make([]BollingerResult, 0, len(candles)-bollingerBandPeriod+1)
	for i := bollingerBandPeriod - 1; i < len(candles); i++ {
		avg := sma(candles, i, bollingerBandPeriod)
		dev := populationStdDev(candles, i, bollingerBandPeriod, avg)

		results = append(results, BollingerResult{
			Timestamp: candles[i].Timestamp,
			Close:     candles[i].Close,
			SMA:       avg,
			Upper:     avg + bollingerBandMultiplier*dev,
			Lower:     avg - bollingerBandMultiplier*dev,
		})
	}
	return FindAlerts(results), nil

}

// calculate the simple moving average of close prices
func sma(candles []domain.RawCandleData, endIdx, bollingerBandPeriod int) float64 {
	sum := 0.0
	for j := endIdx - bollingerBandPeriod + 1; j <= endIdx; j++ {
		sum += candles[j].Close
	}
	return sum / float64(bollingerBandPeriod)
}

// calculate population standard deviation of close prices
func populationStdDev(candles []domain.RawCandleData, endIdx, period int, mean float64) float64 {
	sumSq := 0.0
	for j := endIdx - period + 1; j <= endIdx; j++ {
		diff := candles[j].Close - mean
		sumSq += diff * diff
	}
	return math.Sqrt(sumSq / float64(period))
}

func FindAlerts(results []BollingerResult) []BandAlerts {
	alerts := make([]BandAlerts, 0)
	for _, r := range results {
		date := extractDate(r.Timestamp)
		if r.Close >= r.Upper {
			alerts = append(alerts, BandAlerts{Date: date, Price: r.Close, Band: "Upper"})
		} else if r.Close <= r.Lower {
			alerts = append(alerts, BandAlerts{Date: date, Price: r.Close, Band: "Lower"})
		}
	}
	return alerts
}

func extractDate(timestamp string) string {
	if len(timestamp) >= 10 {
		return timestamp[:10]
	}
	return timestamp
}
