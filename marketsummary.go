package main

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

func printTable(table []IndexSummary) {
	for _, s := range table {
		fmt.Printf("%+v\n", s)
	}
}

func processCell(token *html.Tokenizer, startTag string) (string, error) {
	foundStart := false
	var text string
	for {
		tt := token.Next()
		switch {
		case tt == html.ErrorToken:
			return "", token.Err()
		case tt == html.StartTagToken:
			t := token.Token()
			if t.Data == startTag {
				// fmt.Println("Start Of Cell")
				foundStart = true
			}
		case tt == html.TextToken:
			if foundStart {
				text = string(token.Text())
				foundStart = false
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "td" {
				// fmt.Println("End Of Cell")
				return text, nil
			}
		}
	}
}

func processIndexRows(token *html.Tokenizer) ([]IndexSummary, error) {
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
			return nil, token.Err()
		case tt == html.StartTagToken:
			t := token.Token()
			// fmt.Println(t.Data)
			if t.Data == "tr" {
				var index IndexSummary
				index.IndexName, _ = processCell(token, "a")
				index.LastValue, _ = processCell(token, "td")
				index.Change, _ = processCell(token, "td")
				index.PercentChange, _ = processCell(token, "td")
				indexes = append(indexes, index)
			}
		case tt == html.EndTagToken:
			t := token.Token()
			// fmt.Println(t.Data)
			if t.Data == "tbody" {
				// fmt.Println("End Of table")
				return indexes, nil
			}
		}
	}
}

func parseMarketSummaryIndexesTable(token *html.Tokenizer) ([]IndexSummary, error) {
	// var table []IndexSummary
	for {
		tt := token.Next()

		switch {
		case tt == html.ErrorToken:
			fmt.Println("parseMarketSummaryIndexesTable")
			return nil, token.Err()
		case tt == html.StartTagToken:
			t := token.Token()
			if t.Data == "tbody" {
				return processIndexRows(token)
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "table" {
				//--If we hit the end of the table before the body return
				return nil, errors.New("MarketSummaryMissingTableBody")
			}
		}
	}
}

func findMarketSummaryIndexesTable(body io.ReadCloser) ([]IndexSummary, error) {
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return nil, z.Err()
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "table"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "id" && a.Val == "marketsummaryindexes" {
						return parseMarketSummaryIndexesTable(z)
					}
				}
			}
		}
	}
}
