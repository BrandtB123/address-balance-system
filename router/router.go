package router

import (
	"net/http"
	"unit410/middleware"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/generateReport", middleware.GenerateReport)
	return mux
}
