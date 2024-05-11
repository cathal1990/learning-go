package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type tickerInfo struct {
	A []string `json:"a"`
	B []string `json:"b"`
	C []string `json:"c"`
	V []string `json:"v"`
	P []string `json:"p"`
	T []int    `json:"t"`
	L []string `json:"l"`
	H []string `json:"h"`
	O string   `json:"o"`
}

type tickerQuery struct {
	Error  []string              `json:"error"`
	Result map[string]tickerInfo `json:"result"`

var urlTemplate = "https://api.kraken.com/0/public/Ticker?pair=%s"

func main() {
	var ticker string

	flag.StringVar(&ticker, "ticker", "", "Crypto ticker to query.")
	flag.Parse()

	if ticker == "" {
		fmt.Println("Invalid or no ticker provided")
		os.Exit(1)
	}

	url := fmt.Sprintf(urlTemplate, ticker)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error with request: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	tickerStats := &tickerQuery{}
	err = json.Unmarshal(body, tickerStats)

	if err != nil {
		fmt.Printf("Failed to unmarshal JSON body from response: %v\n", err)
		os.Exit(1)
	}

	for _, info := range tickerStats.Result {
		fmt.Printf("-----------Ticker Statistics for %s------÷------\n", ticker)
		fmt.Printf("Ask: %s\nBid: %s\nPrice: %s\nVolume: %s\nAverage price (24h): %s\nHigh (24h): %s\nLow (24h): %s\n",
			info.A[0], info.B[0], info.C[0], info.V[0], info.P[0], info.H[0], info.L[0])
	}
}
