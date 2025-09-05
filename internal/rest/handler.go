package rest

import (
	"encoding/json"
	"net/http"

	logger_lib "github.com/s21platform/logger-lib"
	api "github.com/s21platform/optionhub-service/internal/generated"
)

type Handler struct {
	dbR DbRepo
}

func New(dbR DbRepo) *Handler {
	return &Handler{dbR: dbR}
}

func (h *Handler) GetOptions(w http.ResponseWriter, r *http.Request, params api.GetOptionsParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	
	options, err := h.dbR.GetOptionsByAttributeId(ctx, params.AttributeId)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to get options")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	resultOptions := make([]api.OptionItem, len(options))
	for i, option := range options {
		resultOptions[i] = api.OptionItem{
			OptionId: option.OptionID,
			AttributeId: option.AttributeID,
			Label: option.Label,
		}
	}

	result := &api.Options{
		Data: resultOptions,
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to marshal options")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}