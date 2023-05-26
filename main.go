package main

import (
	"fmt"
	"time"
	"unit410/db"
	"unit410/middleware"
)

var chains = []string{
	"near",
	"audius",
	"osmosis",
}

func main() {
	db.NewDB()
	err := db.CreateAddressesTable()
	if err != nil {
		fmt.Println("error creating addresses table: ", err.Error())
	}
	err = db.CreateBalancesTable()
	if err != nil {
		fmt.Println("error creating balances table: ", err.Error())
	}

	err = middleware.GetBalances(chains)
	if err != nil {
		fmt.Println("error getting balances: ", err.Error())
	}
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, now.Location())
	fmt.Println(db.GetBalancesByDate(startOfToday))

	// Create a ticker that ticks every 24 hours
	ticker := time.NewTicker(24 * time.Hour)

	// Start a goroutine to execute at midnight
	go func() {
		durationUntilMidnight := time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Now().Location()))
		<-time.After(durationUntilMidnight)

		middleware.GetBalances(chains)
		for range ticker.C {
			middleware.GetBalances(chains)
		}
	}()

	select {}
}
