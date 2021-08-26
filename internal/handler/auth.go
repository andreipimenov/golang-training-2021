package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/models"
	"github.com/andreipimenov/golang-training-2021/internal/restapi/operations/auth"
)

type Auth struct {
	logger  *zerolog.Logger
	service AuthService
}

type AuthService interface {
	Authenticate(string, string) (string, string, error)
	Refresh(string, string) (string, string, error)
}

func NewAuth(logger *zerolog.Logger, srv AuthService) *Auth {
	return &Auth{
		logger:  logger,
		service: srv,
	}
}

func (h *Auth) Handle(req auth.AuthParams) middleware.Responder {
	accessToken, refreshToken, err := h.service.Authenticate(req.Body.Username, req.Body.Password)
	if err != nil {
		h.logger.Error().Err(err).Msg("Authentication error")
		return auth.NewAuthForbidden().WithPayload(&models.Error{"Forbidden"})
	}

	return auth.NewAuthOK().WithPayload(&models.Tokens{AccessToken: accessToken, RefreshToken: refreshToken})
}
