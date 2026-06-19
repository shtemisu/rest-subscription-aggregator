package controller

import "net/http"

func NewRouter(h *SubsHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("GET /subscriptions", h.List)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetById)
	mux.HandleFunc("PUT /subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Delete)

	mux.HandleFunc("GET /subcriptions/summary", h.SumUsingFilter)
	return mux
}
