package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========== DatabaseConfig Tests ==========

func TestDefaultDatabaseConfig(t *testing.T) {
	cfg := DefaultDatabaseConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "postgres", cfg.Username)
	assert.Equal(t, "postgres", cfg.Password)
	assert.Equal(t, "gophkeeper", cfg.DatabaseName)
	assert.Equal(t, "disable", cfg.SSLMode)
	assert.Equal(t, 25, cfg.MaxOpenConns)
	assert.Equal(t, 5, cfg.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, cfg.MaxLifetime)
}

func TestDatabaseConfig_DSN(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:         "localhost",
		Port:         5432,
		Username:     "testuser",
		Password:     "testpass",
		DatabaseName: "testdb",
		SSLMode:      "disable",
	}

	dsn := cfg.DSN()

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expected, dsn)
}

func TestDatabaseConfig_DSN_WithDifferentValues(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:         "db.example.com",
		Port:         5433,
		Username:     "admin",
		Password:     "secret123",
		DatabaseName: "production",
		SSLMode:      "require",
	}

	dsn := cfg.DSN()

	expected := "host=db.example.com port=5433 user=admin password=secret123 dbname=production sslmode=require"
	assert.Equal(t, expected, dsn)
}

// ========== ClientConfig Tests ==========

func TestDefaultClientConfig(t *testing.T) {
	cfg := DefaultClientConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost:8080", cfg.ServerAddress)
	assert.Equal(t, "1.0", cfg.ConfigVersion)
	assert.Equal(t, int64(0), cfg.LastSyncTime)
	assert.Empty(t, cfg.AccessToken)
	assert.Empty(t, cfg.RefreshToken)
}

func TestClientConfig_IsAuthenticated_True(t *testing.T) {
	cfg := DefaultClientConfig()
	cfg.AccessToken = "token"
	cfg.TokenExpiresAt = time.Now().Add(time.Hour).Unix()

	assert.True(t, cfg.IsAuthenticated())
}

func TestClientConfig_IsAuthenticated_False_NoToken(t *testing.T) {
	cfg := DefaultClientConfig()
	cfg.AccessToken = ""
	cfg.TokenExpiresAt = time.Now().Add(time.Hour).Unix()

	assert.False(t, cfg.IsAuthenticated())
}

func TestClientConfig_IsAuthenticated_False_Expired(t *testing.T) {
	cfg := DefaultClientConfig()
	cfg.AccessToken = "token"
	cfg.TokenExpiresAt = time.Now().Add(-time.Hour).Unix()

	assert.False(t, cfg.IsAuthenticated())
}

func TestClientConfig_IsTokenExpired_True(t *testing.T) {
	cfg := DefaultClientConfig()
	cfg.TokenExpiresAt = time.Now().Add(-time.Hour).Unix()

	assert.True(t, cfg.IsTokenExpired())
}

func TestClientConfig_IsTokenExpired_False(t *testing.T) {
	cfg := DefaultClientConfig()
	cfg.TokenExpiresAt = time.Now().Add(time.Hour).Unix()

	assert.False(t, cfg.IsTokenExpired())
}

func TestClientConfig_SetAuthTokens(t *testing.T) {
	cfg := DefaultClientConfig()
	accessToken := "access_token"
	refreshToken := "refresh_token"
	expiresAt := time.Now().Add(time.Hour)
	userID := "user-123"
	login := "testuser"

	cfg.SetAuthTokens(accessToken, refreshToken, expiresAt, userID, login)

	assert.Equal(t, accessToken, cfg.AccessToken)
	assert.Equal(t, refreshToken, cfg.RefreshToken)
	assert.Equal(t, expiresAt.Unix(), cfg.TokenExpiresAt)
	assert.Equal(t, userID, cfg.UserID)
	assert.Equal(t, login, cfg.Login)
}

func TestClientConfig_ClearAuthTokens(t *testing.T) {
	cfg := DefaultClientConfig()
	cfg.AccessToken = "token"
	cfg.RefreshToken = "refresh"
	cfg.TokenExpiresAt = time.Now().Unix()
	cfg.UserID = "user-123"
	cfg.Login = "testuser"

	cfg.ClearAuthTokens()

	assert.Empty(t, cfg.AccessToken)
	assert.Empty(t, cfg.RefreshToken)
	assert.Equal(t, int64(0), cfg.TokenExpiresAt)
	assert.Empty(t, cfg.UserID)
	assert.Empty(t, cfg.Login)
}

func TestClientConfig_UpdateLastSyncTime(t *testing.T) {
	cfg := DefaultClientConfig()
	beforeUpdate := time.Now()

	cfg.UpdateLastSyncTime()

	assert.True(t, cfg.LastSyncTime >= beforeUpdate.Unix())
}

func TestClientConfig_GetLastSyncTime(t *testing.T) {
	cfg := DefaultClientConfig()
	syncTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	cfg.LastSyncTime = syncTime.Unix()

	result := cfg.GetLastSyncTime()

	assert.Equal(t, syncTime.Unix(), result.Unix())
}

func TestClientConfig_GetLastSyncTime_Zero(t *testing.T) {
	cfg := DefaultClientConfig()
	cfg.LastSyncTime = 0

	result := cfg.GetLastSyncTime()

	assert.Equal(t, time.Unix(0, 0), result)
}

// ========== File Operations Tests ==========

func TestGetConfigPath_Success(t *testing.T) {
	path, err := GetConfigPath()

	require.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.Contains(t, path, ".gophkeeper")
	assert.Contains(t, path, "config.json")
}

func TestSaveAndLoadClientConfig(t *testing.T) {
	// Создаем временную конфигурацию
	originalConfig := DefaultClientConfig()
	originalConfig.ServerAddress = "test.example.com:9000"
	originalConfig.AccessToken = "test_token"
	originalConfig.RefreshToken = "test_refresh"
	originalConfig.UserID = "user-456"
	originalConfig.Login = "testuser"

	// Сохраняем
	err := SaveClientConfig(originalConfig)
	require.NoError(t, err)

	// Загружаем
	loadedConfig, err := LoadClientConfig()
	require.NoError(t, err)

	// Проверяем
	assert.Equal(t, originalConfig.ServerAddress, loadedConfig.ServerAddress)
	assert.Equal(t, originalConfig.AccessToken, loadedConfig.AccessToken)
	assert.Equal(t, originalConfig.RefreshToken, loadedConfig.RefreshToken)
	assert.Equal(t, originalConfig.UserID, loadedConfig.UserID)
	assert.Equal(t, originalConfig.Login, loadedConfig.Login)
	assert.Equal(t, originalConfig.ConfigVersion, loadedConfig.ConfigVersion)

	// Cleanup
	configPath, _ := GetConfigPath()
	os.Remove(configPath)
}

func TestLoadClientConfig_CreateDefault(t *testing.T) {
	// Удаляем файл конфигурации если он существует
	configPath, err := GetConfigPath()
	require.NoError(t, err)
	os.Remove(configPath)

	// Загружаем конфигурацию (должна создаться по умолчанию)
	config, err := LoadClientConfig()
	require.NoError(t, err)

	// Проверяем, что создана конфигурация по умолчанию
	assert.Equal(t, "localhost:8080", config.ServerAddress)
	assert.Equal(t, "1.0", config.ConfigVersion)

	// Проверяем, что файл создан
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Cleanup
	os.Remove(configPath)
}

func TestSaveClientConfig_Persistence(t *testing.T) {
	// Создаем и сохраняем конфигурацию
	config1 := DefaultClientConfig()
	config1.ServerAddress = "server1.example.com"
	err := SaveClientConfig(config1)
	require.NoError(t, err)

	// Загружаем и изменяем
	config2, err := LoadClientConfig()
	require.NoError(t, err)
	config2.ServerAddress = "server2.example.com"
	err = SaveClientConfig(config2)
	require.NoError(t, err)

	// Загружаем снова и проверяем, что изменения сохранились
	config3, err := LoadClientConfig()
	require.NoError(t, err)
	assert.Equal(t, "server2.example.com", config3.ServerAddress)

	// Cleanup
	configPath, _ := GetConfigPath()
	os.Remove(configPath)
}

func TestSaveClientConfig_FilePermissions(t *testing.T) {
	config := DefaultClientConfig()
	err := SaveClientConfig(config)
	require.NoError(t, err)

	configPath, err := GetConfigPath()
	require.NoError(t, err)

	// Проверяем права доступа к файлу (должны быть 0600)
	info, err := os.Stat(configPath)
	require.NoError(t, err)

	mode := info.Mode().Perm()
	// На Unix-системах проверяем, что файл доступен только владельцу
	if os.Getenv("GOOS") != "windows" {
		assert.Equal(t, os.FileMode(0600), mode)
	}

	// Cleanup
	os.Remove(configPath)
}

func TestGetConfigPath_DirectoryCreation(t *testing.T) {
	// Получаем путь к конфигурации
	configPath, err := GetConfigPath()
	require.NoError(t, err)

	// Проверяем, что директория существует
	configDir := filepath.Dir(configPath)
	info, err := os.Stat(configDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Проверяем права доступа к директории (должны быть 0700)
	mode := info.Mode().Perm()
	if os.Getenv("GOOS") != "windows" {
		assert.Equal(t, os.FileMode(0700), mode)
	}
}

func TestClientConfig_FullAuthFlow(t *testing.T) {
	cfg := DefaultClientConfig()

	// Изначально не аутентифицирован
	assert.False(t, cfg.IsAuthenticated())

	// Устанавливаем токены
	expiresAt := time.Now().Add(time.Hour)
	cfg.SetAuthTokens("access", "refresh", expiresAt, "user-123", "testuser")

	// Теперь аутентифицирован
	assert.True(t, cfg.IsAuthenticated())
	assert.False(t, cfg.IsTokenExpired())

	// Сохраняем и загружаем
	err := SaveClientConfig(cfg)
	require.NoError(t, err)

	loaded, err := LoadClientConfig()
	require.NoError(t, err)

	// Проверяем, что аутентификация сохранилась
	assert.True(t, loaded.IsAuthenticated())
	assert.Equal(t, "testuser", loaded.Login)
	assert.Equal(t, "user-123", loaded.UserID)

	// Очищаем токены
	loaded.ClearAuthTokens()
	assert.False(t, loaded.IsAuthenticated())

	// Cleanup
	configPath, _ := GetConfigPath()
	os.Remove(configPath)
}
