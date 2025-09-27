package models

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя системы GophKeeper.
// Содержит основную информацию о пользователе, включая идентификатор,
// учетные данные и временные метки для аудита.
type User struct {
	// ID уникальный идентификатор пользователя.
	// Генерируется автоматически при создании пользователя с использованием UUID.
	ID string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Username имя пользователя для аутентификации.
	// Должно быть уникальным в системе и содержать от 3 до 50 символов.
	Username string `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`

	// PasswordHash хеш пароля пользователя.
	// Пароль никогда не хранится в открытом виде, только его криптографический хеш.
	PasswordHash string `json:"-" gorm:"type:varchar(255);not null"`

	// Salt соль для хеширования пароля.
	// Уникальная случайная строка, добавляемая к паролю перед хешированием
	// для защиты от атак по радужным таблицам.
	Salt string `json:"-" gorm:"type:varchar(255);not null"`

	// CreatedAt время создания пользователя.
	// Автоматически устанавливается при регистрации пользователя.
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt время последнего обновления данных пользователя.
	// Автоматически обновляется при любых изменениях записи пользователя.
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// LastLoginAt время последнего входа пользователя в систему.
	// Используется для аудита и безопасности, может быть nil если пользователь никогда не входил.
	LastLoginAt *time.Time `json:"last_login_at" gorm:"type:timestamp"`

	// IsActive флаг активности пользователя.
	// Позволяет временно отключать пользователей без удаления их данных.
	IsActive bool `json:"is_active" gorm:"default:true"`
}

// NewUser создает новый экземпляр пользователя с базовыми значениями.
// Генерирует уникальный ID и устанавливает пользователя как активного.
// Время создания и обновления будут установлены автоматически ORM.
func NewUser(username string) *User {
	return &User{
		ID:       uuid.New().String(),
		Username: username,
		IsActive: true,
	}
}

// SetPassword устанавливает хеш пароля и соль для пользователя.
// Используется при регистрации и смене пароля.
// passwordHash должен быть уже захешированным паролем с солью.
func (u *User) SetPassword(passwordHash, salt string) {
	u.PasswordHash = passwordHash
	u.Salt = salt
}

// UpdateLastLogin обновляет время последнего входа пользователя.
// Вызывается при успешной аутентификации пользователя.
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// Deactivate деактивирует пользователя.
// Деактивированный пользователь не может войти в систему,
// но его данные сохраняются в базе данных.
func (u *User) Deactivate() {
	u.IsActive = false
}

// Activate активирует пользователя.
// Активирует ранее деактивированного пользователя.
func (u *User) Activate() {
	u.IsActive = true
}

// TableName возвращает название таблицы в базе данных для модели User.
// Используется ORM GORM для корректного маппинга модели на таблицу.
func (User) TableName() string {
	return "users"
}