package handler

import "net/http"

func RegisterTransactionRouter(mux *http.ServeMux, transactionHandler *TransactionHandler) {
	mux.HandleFunc("POST /api/accounts/{id}/deposit", transactionHandler.Deposit)
	mux.HandleFunc("POST /api/accounts/{id}/withdraw", transactionHandler.Withdraw)
	mux.HandleFunc("POST /api/accounts/{id}/transfer", transactionHandler.Transfer)
	mux.HandleFunc("POST /api/accounts/{id}/transactions", transactionHandler.ListTransactions)
}
