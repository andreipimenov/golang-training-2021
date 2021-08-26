// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"database/sql"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/config"
	"github.com/andreipimenov/golang-training-2021/internal/handler"
	"github.com/andreipimenov/golang-training-2021/internal/repository"
	"github.com/andreipimenov/golang-training-2021/internal/restapi/operations"
	"github.com/andreipimenov/golang-training-2021/internal/restapi/operations/auth"
	"github.com/andreipimenov/golang-training-2021/internal/restapi/operations/stock"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

//go:generate swagger generate server --target ../../internal --name StockService --spec ../../api/api.yaml --principal interface{}

func configureFlags(api *operations.StockServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.StockServiceAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Configuration error")
	}

	db, err := sql.Open("postgres", cfg.DBConnString)
	if err != nil {
		logger.Fatal().Err(err).Msg("DB initializing error")
	}

	stockRepo := repository.NewCache()
	stockService := service.NewStock(&logger, stockRepo, cfg.ExternalAPIToken)
	stockHandler := handler.NewStock(&logger, stockService)

	authRepo := repository.NewAuth(db)
	authService := service.NewAuth(&logger, authRepo, []byte(cfg.Secret))
	authHandler := handler.NewAuth(&logger, authService)
	refreshHandler := handler.NewRefresh(&logger, authService)

	api.StockGetPriceHandler = stockHandler
	api.AuthAuthHandler = authHandler
	api.AuthRefreshTokenHandler = refreshHandler

	if api.AuthAuthHandler == nil {
		api.AuthAuthHandler = auth.AuthHandlerFunc(func(params auth.AuthParams) middleware.Responder {
			return middleware.NotImplemented("operation auth.Auth has not yet been implemented")
		})
	}
	if api.StockGetPriceHandler == nil {
		api.StockGetPriceHandler = stock.GetPriceHandlerFunc(func(params stock.GetPriceParams) middleware.Responder {
			return middleware.NotImplemented("operation stock.GetPrice has not yet been implemented")
		})
	}
	if api.AuthRefreshTokenHandler == nil {
		api.AuthRefreshTokenHandler = auth.RefreshTokenHandlerFunc(func(params auth.RefreshTokenParams) middleware.Responder {
			return middleware.NotImplemented("operation auth.RefreshToken has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares), cfg)
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(h http.Handler, cfg *config.Config) http.Handler {
	f := handler.JWT([]byte(cfg.Secret))
	return f(h)
}
