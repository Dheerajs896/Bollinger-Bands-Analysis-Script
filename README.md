# Bollinger Bands Analysis Script

A command-line tool built in Go to fetch historical candlestick data from the Upstox API, compute **Bollinger Bands**, and generate alerts when the closing price touches the upper or lower bands.

## Features
- Fetch historical candle data (default: 1-minute intervals) for any symbol.
- Computes Bollinger Bands using **20-period Simple Moving Average (SMA)**  **2 standard deviations**.
- Validates trading symbols against suspended instruments and retrieves instrument keys.
- Saves raw candle data to JSON for persistence.
- CLI-driven with flag validation (symbol, from/to dates).
- Configurable via `.env` (Upstox token, timeouts, paths).
- Robust error handling and informative console output.

## Prerequisites
- **Go 1.21+** installed.
- **Upstox Developer API Access**: Get your access token from Upstox Developer Console.
- Instrument data files:
  - `storage/data/instruments.json`: Active instruments (trading_symbol → instrument_key).
  - `storage/data/suspended_instrument.json`: Suspended symbols list.
  (Download from upstox contract note or use provided samples.)

## Quick Start
1. Clone/open the project.
2. Install dependencies:
   ```
   go mod tidy
   ```
3. Create `.env` (see Configuration).
4. Run for example (JOCIL, July 1-31, 2025):
   ```
   go run main.go --symbol=JOCIL --from=01-07-2025 --to=31-07-2025
   ```

## Configuration
Copy the example `.env` structure:

```
# Upstox API
UPSTOX_ACCESS_TOKEN=""
SUSPENDED_INSTRUMENT_FILE_PATH="./storage/data/suspended_instrument.json"
ACTIVE_INSTRUMENT_FILE_PATH="./storage/data/instruments.json"
URL=""
TIME_UNIT="minutes"
TIME_OUT=10

# Instrument files (absolute/relative paths)
SUSPENDED_INSTRUMENT_FILE_PATH=storage/data/suspended_instrument.json
ACTIVE_INSTRUMENT_FILE_PATH=storage/data/instruments.json
```

**Note**: 
Download complete contract note in json formate and unzip and rename it with **instruments.json** and keep it at **storage/data**
If trading symbol contains space then wrap it with double quate for example MCX 2200 CE 26 MAY 26 then **"MCX 2200 CE 26 MAY 26"** at the time of cli inputs.

## Usage
```
go run main.go --symbol=<SYMBOL> --from=<DD-MM-YYYY> --to=<DD-MM-YYYY>
```

### Examples
```
# Analyze JOCIL for July 2025
go run main.go --symbol=JOCIL --from=01-07-2025 --to=31-07-2025

# Reliance, custom interval (edit .env INTERVAL=15)
go run main.go --symbol=RELIANCE --from=15-06-2025 --to=15-07-2025
```

**Help**:
```
go run main.go --help
```

**Input Validation**:
- Dates: DD-MM-YYYY, to > from.
- Symbol: Uppercase, listed, not suspended.
- Env: `UPSTOX_ACCESS_TOKEN` fallback (but use .env TOKEN).

## How It Works
```
1. Load .env config
2. Load instruments: suspended.json & instruments.json into memory
3. Parse CLI: symbol, from/to → API dates (YYYY-MM-DD)
4. Validate: Symbol not suspended? Get instrument_key?
5. Fetch candles: Upstox API (encoded_key, dates, interval)
6. Save: storage/candles/{SYMBOL}_{YYYY-MM_DD}_{YYYY-MM_DD}.json
7. Compute BB:
   - SMA(close, 20)
   - StdDev = √[Σ(close - SMA)² / 20]
   - Upper = SMA + 2*StdDev
   - Lower = SMA - 2*StdDev
8. Alerts: If close >= Upper → \"Upper\"; <= Lower → \"Lower\"
9. Print: [YYYY-MM-DD] Price touched <BAND> Band at <PRICE>
```

**Bollinger Bands Formula**:
```
Middle Band (SMA) = (Sum of last 20 closes) / 20
Upper Band = SMA + (Multiplier × StdDev)
Lower Band = SMA - (Multiplier × StdDev)
```
*(Fixed: period=20, multiplier=2.0)*

## Sample Output
```
Processing Bollinger Bands for JOCIL from 2025-07-01 to 2025-07-31
Validating Trading Symbol from suspended instrument contracts file ...
Fetching Instrument key from contracts file based on trading symbol...
Fetching historical candle data ...
Saving historical candle data to json file ...
Applying Bollinger Bands...
[2025-07-10] Price touched Upper Band at 1850.50
[2025-07-15] Price touched Lower Band at 1720.25
Finished
```

**Saved File Example** (`storage/candles/JOCIL_2025-07-01_2025-07-31.json`):
```json
[
  ["2025-07-01T09:30:00+05:30", 1800.0, 1810.0, 1795.0, 1805.0, 100000],
  ...
]
```

## Project Structure
```
bollinger-bands-script/
├── main.go                 # Entry point
├── go.mod / go.sum        # Dependencies
├── internals/
│   ├── config/            # .env loader
│   ├── domain/            # Models, API client, storage
│   ├── indicators/        # BB calculations
│   ├── service/           # Instrument loaders
│   └── pkg/utils/         # CLI validation
├── storage/
│   ├── data/              # instruments.json, suspended_instrument.json
│   ├── candles/           # Output JSONs (e.g., JOCIL_2025-07-01_2025-07-31.json)
│   └── note.txt
└── README.md              # This file
```

## Limitations & Future Improvements
- Instruments loaded to memory; optimize with Redis (commented in code).

## Troubleshooting
- **No candles**: Check dates, token, instrument_key.
- **Suspended**: Symbol in suspended.json.
- **JSON parse err**: Update instrument files.
- **Token invalid**: Refresh Upstox token.