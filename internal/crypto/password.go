package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

// Константы для хеширования паролей
const (
	// DefaultCost стоимость bcrypt по умолчанию для хеширования паролей.
	// Значение 12 обеспечивает хороший баланс между безопасностью и производительностью.
	DefaultCost = 12

	// SaltSize размер соли в байтах для генерации ключей.
	// 16 байт (128 бит) достаточно для криптографической стойкости.
	SaltSize = 16

	// PBKDF2Iterations количество итераций для PBKDF2.
	// 100,000 итераций обеспечивают защиту от атак перебора.
	PBKDF2Iterations = 100000
)

// HashPassword хеширует пароль с использованием bcrypt.
// Возвращает хеш пароля и соль в одной строке.
// Использует стандартную стоимость bcrypt для безопасности.
// Предназначено для хранения паролей пользователей в базе данных.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// VerifyPassword проверяет, соответствует ли пароль его хешу.
// Использует bcrypt для безопасного сравнения пароля с хешем.
// Возвращает true, если пароль корректен, false в противном случае.
func VerifyPassword(password, hashedPassword string) bool {
	if password == "" || hashedPassword == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateSalt генерирует криптографически стойкую соль.
// Возвращает случайную последовательность байт размером SaltSize.
// Используется для генерации ключей из паролей через PBKDF2.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey создает ключ из пароля и соли с использованием PBKDF2.
// Использует SHA-256 как функцию хеширования и заданное количество итераций.
// keyLen определяет длину результирующего ключа в байтах.
// Используется для преобразования пользовательских паролей в ключи шифрования.
func DeriveKey(password string, salt []byte, keyLen int) []byte {
	return pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, keyLen, sha256.New)
}

// GenerateRandomBytes генерирует криптографически стойкую случайную последовательность.
// Принимает длину требуемой последовательности в байтах.
// Возвращает случайные байты или ошибку, если генерация не удалась.
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, ErrInvalidLength
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// CreateMasterKey создает мастер-ключ пользователя из его пароля.
// Генерирует уникальную соль для каждого пользователя и создает ключ через PBKDF2.
// Возвращает ключ и соль, которые должны быть сохранены для пользователя.
// Мастер-ключ используется для шифрования данных пользователя.
func CreateMasterKey(password string) (key []byte, salt []byte, err error) {
	if password == "" {
		return nil, nil, ErrEmptyPassword
	}

	// Генерируем уникальную соль для пользователя
	salt, err = GenerateSalt()
	if err != nil {
		return nil, nil, err
	}

	// Создаем мастер-ключ из пароля и соли
	key = DeriveKey(password, salt, 32) // 32 байта для AES-256

	return key, salt, nil
}

// RestoreMasterKey восстанавливает мастер-ключ пользователя из пароля и сохраненной соли.
// Используется при аутентификации пользователя для восстановления его ключа шифрования.
// Пароль должен совпадать с тем, который использовался при создании ключа.
func RestoreMasterKey(password string, salt []byte) []byte {
	return DeriveKey(password, salt, 32)
}

// ValidatePasswordStrength проверяет надежность пароля.
// Возвращает true, если пароль соответствует минимальным требованиям безопасности.
// Требования: минимум 8 символов, наличие букв и цифр.
func ValidatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		switch {
		case (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			 (char >= 'а' && char <= 'я') || (char >= 'А' && char <= 'Я'):
			hasLetter = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}