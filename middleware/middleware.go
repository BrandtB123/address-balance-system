package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unit410/api"
	"unit410/db"
)

var chains = []string{
	"near",
	"audius",
	"osmosis",
}

func GenerateReport(w http.ResponseWriter,
	r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")

	err := GetBalances(chains)
	if err != nil {
		w.WriteHeader(400)
	}

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, now.Location())
	balances, err := db.GetBalancesByDate(startOfToday)
	if err != nil {
		w.WriteHeader(400)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(balances)
}

// API method to handle calculateAllocation request
func GetReport(w http.ResponseWriter,
	r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	vars := r.URL.Query()
	chain, ok := vars["chain"]
	if !ok {
		fmt.Println("Must need chain")
		w.WriteHeader(400)
	}

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, now.Location())
	balances, err := db.GetBalancesByDateAndAsset(chain[0], startOfToday)
	if err != nil {
		fmt.Println("error getting balances: ", err.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(balances)
}

// API method to handle calculateAllocation request
func GetAddressBalance(w http.ResponseWriter,
	r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	vars := r.URL.Query()
	address, ok := vars["address"]
	if !ok {
		fmt.Println("Must need chain")
		w.WriteHeader(400)
	}

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, now.Location())
	balances, err := db.GetBalancesByDateAndAddress(address[0], startOfToday)
	if err != nil {
		fmt.Println("error getting balances: ", err.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(balances)
}

func GetBalances(chains []string) error {
	for _, chain := range chains {
		var balanceAPI api.API
		switch chain {
		case "near":
			balanceAPI = &api.NearAPI{}
		case "audius":
			balanceAPI = &api.AudiusAPI{}
		case "osmosis":
			balanceAPI = &api.OsmosisAPI{}
		}
		err := balanceAPI.GetData()
		if err != nil {
			return err
		}
	}
	return nil
}
