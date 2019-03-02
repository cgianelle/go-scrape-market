package main

import (
	"fmt"
	"net/http"
)

type IndexSummary struct {
	IndexName     string
	LastValue     string
	Change        string
	PercentChange string
}

func main() {
	const MARKET_URL = "https://www.marketwatch.com/tools/marketsummary"
	fmt.Println("Go Scrape Market Starting")

	response, err := http.Get(MARKET_URL)
	if err != nil {
		fmt.Println(err)
	} else {
		defer response.Body.Close()
		table, marketError := findMarketSummaryIndexesTable(response.Body)
		if marketError != nil {
			fmt.Println(marketError)
		} else {
			printTable(table)
		}
	}
}
