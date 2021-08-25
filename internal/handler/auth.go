package handler

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andreipimenov/golang-training-2021/internal/pb"
)

const (
	AuthPath = "/auth"
)

type Auth struct {
	pb.UnimplementedAuthServer
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

func (a *Auth) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	accessToken, refreshToken, err := a.service.Authenticate(in.Username, in.Password)
	if err != nil {
		a.logger.Error().Err(err).Msg("Authentication error")
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	return &pb.AuthenticateResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	newAccessToken, newRefreshToken, err := a.service.Refresh(in.AccessToken, in.RefreshToken)
	if err != nil {
		a.logger.Error().Err(err).Msg("Token refreshing error")
		return nil, status.Error(codes.Unauthenticated, "Invalid access and refresh tokens pair")
	}

	return &pb.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
