package main

import (
	"bollinger-bands-script/internals/config"
	"bollinger-bands-script/internals/domain"
	"bollinger-bands-script/internals/indicators"
	"bollinger-bands-script/internals/pkg/utils"
	"bollinger-bands-script/internals/service"
	"fmt"
	"os"
)

func main() {
	//Load config from .env
	cfg := config.Load()

	// Load Data in Memory
	// Load suspended and active instrument in memory
	// we can load this data in redis inorder to remove application server memory uses
	suspendedInstrument, loadSusErr := service.LoadSuspendedInstrument(cfg.Instrument)
	activeInstrument, loadActErr := service.LoadActiveInstrument(cfg.Instrument)

	// If any file not loaded properly terminate the application..
	if loadSusErr != nil || loadActErr != nil {
		fmt.Println("Error in loading BOD files..........")
		fmt.Println(loadSusErr, " ", loadActErr)
		os.Exit(0)
	}

	// Validate Input fields
	inpurParams, validationError := utils.ValidateInputs()
	if validationError != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", validationError)
		os.Exit(1)
	}

	fmt.Printf("Processing Bollinger Bands for %s from %s to %s\n", inpurParams.Symbol, inpurParams.FromDate, inpurParams.ToDate)

	//Check Input instrument is valid as per the contract note
	fmt.Println("Validating Trading Symbol from suspended instrument contracts file ...")
	if utils.IsInstrumentSuspended(inpurParams.Symbol, suspendedInstrument) {
		fmt.Println("Input Symbol is suspended for trading...")
		os.Exit(0)
	}

	//Get The Instrument key using trading symbol from contract note
	fmt.Println("Fetching Instrument key from contracts file based on trading symbol...")
	instrumentKey := utils.GetInstrumentKeyFromTradingSymbol(inpurParams.Symbol, activeInstrument)
	if instrumentKey == "" {
		fmt.Println("Input Symbol is not present in contract note. Kindly verify the symbol")
		os.Exit(0)
	}

	//fetch historical candle data using upstox api
	fmt.Println("Fetching historical candle data ...")
	candles, err := domain.FetchCandles(cfg, inpurParams, instrumentKey)
	if err != nil {
		fmt.Println("Error while calling upstox historical candle api : ", err)
		os.Exit(0)
	}

	fmt.Println("Saving historical candle data to json file ...")
	domain.Save(candles, inpurParams)

	// implemetion of bollinger band calculations
	fmt.Println("Applying Bollinger Bands...")
	indicatorPriceAlert, err := indicators.BollingerBandsCal(candles)
	if err != nil {
		fmt.Println("BollingerBandsCal Error:- ", err)
		os.Exit(0)
	}

	// Print all alerts of bollinger band which touches the bands of indecators
	for _, a := range indicatorPriceAlert {
		fmt.Printf("[%s] Price touched %s Band at %.2f\n", a.Date, a.Band, a.Price)
	}

	// Process completed so give the ack...
	fmt.Println("Finished")
}
