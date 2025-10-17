// Package main содержит точку входа для серверного приложения GophKeeper.
// Сервер реализует gRPC API для управления пользователями и их данными.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fylgushev/go-diplom-final/internal/domain/services"
	"github.com/fylgushev/go-diplom-final/internal/infrastructure/config"
	"github.com/fylgushev/go-diplom-final/internal/infrastructure/repository"
	grpcServer "github.com/fylgushev/go-diplom-final/internal/interfaces/grpc"
	"github.com/fylgushev/go-diplom-final/pkg/version"
)

// main основная функция серверного приложения.
func main() {
	fmt.Println("GophKeeper Server")
	fmt.Println("==================")
	fmt.Println(version.FormatForCLI())
	fmt.Println()

	// Инициализация базы данных
	dbConfig := config.DefaultDatabaseConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Выполнение миграций
	if err := config.Migrate(db); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	// Создание репозиториев
	userRepo := repository.NewUserRepository(db)

	// Создание JWT менеджера
	jwtConfig := services.DefaultConfig()
	// todo: получать ключ из переменных окружения
	jwtConfig.SigningKey = []byte("super-secret-jwt-key-change-in-production")
	tokenManager := services.NewTokenManagerFromConfig(jwtConfig)

	// Создание сервисов
	authService := services.NewAuthService(userRepo, tokenManager)

	// Конфигурация gRPC сервера
	serverConfig := grpcServer.ServerConfig{
		Port:           8080,
		TokenManager:   tokenManager,
		AuthService:    authService,
		MaxMessageSize: 4 * 1024 * 1024, // 4MB
	}

	// Создание и запуск gRPC сервера
	server := grpcServer.NewServer(serverConfig)

	// Настройка graceful shutdown
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание сигнала для остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Получен сигнал остановки сервера...")

	// Остановка сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка остановки сервера: %v", err)
	}

	log.Println("Сервер остановлен")
}
