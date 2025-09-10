package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

// AccountUsecase defines the account service interface
type AccountUsecase interface {
	Create(ctx context.Context, account *entity.Account) (*entity.Account, error)
	GetByID(ctx context.Context, id int64) (*entity.Account, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]entity.Account, int, error)
}
type AccountHandler struct {
	*BaseHandler
	service AccountUsecase
	logger  *logrus.Entry
}

func NewAccountHandler(baseHandler *BaseHandler, service AccountUsecase, logger *logrus.Entry) *AccountHandler {
	return &AccountHandler{
		BaseHandler: baseHandler,
		service:     service,
		logger:      logger,
	}
}

// Create godoc
// @Summary Создать новый счет
// @Description Создает новый аккаунт с балансом и валютой
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body entity.Account true "Данные аккаунта"
// @Success 201 {object} entity.Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input entity.Account
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	account, err := h.service.Create(ctx, &input)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create account")
		h.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondWithJSON(w, http.StatusCreated, account)

}

// GetByID godoc
// @Summary Получить счет по ID
// @Tags accounts
// @Produce json
// @Param id path int true "ID аккаунта"
// @Success 200 {object} entity.Account
// @Failure 404 {object} ErrorResponse
// @Router /accounts/{id} [get]
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

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	account, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.RespondWithJSON(w, http.StatusOK, account)

}

// Delete godoc
// @Summary Удалить счет
// @Tags accounts
// @Param id path int true "ID аккаунта"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /accounts/{id} [delete]
func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := h.GetIDFromPath(r)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.WithError(err).WithField("account_id", id).Error("Failed to delete account")
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	h.RespondWithJSON(w, http.StatusNoContent, nil)
}

// List godoc
// @Summary Список счетов
// @Tags accounts
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество элементов"
// @Success 200 {object} map[string]interface{}
// @Router /accounts [get]
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	account, total, err := h.service.List(ctx, page, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to fetch account")
		h.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": account,
		"meta": map[string]interface{}{
			"total":        total,
			"current_page": page,
			"last_page":    (total + limit - 1) / limit,
		},
	})
}
