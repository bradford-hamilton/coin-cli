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
	res, err := http.Get("https://api.coincap.io/v2/assets?limit=10")
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
		case "1":
			renderCoin(a.Data[0])
		case "2":
			renderCoin(a.Data[1])
		case "3":
			renderCoin(a.Data[2])
		case "4":
			renderCoin(a.Data[3])
		case "5":
			renderCoin(a.Data[4])
		case "6":
			renderCoin(a.Data[5])
		case "7":
			renderCoin(a.Data[6])
		case "8":
			renderCoin(a.Data[7])
		case "9":
			renderCoin(a.Data[8])
		case "10":
			renderCoin(a.Data[9])
		}
	}
}

func renderCoin(coinData CoinData) {
	ui.Clear()

	s, err := strconv.ParseFloat(coinData.PriceUsd, 32)
	if err == nil {
		Log.Println("parse float err: ", err)
	}

	p := ui.NewParagraph(fmt.Sprintf("Price (USD) %s", fmt.Sprintf("%.2f", s)))
	p.Height = 3
	p.Width = 24
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = coinData.Symbol
	p.BorderFg = ui.ColorCyan

	ui.Render(p)
}
