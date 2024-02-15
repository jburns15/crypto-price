package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/fatih/color"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type APIResponse struct {
	ApiResponse []TokenData
}

type TokenData struct {
	ID                                 string  `json:"id"`
	Symbol                             string  `json:"symbol"`
	Name                               string  `json:"name"`
	Image                              string  `json:"image"`
	CurrentPrice                       float64 `json:"current_price"`
	MarketCap                          int64   `json:"market_cap"`
	MarketCapRank                      int     `json:"market_cap_rank"`
	FullyDilutedValuation              int64   `json:"fully_diluted_valuation"`
	TotalVolume                        int64   `json:"total_volume"`
	High24h                            float64 `json:"high_24h"`
	Low24h                             float64 `json:"low_24h"`
	PriceChange24h                     float64 `json:"price_change_24h"`
	PriceChangePercentage24h           float64 `json:"price_change_percentage_24h"`
	MarketCapChange24h                 float64 `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h       float64 `json:"market_cap_change_percentage_24h"`
	CirculatingSupply                  float64 `json:"circulating_supply"`
	TotalSupply                        float64 `json:"total_supply"`
	Ath                                float64 `json:"ath"`
	AthChangePercentage                float64 `json:"ath_change_percentage"`
	AthDate                            string  `json:"ath_date"`
	Atl                                float64 `json:"atl"`
	AtlChangePercentage                float64 `json:"atl_change_percentage"`
	AtlDate                            string  `json:"atl_date"`
	LastUpdated                        string  `json:"last_updated"`
	PriceChangePercentage1hInCurrency  float64 `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24hInCurrency float64 `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7dInCurrency  float64 `json:"price_change_percentage_7d_in_currency"`
}

func main() {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/markets", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("page", "1")
	q.Add("per_page", "100")
	q.Add("vs_currency", "usd")
	q.Add("order", "market_cap_desc")
	q.Add("sparkine", "false")
	q.Add("price_change_percentage", "1h,24h,7d")
	q.Add("local", "en")

	req.Header.Set("Accepts", "application/json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	var tokens []TokenData
	err = json.Unmarshal([]byte(respBody), &tokens)
	if err != nil {
		log.Fatal(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 4, 1, ' ', 0)
	fmt.Fprintf(w, "Name\tCurrentPrice\t1hPercentChange\t24hPercentChange\t7DayPercentChange\tHigh24h\tLow24h\t24hVolume\t\n")
	for _, token := range tokens {
		fmt.Fprintf(w, "%s\t%s\t%f\t%f\t%f\t%s\t%s\t%s\t\n",
			token.Name,
			formatCurrency(token.CurrentPrice),
			token.PriceChangePercentage1hInCurrency,
			token.PriceChangePercentage24h,
			token.PriceChangePercentage7dInCurrency,
			formatCurrency(token.High24h),
			formatCurrency(token.Low24h),
			formatVolume(token.TotalVolume),
		)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func formatCurrency(f float64) string {
	if f > 1.00 {
		printer := message.NewPrinter(language.English)
		return printer.Sprintf("$%.2f", f)
	} else {
		return fmt.Sprintf("$%f", f)
	}
}

func formatVolume(vol int64) string {
	printer := message.NewPrinter(language.English)
	return printer.Sprintf("$%d", vol)
}

func setColor(f float64) string {
	s := strconv.FormatFloat(f, 'f', 4, 64)
	if f == 0 {
		return s
	}
	if f > 0 {
		green := color.New(color.FgGreen)
		return green.Sprintf(s)
	} else {
		red := color.New(color.FgRed)
		return red.Sprintf(s)
	}
}
