package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

const (
	AuthPath = "/auth"
)

type Auth struct {
	logger  *zerolog.Logger
	service AuthService
}

type AuthService interface {
	Authenticate(string, string) (string, string, error)
}

func NewAuth(logger *zerolog.Logger, srv AuthService) *Auth {
	return &Auth{
		logger:  logger,
		service: srv,
	}
}

func (h *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &model.AuthRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
		return
	}

	accessToken, refreshToken, err := h.service.Authenticate(req.Username, req.Password)
	if err != nil {
		h.logger.Error().Err(err).Msg("Authentication error")
		writeResponse(w, http.StatusForbidden, model.Error{Error: "Forbidden"})
		return
	}

	res := &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	writeResponse(w, http.StatusOK, res)
}
