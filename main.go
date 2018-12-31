package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	ui "github.com/gizak/termui"
)

// Log sets up our logger for debugging
var Log *log.Logger

// Assets is the response shape for coincap assets request
type Assets struct {
	data []struct {
		ID                string `json:"id"`
		Rank              string `json:"rank"`
		Symbol            string `json:"symbol"`
		Name              string `json:"name"`
		Supply            string `json:"supply"`
		MaxSupply         string `json:"maxSupply"`
		MarketCapUsd      string `json:"marketCapUsd"`
		VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
		PriceUsd          string `json:"priceUsd"`
		ChangePercent24Hr string `json:"changePercent24Hr"`
		Vwap24Hr          string `json:"vwap24Hr"`
	}
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

	res, err := http.Get("https://api.coincap.io/v2/assets?limit=10")
	if err != nil {
		Log.Println("HTTP request err: ", err)
	}
	defer res.Body.Close()

	a := &Assets{}
	r := json.NewDecoder(res.Body).Decode(a)
	Log.Println("HTTP response: ", r)

	strs := []string{
		"[1] [BTC](fg-blue)",
		"[2] [LTC](fg-blue)",
		"[3] [ETH](fg-blue)",
	}

	ls := ui.NewList()
	ls.Items = strs
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Choose coin for price"
	ls.Height = 5
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
			getPrice("BTC")
		case "2":
			getPrice("LTC")
		case "3":
			getPrice("ETH")
		}
	}
}

func getPrice(ticker string) {
	ui.Clear()

	p := ui.NewParagraph("         3395")
	p.Height = 3
	p.Width = 24
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = ticker
	p.BorderFg = ui.ColorCyan

	ui.Render(p)
}
