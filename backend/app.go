package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/exchange", getRates).Methods("GET")
	err := http.ListenAndServe(":6000", router)
	if err != nil {
		fmt.Println(err)
	}
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

func getRates(w http.ResponseWriter, r *http.Request) {
	// Connect to MySQL database
	db, err := sql.Open("mysql", "default:ether@tcp(172.18.0.1:3306)/currency-exchange")

	w.Header().Set("Content-Type", "application/json")
	result, err := db.Query("SELECT EUR.date AS date, EUR.price AS EUR, GBP.price AS GBP FROM EUR INNER JOIN GBP ON EUR.date = GBP.date ORDER BY EUR.date DESC;")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var RatesExposed []CurrencyRatesExposed
	for result.Next() {
		var CurrencyRatesExposed CurrencyRatesExposed
		err := result.Scan(&CurrencyRatesExposed.Date, &CurrencyRatesExposed.EUR, &CurrencyRatesExposed.GBP)
		if err != nil {
			panic(err.Error())
		}
		RatesExposed = append(RatesExposed, CurrencyRatesExposed)
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(RatesExposed)

	if err != nil {
		panic(err.Error())
	}
}
