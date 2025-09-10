package rest

import (
	"net/http"

	"github.com/serikdev/CashFlow/internal/port/rest/handler"
)

type Handlers struct {
	AccountHandler *handler.AccountHandler
}

func NewRouter(handlers *Handlers) http.Handler {
	mux := http.NewServeMux()

	if handlers.AccountHandler != nil {
		handler.RegisterAccountRouter(mux, handlers.AccountHandler)
	}

	return mux
}
