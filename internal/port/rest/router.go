package rest

import (
	"net/http"

	"github.com/serikdev/CashFlow/internal/port/rest/handler"
)

type Handlers struct {
	AccountHandler     *handler.AccountHandler
	TransactionHandler *handler.TransactionHandler
}

func NewRouter(handlers *Handlers) http.Handler {
	mux := http.NewServeMux()

	if handlers.AccountHandler != nil {
		handler.RegisterAccountRouter(mux, handlers.AccountHandler)
	}
	if handlers.TransactionHandler != nil {
		handler.RegisterTransactionRouter(mux, handlers.TransactionHandler)
	}

	return mux
}
