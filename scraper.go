package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type IndexSummary struct {
	IndexName     string
	LastValue     string
	Change        string
	PercentChange string
}

func printTable(table []IndexSummary) {
	for _, s := range table {
		fmt.Printf("%+v\n", s)
	}
}

func processCell(token *html.Tokenizer) string {
	foundTD := false
	var text string
	for {
		tt := token.Next()
		switch {
		case tt == html.ErrorToken:
			return ""
		case tt == html.StartTagToken:
			t := token.Token()
			if t.Data == "td" {
				// fmt.Println("Start Of Cell")
				foundTD = true
			}
		case tt == html.TextToken:
			if foundTD {
				text = string(token.Text())
				foundTD = false
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "td" {
				// fmt.Println("End Of Cell")
				return text
			}
		}
	}
}

func processIndexLinkedCell(token *html.Tokenizer) string {
	foundA := false
	var text string
	for {
		tt := token.Next()
		switch {
		case tt == html.ErrorToken:
			return ""
		case tt == html.StartTagToken:
			t := token.Token()
			if t.Data == "a" {
				// fmt.Println("Start Of Cell")
				foundA = true
			}
		case tt == html.TextToken:
			if foundA {
				text = string(token.Text())
				foundA = false
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "td" {
				// fmt.Println("End Of Cell")
				return text
			}
		}
	}
}

func processIndexRows(token *html.Tokenizer) []IndexSummary {
	var indexes []IndexSummary
	// processes row between the <TR> and the </TR> tags
	// format:
	// <td>
	//		<a>Index Title</a>
	// 		<span> </span>
	// </td>
	// <td>LAST</td>
	// <td>CHANGE</td>
	// <td>% Change</td>
	for {
		tt := token.Next()
		switch {
		case tt == html.ErrorToken:
			return nil
		case tt == html.StartTagToken:
			t := token.Token()
			// fmt.Println(t.Data)
			if t.Data == "tr" {
				var index IndexSummary
				index.IndexName = processIndexLinkedCell(token)
				index.LastValue = processCell(token)
				index.Change = processCell(token)
				index.PercentChange = processCell(token)
				indexes = append(indexes, index)
			}
		case tt == html.EndTagToken:
			t := token.Token()
			// fmt.Println(t.Data)
			if t.Data == "tbody" {
				// fmt.Println("End Of table")
				return indexes
			}
		}
	}
}

func parseMarketSummaryIndexesTable(token *html.Tokenizer) []IndexSummary {
	var table []IndexSummary
	for {
		tt := token.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			fmt.Println("Error or End of Document")
			return nil
		case tt == html.StartTagToken:
			t := token.Token()
			fmt.Print(t.Data)
			fmt.Print(" ")
			if t.Data == "tbody" {
				table = processIndexRows(token)
			}
		case tt == html.EndTagToken:
			t := token.Token()
			fmt.Println(t.Data)
			if t.Data == "table" {
				fmt.Println("Found end of table, exiting")
				return table
			}
		}
	}
}

func findMarketSummaryIndexesTable(body io.ReadCloser) {
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "table"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "id" && a.Val == "marketsummaryindexes" {
						fmt.Println("Found id:", a.Val)
						fmt.Println(t)
						table := parseMarketSummaryIndexesTable(z)
						printTable(table)
					}
				}
			}
		}
	}
}

func main() {
	const MARKET_URL = "https://www.marketwatch.com/tools/marketsummary"
	fmt.Println("Go Scrape Market Starting")

	response, err := http.Get(MARKET_URL)
	if err != nil {
		fmt.Println(err)
	} else {
		defer response.Body.Close()
		findMarketSummaryIndexesTable(response.Body)
	}
}
