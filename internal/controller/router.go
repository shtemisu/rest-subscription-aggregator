package controller

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "subAggregator/docs"
)

func NewRouter(h *SubsHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("GET /subscriptions", h.List)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetById)
	mux.HandleFunc("PUT /subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Delete)

	mux.HandleFunc("GET /subscriptions/summary", h.SumUsingFilter)
	return mux
}
