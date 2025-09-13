package handler

import (
	"encoding/json"
	"net/http"

	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/serikdev/CashFlow/internal/port/rest/handler/dto"
	"github.com/sirupsen/logrus"
)

type TransactionUsecase interface {
	Deposit(accountID int64, amount float64) (*entity.Transaction, error)
	Withdraw(accountID int64, amount float64) (*entity.Transaction, error)
	Transfer(fromAccountID, toAccountID int64, amount float64) (*entity.Transaction, error)
	ListTransactions(accountID int64) ([]entity.Transaction, error)
}

type TransactionHandler struct {
	*BaseHandler
	service TransactionUsecase
	logger  *logrus.Entry
}

func NewTransactionHandler(baseHandler *BaseHandler, service TransactionUsecase, logger *logrus.Entry) *TransactionHandler {
	return &TransactionHandler{
		BaseHandler: baseHandler,
		service:     service,
		logger:      logger,
	}
}

// Deposit godoc
// @Summary Пополнение счета
// @Description Пополняет баланс аккаунта
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "ID аккаунта"
// @Param request body dto.DepositRequest true "Сумма пополнения"
// @Success 201 {object} entity.Transaction
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /accounts/{id}/deposit [post]
func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := h.GetIDFromPath(r)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload dto.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tx, err := h.service.Deposit(id, payload.Amount)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondWithJSON(w, http.StatusCreated, tx)
}

// Withdraw godoc
// @Summary Снятие со счета
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "ID аккаунта"
// @Param request body dto.WithdrawRequest true "Сумма снятия"
// @Success 201 {object} entity.Transaction
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /accounts/{id}/withdraw [post]
func (h *TransactionHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := h.GetIDFromPath(r)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload dto.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tx, err := h.service.Withdraw(id, payload.Amount)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.RespondWithJSON(w, http.StatusCreated, tx)
}

// Transfer godoc
// @Summary Перевод средств
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "ID аккаунта-отправителя"
// @Param request body dto.TransferRequest true "Перевод"
// @Success 201 {object} entity.Transaction
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /accounts/{id}/transfer [post]
func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	fromID, err := h.GetIDFromPath(r)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	tx, err := h.service.Transfer(fromID, payload.ToAccountID, payload.Amount)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.RespondWithJSON(w, http.StatusCreated, tx)
}

// ListTransactions godoc
// @Summary История транзакций
// @Tags transactions
// @Produce json
// @Param id path int true "ID аккаунта"
// @Success 200 {array} entity.Transaction
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /accounts/{id}/transactions [get]
func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := h.GetIDFromPath(r)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	txs, err := h.service.ListTransactions(id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondWithJSON(w, http.StatusOK, txs)
}
