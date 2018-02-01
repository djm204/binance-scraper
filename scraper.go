package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

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
	fmt.Println("Exchange Bot v0.03")
	var output, err = GetBinanceInfo("http://api.binance.com/api/v1/ticker/24hr")
	coinsToPostAbout := []string{
		"ETHBTC",
		"LTCBTC",
		"XLMBTC",
		"XRPBTC",
		"TRXBTC",
		"NEOBTC",
		"VENBTC",
		"BNBBTC",
		"OMGBTC",
		"BCCBTC",
		"DASHBTC",
		"POWRBTC",
		"ZECBTC",
		"ADABTC",
		"ICXBTC",
		"IOSTBTC",
		"NAVBTC",
		"KMDBTC",
	}

	coinsThatWillBePosted := []RateInfo{}

	if err != nil {
		fmt.Println("error:", err)
	}

	for k := 0; k < len(coinsToPostAbout)-1; k++ {
		for x := 0; x < len(output)-1; x++ {
			if coinsToPostAbout[k] == output[x].Symbol {
				coinsThatWillBePosted = append(coinsThatWillBePosted, output[x])
			}
		}
	}

	for true {
		PostToTwitter(coinsThatWillBePosted)
	}

}

//PostToTwitter posts to twitter
func PostToTwitter(coinsThatWillBePosted []RateInfo) {
	myClient := Configure()

	for i := 0; i < len(coinsThatWillBePosted)-1; i++ {
		thing := coinsThatWillBePosted[i]
		// if thing.Symbol == "123456" || thing.LastPrice == "0.00000000" || thing.PriceChangePercent == "0" {
		// 	continue
		// }
		fmt.Println("Tweet #:", i)
		var tweetText = "Symbol: $" + thing.Symbol + "\n" +
			"Price Change Percentage: " + thing.PriceChangePercent + "%\n" +
			"Last Price: " + thing.LastPrice + "\n" +
			"Volume: " + thing.Volume + "\n" +
			"Available On: " + "https://www.binance.com/?ref=19085583 \n" +
			PrintHashTags(thing.Symbol)

		fmt.Println(tweetText + "\n")

		Tweet(tweetText, myClient)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n\n\n*******Completed at: " + time.Now().Format("2006-01-02 15:04:05") + "*******\n\n\n")

	time.Sleep(5 * time.Second)
}

//PrintHashTags prints hash tags
func PrintHashTags(symbol string) string {
	return "#" + symbol + " #Binance" + " $BTC $" + strings.Replace(symbol, "BTC", "", -1)
}

//Configure Twitter API
func Configure() *twitter.Client {
	//Load ENV File
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("TWITTER_CONSUMER_API_KEY")
	apiSecret := os.Getenv("TWITTER_CONSUMER_API_SECRET")
	accountToken := os.Getenv("TWITTER_ACCOUNT_ACCESS_TOKEN")
	accountSecret := os.Getenv("TWITTER_ACCOUNT_ACCESS_SECRET")
	config := oauth1.NewConfig(apiKey, apiSecret)
	token := oauth1.NewToken(accountToken, accountSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return client
}

//Tweet a tweet
func Tweet(tweetText string, client *twitter.Client) bool {
	tweet, resp, err := client.Statuses.Update(tweetText, nil)

	fmt.Println("Tweet Status: ", resp.Status)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(tweet, "Error tweeting")
		return false
	}

	return true
}

//GetBinanceInfo from binance
func GetBinanceInfo(url string) ([]RateInfo, error) {
	var myClient = &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-agent", "crypto-reports-anyalizer")

	resp, err := myClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err, resp.Status)
	}

	body, readError := ioutil.ReadAll(resp.Body)
	if readError != nil {
		panic(readError)
	}

	rI, err := GetRates([]byte(body))
	if err != nil {
		panic(err)
	}
	return rI, err
}

//GetRates gets rates from server
func GetRates(body []byte) ([]RateInfo, error) {
	var bI []RateInfo
	err := json.Unmarshal(body, &bI)
	if err != nil {
		fmt.Println("error:", err)
	}
	return bI, err
}
