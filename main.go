package main

import (
	"fmt"
	"time"
	"unit410/api"
	"unit410/db"
)

func getBalances(chain string) error {
	var a api.API
	switch chain {
	case "near":
		a = &api.NearAPI{}
	case "audius":
		a = &api.AudiusAPI{}
	case "osmosis":
		a = &api.OsmosisAPI{}
	}
	return a.GetData()
}

var chains = []string{
	"near",
	"audius",
	"osmosis",
}

func main() {
	db.NewDB()

	db.CreateAddressesTable()
	db.CreateBalancesTable()
	// for _, chain := range chains {
	// 	getBalances(chain)
	// }
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, now.Location())
	fmt.Println(db.GetBalancesByDate(startOfToday))

	// // Create a ticker that ticks every 24 hours (1 day).
	// ticker := time.NewTicker(24 * time.Hour)

	// // Start a goroutine to execute your function at midnight.
	// go func() {
	// 	durationUntilMidnight := time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Now().Location()))
	// 	<-time.After(durationUntilMidnight)

	// 	for _, chain := range chains {
	// 		getBalances(chain)
	// 	}
	// 	for range ticker.C {
	// 		for _, chain := range chains {
	// 			getBalances(chain)
	// 		}
	// 	}
	// }()

	// select {}
}
