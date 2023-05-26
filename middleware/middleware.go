package middleware

import (
	"encoding/json"
	"net/http"
	"unit410/api"
)

// API method to handle calculateAllocation request
func GenerateReport(w http.ResponseWriter,
	r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	// do stuff

	json.NewEncoder(w).Encode(nil)
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
		return balanceAPI.GetData()
	}
	return nil
}

var chains = []string{
	"near",
	"audius",
	"osmosis",
}
