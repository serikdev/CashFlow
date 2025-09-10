package dto // data transfer object

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type TransferRequest struct {
	ToAccountID int64   `json:"to_account_id"`
	Amount      float64 `json:"amount"`
}
