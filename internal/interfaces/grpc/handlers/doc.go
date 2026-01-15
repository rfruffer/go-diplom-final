// Package handlers предоставляет gRPC обработчики для API GophKeeper.
//
// Этот пакет содержит реализации gRPC сервисов, которые обрабатывают
// входящие запросы от клиентов и взаимодействуют с бизнес-логикой.
//
// Основные компоненты:
//   - AuthHandler: обработка запросов аутентификации (регистрация, вход, обновление токенов)
//   - SyncHandler: обработка запросов синхронизации данных между клиентами
//
// Пример использования:
//
//	authHandler := handlers.NewAuthHandler(authService)
//	authServer := grpc.NewServer()
//	auth.RegisterAuthServiceServer(authServer, authHandler)
package handlers
