package middleware

import (
	"encoding/json"
	"net/http"
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
