package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

const (
	RefreshPath = "/refresh"
)

type Refresh struct {
	logger  *zerolog.Logger
	service AuthService
}

func NewRefresh(logger *zerolog.Logger, srv AuthService) *Refresh {
	return &Refresh{
		logger:  logger,
		service: srv,
	}
}

func (h *Refresh) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &model.Tokens{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
		return
	}

	newAccessToken, newRefreshToken, err := h.service.Refresh(req.AccessToken, req.RefreshToken)
	if err != nil {
		h.logger.Error().Err(err).Msg("Token refreshing error")
		writeResponse(w, http.StatusForbidden, model.Error{Error: "Forbidden"})
		return
	}

	res := &model.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	writeResponse(w, http.StatusOK, res)
}
