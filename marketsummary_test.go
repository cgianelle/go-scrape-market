package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func openHTML(file string) (io.ReadCloser, error) {
	htmlFile, err := os.Open(file)
	if err != nil {
		return nil, err
	} else {
		reader := ioutil.NopCloser(htmlFile)
		return reader, nil
	}
}

func TestOpen(t *testing.T) {
	const file = "marketwatch.html"
	_, error := openHTML(file)
	if error != nil {
		t.Error(error)
	}
}

func TestParseGoodMarketWatch(t *testing.T) {
	reader, error := openHTML("marketwatch.html")
	if error != nil {
		t.Error(error)
	} else {
		table, msError := findMarketSummaryIndexesTable(reader)
		if msError != nil {
			t.Error(msError)
		} else {
			if len(table) != 12 {
				t.Errorf("Expected table length to be 12, got %d", len(table))
			}
			first := table[0]
			last := table[len(table)-1]
			if first.IndexName != "Dow Jones Industrial Average" &&
				first.LastValue != "26,026" && first.Change != "+110.32" &&
				first.PercentChange != "+0.43%" {
				t.Errorf("Expected first index to be Dow Jones Industrial Average, but got %v", first)
			}
			if last.IndexName != "CBOE 10 Year Treasury Note..." &&
				last.LastValue != "27.55" && last.Change != "+0.44" &&
				last.PercentChange != "+1.62%" {
				t.Errorf("Expected first index to be CBOE 10 Year Treasury Note..., but got %v", last)
			}
		}
	}
}

func TestParseNoTableHTML(t *testing.T) {
	reader, error := openHTML("notable.html")
	if error != nil {
		t.Error(error)
	} else {
		_, msError := findMarketSummaryIndexesTable(reader)
		if msError == nil {
			t.Error("Expected this to fail due to html missing tables")
		}
		if msError != nil && msError.Error() != "EOF" {
			t.Errorf("Expected this to fail because of EOF, but got %v", msError.Error())
		}
	}
}

func TestParseMissingMarketSummaryTable(t *testing.T) {
	reader, error := openHTML("marketwatch.1.html")
	if error != nil {
		t.Error(error)
	} else {
		_, msError := findMarketSummaryIndexesTable(reader)
		if msError == nil {
			t.Error("Expected this to fail due to html missing table body")
		}
		if msError != nil && msError.Error() != "EOF" {
			t.Errorf("Expected this to fail because of EOF, but got %v", msError.Error())
		}
	}
}

func TestParseMissingTableBody(t *testing.T) {
	reader, error := openHTML("marketwatch.2.html")
	if error != nil {
		t.Error(error)
	} else {
		_, msError := findMarketSummaryIndexesTable(reader)
		if msError == nil {
			t.Error("Expected this to fail due to html missing tables")
		}
		if msError != nil && msError.Error() != "MarketSummaryMissingTableBody" {
			t.Errorf("Expected this to fail because of MarketSummaryMissingTableBody, but got %v", msError.Error())
		}
	}
}
