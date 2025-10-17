package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/fylgushev/go-diplom-final/internal/infrastructure/config"
	"github.com/fylgushev/go-diplom-final/pkg/proto/auth"
)

// Client представляет gRPC клиента
type Client struct {
	conn       *grpc.ClientConn
	authClient auth.AuthServiceClient
	config     *config.ClientConfig
	address    string
}

// NewClient создает новый gRPC клиент
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Connect устанавливает соединение с сервером
func (c *Client) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к серверу %s: %w", c.address, err)
	}

	c.conn = conn
	c.authClient = auth.NewAuthServiceClient(conn)

	return nil
}

// Close закрывает соединение с сервером
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// LoadConfig загружает конфигурацию клиента
func (c *Client) LoadConfig() error {
	config, err := config.LoadClientConfig()
	if err != nil {
		return fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
	}
	c.config = config
	return nil
}

// SaveConfig сохраняет конфигурацию клиента
func (c *Client) SaveConfig() error {
	if c.config == nil {
		return fmt.Errorf("конфигурация не загружена")
	}
	return config.SaveClientConfig(c.config)
}

// Register регистрирует нового пользователя
func (c *Client) Register(ctx context.Context, login, password string) error {
	req := &auth.RegisterRequest{
		Login:    login,
		Password: password,
	}

	resp, err := c.authClient.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка регистрации: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("регистрация неуспешна: %s", resp.Message)
	}

	// Сохраняем токены в конфигурации
	c.config.SetAuthTokens(
		resp.AccessToken,
		resp.RefreshToken,
		resp.ExpiresAt,
		"", // UserID будет получен при логине
		login,
	)

	return c.SaveConfig()
}

// Login аутентифицирует пользователя
func (c *Client) Login(ctx context.Context, login, password string) error {
	req := &auth.LoginRequest{
		Login:    login,
		Password: password,
	}

	resp, err := c.authClient.Login(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка входа: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("вход неуспешен: %s", resp.Message)
	}

	// Сохраняем токены и информацию о пользователе
	userID := ""
	if resp.User != nil {
		userID = resp.User.ID
	}

	c.config.SetAuthTokens(
		resp.AccessToken,
		resp.RefreshToken,
		resp.ExpiresAt,
		userID,
		login,
	)

	return c.SaveConfig()
}

// Logout выходит из системы
func (c *Client) Logout() error {
	c.config.ClearAuthTokens()
	return c.SaveConfig()
}

// RefreshTokens обновляет токены доступа
func (c *Client) RefreshTokens(ctx context.Context) error {
	if c.config.RefreshToken == "" {
		return fmt.Errorf("refresh token отсутствует")
	}

	req := &auth.RefreshTokenRequest{
		RefreshToken: c.config.RefreshToken,
	}

	resp, err := c.authClient.RefreshToken(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка обновления токенов: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("обновление токенов неуспешно: %s", resp.Message)
	}

	// Обновляем токены в конфигурации
	c.config.AccessToken = resp.AccessToken
	c.config.RefreshToken = resp.RefreshToken
	c.config.TokenExpiresAt = resp.ExpiresAt.Unix()

	return c.SaveConfig()
}

// addAuthHeader добавляет заголовок авторизации к контексту
func (c *Client) addAuthHeader(ctx context.Context) context.Context {
	if c.config.AccessToken != "" {
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + c.config.AccessToken,
		})
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

// ensureAuthenticated проверяет аутентификацию и обновляет токены при необходимости
func (c *Client) ensureAuthenticated(ctx context.Context) error {
	if !c.config.IsAuthenticated() {
		if c.config.RefreshToken == "" {
			return fmt.Errorf("не аутентифицирован. Используйте команду login")
		}

		// Пытаемся обновить токены
		if err := c.RefreshTokens(ctx); err != nil {
			return fmt.Errorf("не удалось обновить токены: %w. Используйте команду login", err)
		}
	}
	return nil
}

// GetConfig возвращает текущую конфигурацию
func (c *Client) GetConfig() *config.ClientConfig {
	return c.config
}