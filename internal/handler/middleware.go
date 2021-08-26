package handler

import (
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

func JWT(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/auth":
			case "/refresh":
			default:
				authHeader := r.Header.Get("Authorization")
				if len(authHeader) == 0 {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
				h := strings.SplitN(authHeader, " ", 2)
				if len(h) != 2 {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
				if strings.ToLower(h[0]) != "bearer" {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
				_, err := jwt.ParseString(h[1], jwt.WithVerify(jwa.HS256, secret), jwt.WithValidate(true))
				if err != nil {
					writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
