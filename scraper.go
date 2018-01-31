package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//GetBinanceInfo from binance
func GetBinanceInfo(url string) (interface{}, error) {
	var myClient = &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-agent", "crypto-reports-anyalizer")

	resp, err := myClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		panic(resp)
	}

	body, readError := ioutil.ReadAll(resp.Body)
	if readError != nil {
		panic(readError)
	}

	rI, err := GetRates([]byte(body))
	// binanceJSON := json.Unmarshal(body, target)
	if err != nil {
		panic(err)
	}
	return rI, err
}

//GetRates gets rates from server
func GetRates(body []byte) (interface{}, error) {
	var bI []interface{}
	err := json.Unmarshal(body, &bI)
	if err != nil {
		fmt.Println("error:", err)
	}
	return bI, err
}

//BinanceInfo rateInfo
type BinanceInfo struct {
	Rates []RateInfo `json:"rateInfo"`
}

//RateInfo holds all of the info about a coin
type RateInfo struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int    `json:"openTime"`
	CloseTime          int    `json:"closeTime"`
	FirstID            int    `json:"firstId"`
	LastID             int    `json:"lastId"`
	Count              int    `json:"count"`
}

func main() {
	fmt.Println("Exchange Watcher Bot v0.02")
	fmt.Println(GetBinanceInfo("http://api.binance.com/api/v1/ticker/24hr"))
}
