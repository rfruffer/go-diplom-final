package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/fylgushev/go-diplom-final/internal/domain/services"
	"github.com/fylgushev/go-diplom-final/internal/interfaces/grpc/handlers"
	"github.com/fylgushev/go-diplom-final/internal/interfaces/grpc/middleware"
	"github.com/fylgushev/go-diplom-final/pkg/proto/auth"
)

// Server представляет gRPC сервер приложения
type Server struct {
	grpcServer   *grpc.Server
	authHandler  *handlers.AuthHandler
	tokenManager *services.TokenManager
	port         int
}

// ServerConfig конфигурация gRPC сервера
type ServerConfig struct {
	Port           int
	TokenManager   *services.TokenManager
	AuthService    *services.AuthService
	MaxMessageSize int
}

// NewServer создает новый экземпляр gRPC сервера
func NewServer(config ServerConfig) *Server {
	// Создаем middleware для аутентификации
	authInterceptor := middleware.NewAuthInterceptor(config.TokenManager)

	// Настраиваем keepalive параметры
	kaep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: false,
	}

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
	}

	// Создаем gRPC сервер с middleware
	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.MaxRecvMsgSize(config.MaxMessageSize),
		grpc.MaxSendMsgSize(config.MaxMessageSize),
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor()),
	)

	// Создаем handlers
	authHandler := handlers.NewAuthHandler(config.AuthService)

	// Регистрируем сервисы
	auth.RegisterAuthServiceServer(grpcServer, authHandler)

	return &Server{
		grpcServer:   grpcServer,
		authHandler:  authHandler,
		tokenManager: config.TokenManager,
		port:         config.Port,
	}
}

// Start запускает gRPC сервер
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.port, err)
	}

	log.Printf("gRPC сервер запущен на порту %d", s.port)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

// Shutdown корректно останавливает gRPC сервер
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Начинаем остановку gRPC сервера...")

	// Канал для отслеживания завершения остановки
	done := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	// Ждем завершения или таймаута
	select {
	case <-done:
		log.Println("gRPC сервер корректно остановлен")
		return nil
	case <-ctx.Done():
		log.Println("Таймаут остановки gRPC сервера, принудительная остановка...")
		s.grpcServer.Stop()
		return ctx.Err()
	}
}