// Package config предоставляет конфигурацию для сервера и клиента GophKeeper.
//
// Этот пакет содержит структуры конфигурации и функции для работы с настройками
// приложения, включая конфигурацию базы данных и клиентские настройки.
//
// Основные компоненты:
//   - DatabaseConfig: конфигурация подключения к PostgreSQL
//   - ClientConfig: конфигурация клиентского приложения (токены, настройки)
//
// Пример использования DatabaseConfig:
//
//	cfg := config.DefaultDatabaseConfig()
//	cfg.Host = "localhost"
//	cfg.Port = 5432
//	db, err := cfg.Connect()
//
// Пример использования ClientConfig:
//
//	cfg, err := config.LoadClientConfig()
//	if err != nil {
//	    cfg = config.DefaultClientConfig()
//	}
//	cfg.SetAuthTokens(accessToken, refreshToken, expiresAt, userID, login)
//	err = config.SaveClientConfig(cfg)
package config
