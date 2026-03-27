package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Instrument       InstrumentConfig
	HistoricalCandle UpstockHistoricalCandleConfig
}

type InstrumentConfig struct {
	Sus_Instrument_File_Path string
	Act_Instrument_File_Path string
}

type UpstockHistoricalCandleConfig struct {
	Url      string
	Token    string
	Interval string
	Timeout  int
	TimeUnit string
}

// Load configuration from .env file and enviroment veriables
func Load() *Config {

	// check env file exist
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, using environment variables")
	}

	return &Config{
		Instrument: InstrumentConfig{
			Sus_Instrument_File_Path: getEnv("SUSPENDED_INSTRUMENT_FILE_PATH", ""),
			Act_Instrument_File_Path: getEnv("ACTIVE_INSTRUMENT_FILE_PATH", ""),
		},
		HistoricalCandle: UpstockHistoricalCandleConfig{
			Url:      getEnv("URL", "https://api.upstox.com/v3/historical-candle"),
			Token:    getEnv("TOKEN", ""),
			Interval: getEnv("INTERVAL", "1"),
			TimeUnit: getEnv("TIME_UNIT", "minutes"),
			Timeout:  getEnvAsInt("TIME_OUT", 10),
		},
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
