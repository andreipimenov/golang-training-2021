package handler

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	authMethod    = "/pb.Auth/Authenticate"
	refreshMethod = "/pb.Auth/Refresh"
)

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret []byte) *AuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
	}
}

func (a *AuthMiddleware) AuthFunc(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	switch method {
	case authMethod:
	case refreshMethod:
	default:
		tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		_, err = jwt.ParseString(tokenString, jwt.WithVerify(jwa.HS256, a.secret), jwt.WithValidate(true))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
		}
	}
	return ctx, nil
}
