// Package services предоставляет бизнес-логику для приложения GophKeeper.
//
// Этот пакет содержит сервисные компоненты, реализующие основную бизнес-логику
// приложения, включая аутентификацию, управление данными, синхронизацию и валидацию.
//
// Основные компоненты:
//   - AuthService: аутентификация и авторизация пользователей
//   - DataService: управление зашифрованными данными пользователей
//   - SyncService: синхронизация данных между клиентами
//   - DataSerializer: сериализация и десериализация различных типов данных
//   - Validator: валидация пользовательских данных
//   - TokenManager: управление JWT токенами
//
// Пример использования AuthService:
//
//	authService := services.NewAuthService(userRepo, tokenManager)
//	response, err := authService.Register(ctx, &auth.RegisterRequest{
//	    Login: "user@example.com",
//	    Password: "securepassword",
//	})
//
// Пример использования DataService:
//
//	masterKey := []byte("encryption-key")
//	dataService, err := services.NewDataService(dataRepo, masterKey)
//	loginData := entities.NewLoginPasswordData("user", "password")
//	item, err := dataService.StoreLoginPassword(ctx, userID, "My Login", loginData, nil)
//
// Пример использования Validator:
//
//	validator := services.NewValidator()
//	err := validator.ValidateLoginPassword(login, password, url)
package services
