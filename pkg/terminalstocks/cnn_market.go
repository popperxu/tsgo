// Copyright (c) 2013-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package TerminalStocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const marketURL = `https://query1.finance.yahoo.com/v7/finance/quote?crumb=%s&symbols=%s`
const marketURLQueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

// Market stores current market information displayed in the top three lines of
// the screen. The market data is fetched and parsed from the HTML page above.
type Market struct {
	vendor    APISourceType
	IsClosed  bool              // True when U.S. markets are closed.
	Dow       map[string]string // Hash of Dow Jones indicators.
	Nasdaq    map[string]string // Hash of NASDAQ indicators.
	Sp500     map[string]string // Hash of S&P 500 indicators.
	Tokyo     map[string]string
	HongKong  map[string]string
	London    map[string]string
	Frankfurt map[string]string
	Yield     map[string]string
	Oil       map[string]string
	Yen       map[string]string
	Euro      map[string]string
	Gold      map[string]string
	Szzs      map[string]string
	Szcz      map[string]string
	Hs300     map[string]string
	Cybz      map[string]string
	regex     *regexp.Regexp // Regex to parse market data from HTML.
	errors    string         // Error(s), if any.
	url       string         // URL with symbols to fetch data
	cookies   string         // cookies for auth
	crumb     string         // crumb for the cookies, to be applied as a query param
}

// Returns new initialized Market struct.
func NewMarket(vendor APISourceType) *Market {
	market := &Market{}
	market.vendor = vendor
	market.IsClosed = false
	market.Dow = make(map[string]string)
	market.Nasdaq = make(map[string]string)
	market.Sp500 = make(map[string]string)

	market.Tokyo = make(map[string]string)
	market.HongKong = make(map[string]string)
	market.London = make(map[string]string)
	market.Frankfurt = make(map[string]string)

	market.Yield = make(map[string]string)
	market.Oil = make(map[string]string)
	market.Yen = make(map[string]string)
	market.Euro = make(map[string]string)
	market.Gold = make(map[string]string)

	market.Szzs = make(map[string]string)
	market.Szcz = make(map[string]string)

	market.Hs300 = make(map[string]string)
	market.Cybz = make(map[string]string)

	market.errors = ``

	const any = `\s*(?:.+?)`
	const change = `>([\+\-]?[\d\.,]+)<\/span>`
	const price = `>([\d\.,]+)<\/span>`
	const percent = `>([\+\-]?[\d\.,]+%?)<`

	rules := []string{
		`>Dow<`, any, percent, any, price, any, change, any,
		`>Nasdaq<`, any, percent, any, price, any, change, any,
		`">S&P<`, any, percent, any, price, any, change, any,
		`>10\-year yield<`, any, price, any, percent, any,
		`>Oil<`, any, price, any, percent, any,
		`>Yen<`, any, price, any, percent, any,
		`>Euro<`, any, price, any, percent, any,
		`>Gold<`, any, price, any, percent, any,
		`>Nikkei 225<`, any, percent, any, price, any, change, any,
		`>Hang Seng<`, any, percent, any, price, any, change, any,
		`>FTSE 100<`, any, percent, any, price, any, change, any,
		`>DAX<`, any, percent, any, price, any, change, any,
	}

	market.regex = regexp.MustCompile(strings.Join(rules, ``))

	market.cookies = fetchCookies()
	market.crumb = fetchCrumb(market.cookies)
	market.url = fmt.Sprintf(marketURL, market.crumb, `^DJI,^IXIC,^GSPC,^N225,^HSI,^FTSE,^GDAXI,^TNX,CL=F,JPY=X,EUR=X,GC=F`) + marketURLQueryParts

	return market
}

// Fetch downloads HTML page from the 'marketURL', parses it, and stores resulting data
// in internal hashes. If download or data parsing fails Fetch populates 'market.errors'.
func (market *Market) Fetch() (self *Market) {
	switch market.vendor {
	case API_VENDOR_YAHOO:
		return market.FetchYahoo()
	case API_VENDOR_QQ:
		mkt1 := market.FetchYahoo()
		mkt2 := market.FetchQQ()
		mkt1.Szzs = mkt2.Szzs
		mkt1.Szcz = mkt2.Szcz
		mkt1.Hs300 = mkt2.Hs300
		mkt1.Cybz = mkt2.Cybz
		return mkt1
	case API_VENDOR_SINA:
		mkt1 := market.FetchYahoo()
		mkt2 := market.FetchQQ()
		mkt1.Szzs = mkt2.Szzs
		mkt1.Szcz = mkt2.Szcz
		mkt1.Hs300 = mkt2.Hs300
		mkt1.Cybz = mkt2.Cybz
		return mkt1
	case API_VENDOR_NETEASE:
		mkt1 := market.FetchYahoo()
		mkt2 := market.FetchQQ()
		mkt1.Szzs = mkt2.Szzs
		mkt1.Szcz = mkt2.Szcz
		return mkt1
	case API_VENDOR_EASTMONEY:
		mkt1 := market.FetchYahoo()
		mkt2 := market.FetchQQ()
		mkt1.Szzs = mkt2.Szzs
		mkt1.Szcz = mkt2.Szcz
		mkt1.Hs300 = mkt2.Hs300
		mkt1.Cybz = mkt2.Cybz
		return mkt1
	case API_VENDOR_LIMITUP_EASTMONEY:
		mkt1 := market.FetchYahoo()
		mkt2 := market.FetchQQ()
		mkt1.Szzs = mkt2.Szzs
		mkt1.Szcz = mkt2.Szcz
		mkt1.Hs300 = mkt2.Hs300
		mkt1.Cybz = mkt2.Cybz
		return mkt1
	case API_VENDOR_LHB_EASTMONEY:
		mkt1 := market.FetchYahoo()
		mkt2 := market.FetchQQ()
		mkt1.Szzs = mkt2.Szzs
		mkt1.Szcz = mkt2.Szcz
		mkt1.Hs300 = mkt2.Hs300
		mkt1.Cybz = mkt2.Cybz
		return mkt1
	default:
		return market.FetchYahoo()
	}
}

// Fetch downloads HTML page from the 'marketURL', parses it, and stores resulting data
// in internal hashes. If download or data parsing fails Fetch populates 'market.errors'.
func (market *Market) FetchYahoo() (self *Market) {
	self = market // <-- This ensures we return correct market after recover() from panic().
	defer func() {
		if err := recover(); err != nil {
			market.errors = fmt.Sprintf("Error fetching market data...\n%s", err)
		} else {
			market.errors = ""
		}
	}()

	client := http.Client{}
	request, err := http.NewRequest("GET", market.url, nil)
	if err != nil {
		panic(err)
	}

	request.Header = http.Header{
		"Accept":          {"*/*"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Connection":      {"keep-alive"},
		"Content-Type":    {"application/json"},
		"Cookie":          {market.cookies},
		"Host":            {"query1.finance.yahoo.com"},
		"Origin":          {"https://finance.yahoo.com"},
		"Referer":         {"https://finance.yahoo.com"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
		"TE":              {"trailers"},
		"User-Agent":      {userAgent},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	body = market.isMarketOpen(body)
	return market.extract(body)
}

// Ok returns two values: 1) boolean indicating whether the error has occured,
// and 2) the error text itself.
func (market *Market) Ok() (bool, string) {
	return market.errors == ``, market.errors
}

// -----------------------------------------------------------------------------
func (market *Market) isMarketOpen(body []byte) []byte {
	// TBD -- CNN page doesn't seem to have market open/close indicator.
	return body
}

// -----------------------------------------------------------------------------
func assign(results []map[string]interface{}, position int, changeAsPercent bool) map[string]string {
	out := make(map[string]string)
	out[`change`] = float2Str(results[position]["regularMarketChange"].(float64))
	out[`latest`] = float2Str(results[position]["regularMarketPrice"].(float64))
	if changeAsPercent {
		out[`change`] = float2Str(results[position]["regularMarketChangePercent"].(float64)) + `%`
	} else {
		out[`percent`] = float2Str(results[position]["regularMarketChangePercent"].(float64))
	}
	return out
}

// -----------------------------------------------------------------------------
func (market *Market) extract(body []byte) *Market {
	d := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		panic(err)
	}
	results := d["quoteResponse"]["result"]
	market.Dow = assign(results, 0, false)
	market.Nasdaq = assign(results, 1, false)
	market.Sp500 = assign(results, 2, false)
	market.Tokyo = assign(results, 3, false)
	market.HongKong = assign(results, 4, false)
	market.London = assign(results, 5, false)
	market.Frankfurt = assign(results, 6, false)
	market.Yield[`name`] = `10-year Yield`
	market.Yield = assign(results, 7, false)

	market.Oil = assign(results, 8, true)
	market.Yen = assign(results, 9, true)
	market.Euro = assign(results, 10, true)
	market.Gold = assign(results, 11, true)

	return market
}
