package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

type AccountUsecase interface {
	Create(ctx context.Context, account *entity.Account) (*entity.Account, error)
	GetByID(ctx context.Context, id int64) (*entity.Account, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]entity.Account, int, error)
}
type AccountHandler struct {
	BaseHandler
	service AccountUsecase
	logger  *logrus.Entry
}

func NewAccountHandler(baseHandler BaseHandler, service AccountUsecase, logger *logrus.Entry) *AccountHandler {
	return &AccountHandler{
		BaseHandler: baseHandler,
		service:     service,
		logger:      logger,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not alowed")
		return
	}

	var input entity.Account
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	account, err := h.usecase.Create(ctx, input)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create account")
		h.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondWithJSON(w, http.StatusCreated, account)

}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	id, err := h.GetIDFromPath(r)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.usecase.GetAccount(id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.RespondWithJSON(w, http.StatusOK, account)

}
