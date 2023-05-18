package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	askRates()
	for range time.Tick(time.Hour * 10) {
        askRates()
    }
}
func askRates() {
	currentTime := time.Now().Format("2006-01-02")
	tenDaysAgo := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	test := retrieveExchangeRates(tenDaysAgo, currentTime)
	fmt.Println(test)
	persistExchangeRates(test)
}

type RateMap map[string]float64

type CurrencyExchangeResponse struct {
	Base       string             `json:"base"`
	StartDate  string             `json:"start_date"`
	EndDate    string             `json:"end_date"`
	Success    bool               `json:"success"`
	TimeSeries bool               `json:"timeseries"`
	Rates      map[string]RateMap `json:"rates"`
}
type CurrencyRatesExposed struct {
	Date string `json:"date"`
	EUR  string `json:"EUR"`
	GBP  string `json:"GBP"`
}

func retrieveExchangeRates(startDate, endDate string) map[string]RateMap {
	//the api request code
	url := fmt.Sprintf("https://api.apilayer.com/fixer/timeseries?base=USD&symbols=EUR,GBP&start_date=%s&end_date=%s", startDate, endDate)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", "cseAjuQ019VT2UC2FcZDwXABTxPlceJw")

	if err != nil {
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	//parse json
	var exchangeRates CurrencyExchangeResponse
	err = json.Unmarshal(body, &exchangeRates)
	if err != nil {
		return nil
	}
	return exchangeRates.Rates
}

func persistExchangeRates(rates map[string]RateMap) error {
	// Connect to MySQL database
	db, err := sql.Open("mysql", "default:ether@tcp(172.18.0.1:3306)/currency-exchange")
	if err != nil {
		return err
	}
	defer db.Close()

	// Iterate over each currency in the rates map
	for date, ratesMap := range rates {
		for currency, rate := range ratesMap {
			// Create table if it doesn't exist
			_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (date DATE UNIQUE, price FLOAT NOT NULL)", currency))
			if err != nil {
				fmt.Println(err)
			}

			_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (date, price) VALUES (?, ?)", currency), date, rate)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
