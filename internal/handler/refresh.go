package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/models"
	"github.com/andreipimenov/golang-training-2021/internal/restapi/operations/auth"
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

func (h *Refresh) Handle(req auth.RefreshTokenParams) middleware.Responder {
	newAccessToken, newRefreshToken, err := h.service.Refresh(req.Body.AccessToken, req.Body.RefreshToken)
	if err != nil {
		h.logger.Error().Err(err).Msg("Token refreshing error")
		return auth.NewRefreshTokenForbidden().WithPayload(&models.Error{"Forbidden"})
	}

	return auth.NewRefreshTokenOK().WithPayload(&models.Tokens{AccessToken: newAccessToken, RefreshToken: newRefreshToken})
}
