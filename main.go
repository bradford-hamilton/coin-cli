package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	ui "github.com/gizak/termui"
)

// Log sets up our logger for debugging
var Log *log.Logger

// Assets is the response shape for assets request
type Assets struct {
	Data []CoinData
}

// CoinData is a single asset shape from request
type CoinData struct {
	ID                string
	Rank              string
	Symbol            string
	Name              string
	Supply            string
	MaxSupply         string
	MarketCapUsd      string
	VolumeUsd24Hr     string
	PriceUsd          string
	ChangePercent24Hr string
	Vwap24Hr          string
}

func main() {
	// debug logging
	f, err := os.OpenFile("debugLog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	Log = log.New(f, "", 0)
	Log.SetOutput(f)

	// start ui thread
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// TODO: add a ui loading
	res, err := http.Get("https://api.coincap.io/v2/assets?limit=9")
	if err != nil {
		Log.Println("HTTP request err: ", err)
	}
	defer res.Body.Close()

	a := &Assets{}
	if err = json.NewDecoder(res.Body).Decode(a); err != nil {
		Log.Println("json decoder err: ", err)
	}

	coins := []string{}
	for i, v := range a.Data {
		coins = append(coins, fmt.Sprintf("[%d] [%s](fg-blue)", i+1, v.Symbol))
	}

	ls := ui.NewList()
	ls.Items = coins
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Choose coin"
	ls.Height = 12
	ls.Width = 24
	ls.Y = 0

	ui.Render(ls)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			renderCoin(a.Data, e.ID)
		}
	}
}

func renderCoin(coinData []CoinData, ind string) {
	ui.Clear()

	i, _ := strconv.Atoi(ind)
	i--

	price, err := strconv.ParseFloat(coinData[i].PriceUsd, 32)
	if err == nil {
		Log.Println("parse float err: ", err)
	}

	p := ui.NewParagraph(fmt.Sprintf("Price (USD) %s", fmt.Sprintf("%.2f", price)))
	p.Height = 3
	p.Width = 24
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = coinData[i].Symbol
	p.BorderFg = ui.ColorCyan

	ui.Render(p)
}
