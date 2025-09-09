package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

type BaseHandler struct {
	logger *logrus.Entry
}

func NewBaseHandler(logger *logrus.Entry) BaseHandler {
	return BaseHandler{logger: logger}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (b *BaseHandler) RespondWithError(w http.ResponseWriter, code int, message string) {
	b.logger.WithFields(logrus.Fields{
		"code":    code,
		"message": message,
	}).Error("API error response")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := ErrorResponse{
		Error:   http.StatusText(code),
		Code:    code,
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}

func (b *BaseHandler) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			b.logger.WithError(err).Error("Failed encoding JSON")
		}
	}
}

func (b *BaseHandler) GetIDFromPath(r *http.Request) (int64, error) {
	idStr := path.Base(r.URL.Path)
	if idStr == "/" || idStr == "." || idStr == "" {
		return 0, fmt.Errorf("no ID found path")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	return id, nil
}
