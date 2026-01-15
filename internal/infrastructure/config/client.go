package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ClientConfig конфигурация клиента
type ClientConfig struct {
	ServerAddress   string `json:"server_address"`
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	TokenExpiresAt  int64  `json:"token_expires_at"`
	UserID          string `json:"user_id"`
	Login           string `json:"login"`
	LastSyncTime    int64  `json:"last_sync_time"`
	ConfigVersion   string `json:"config_version"`
}

// DefaultClientConfig возвращает конфигурацию клиента по умолчанию
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ServerAddress:  "localhost:8080",
		ConfigVersion:  "1.0",
		LastSyncTime:   0,
	}
}

// GetConfigPath возвращает путь к файлу конфигурации
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("не удалось определить домашний каталог: %w", err)
	}

	configDir := filepath.Join(homeDir, ".gophkeeper")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("не удалось создать каталог конфигурации: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadClientConfig загружает конфигурацию из файла
func LoadClientConfig() (*ClientConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Если файл не существует, создаем конфигурацию по умолчанию
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultClientConfig()
		if err := SaveClientConfig(config); err != nil {
			return nil, fmt.Errorf("не удалось создать конфигурацию по умолчанию: %w", err)
		}
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}

	var config ClientConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("не удалось разобрать конфигурацию: %w", err)
	}

	return &config, nil
}

// SaveClientConfig сохраняет конфигурацию в файл
func SaveClientConfig(config *ClientConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать конфигурацию: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("не удалось записать файл конфигурации: %w", err)
	}

	return nil
}

// IsAuthenticated проверяет, аутентифицирован ли пользователь
func (c *ClientConfig) IsAuthenticated() bool {
	return c.AccessToken != "" && c.TokenExpiresAt > time.Now().Unix()
}

// IsTokenExpired проверяет, истек ли токен
func (c *ClientConfig) IsTokenExpired() bool {
	return c.TokenExpiresAt <= time.Now().Unix()
}

// SetAuthTokens устанавливает токены аутентификации
func (c *ClientConfig) SetAuthTokens(accessToken, refreshToken string, expiresAt time.Time, userID, login string) {
	c.AccessToken = accessToken
	c.RefreshToken = refreshToken
	c.TokenExpiresAt = expiresAt.Unix()
	c.UserID = userID
	c.Login = login
}

// ClearAuthTokens очищает токены аутентификации
func (c *ClientConfig) ClearAuthTokens() {
	c.AccessToken = ""
	c.RefreshToken = ""
	c.TokenExpiresAt = 0
	c.UserID = ""
	c.Login = ""
}

// UpdateLastSyncTime обновляет время последней синхронизации
func (c *ClientConfig) UpdateLastSyncTime() {
	c.LastSyncTime = time.Now().Unix()
}

// GetLastSyncTime возвращает время последней синхронизации
func (c *ClientConfig) GetLastSyncTime() time.Time {
	return time.Unix(c.LastSyncTime, 0)
}