package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/andreipimenov/golang-training-2021/internal/config"
	"github.com/andreipimenov/golang-training-2021/internal/handler"
	"github.com/andreipimenov/golang-training-2021/internal/pb"
	"github.com/andreipimenov/golang-training-2021/internal/repository"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Configuration error")
	}

	db, err := sql.Open("postgres", cfg.DBConnString)
	if err != nil {
		logger.Fatal().Err(err).Msg("DB initializing error")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal().Err(err).Msg("DB pinging error")
	}

	stockRepo := repository.NewCache()
	stockService := service.NewStock(&logger, stockRepo, cfg.ExternalAPIToken)
	stockHandler := handler.NewStock(&logger, stockService)

	authRepo := repository.NewAuth(db)
	authService := service.NewAuth(&logger, authRepo, []byte(cfg.Secret))
	authHandler := handler.NewAuth(&logger, authService)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(),
		grpc_auth.UnaryServerInterceptor(handler.NewAuthMiddleware([]byte(cfg.Secret)).AuthFunc),
	)))
	pb.RegisterAuthServer(grpcServer, authHandler)
	pb.RegisterStockServer(grpcServer, stockHandler)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal().Err(err).Msg("Listening gRPC error")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		logger.Info().Msgf("GRPC server is listening on :%d", cfg.GRPCPort)
		err := grpcServer.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
			logger.Fatal().Err(err).Msg("GRPC server error")
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterAuthHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", cfg.GRPCPort), opts)
	if err != nil {
		logger.Fatal().Err(err).Msg("Registering gRPC gateway endpoint error")
	}
	err = pb.RegisterStockHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", cfg.GRPCPort), opts)
	if err != nil {
		logger.Fatal().Err(err).Msg("Registering gRPC gateway endpoint error")
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		logger.Info().Msgf("GRPC gateway server is listening on :%d", cfg.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("GRPC gateway server error")
		}
	}()

	<-shutdown

	logger.Info().Msg("Shutdown signal received")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("GRPC gateway server shutdown error")
	}

	grpcServer.GracefulStop()

	logger.Info().Msg("Server stopped gracefully")
}
