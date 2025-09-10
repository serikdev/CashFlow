package handler

import "net/http"

func RegisterAccountRouter(mux *http.ServeMux, accountHandler *AccountHandler) {
	mux.HandleFunc("POST /api/accounts", accountHandler.Create)
	mux.HandleFunc("GET /api/accounts/", accountHandler.GetByID)
	mux.HandleFunc("DELETE /api/accounts/", accountHandler.Delete)
	mux.HandleFunc("GET /api/accounts", accountHandler.List)
}
