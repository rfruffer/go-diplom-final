package infrastructure_test

import (
	"bytes"
	"testing"

	"github.com/fylgushev/go-diplom-final/internal/infrastructure/crypto"
)

// TestHashPassword проверяет хеширование паролей.
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "обычный пароль",
			password: "mypassword123",
			wantErr:  false,
		},
		{
			name:     "сложный пароль",
			password: "Complex!Password@2024#",
			wantErr:  false,
		},
		{
			name:     "русский пароль",
			password: "русский_пароль_123",
			wantErr:  false,
		},
		{
			name:     "пустой пароль",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := crypto.HashPassword(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("HashPassword() ожидалась ошибка для пустого пароля")
				}
				return
			}

			if err != nil {
				t.Errorf("HashPassword() неожиданная ошибка: %v", err)
				return
			}

			// Проверяем, что хеш не пустой
			if hashedPassword == "" {
				t.Errorf("HashPassword() вернул пустой хеш")
			}

			// Проверяем, что хеш отличается от исходного пароля
			if hashedPassword == tt.password {
				t.Errorf("Хеш не должен совпадать с исходным паролем")
			}

			// Проверяем, что хеш выглядит как bcrypt хеш
			if len(hashedPassword) < 50 {
				t.Errorf("Хеш слишком короткий, возможно это не bcrypt")
			}
		})
	}
}

// TestVerifyPassword проверяет верификацию паролей.
func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("Не удалось создать хеш пароля: %v", err)
	}

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expected       bool
	}{
		{
			name:           "правильный пароль",
			password:       password,
			hashedPassword: hashedPassword,
			expected:       true,
		},
		{
			name:           "неправильный пароль",
			password:       "wrongpassword",
			hashedPassword: hashedPassword,
			expected:       false,
		},
		{
			name:           "пустой пароль",
			password:       "",
			hashedPassword: hashedPassword,
			expected:       false,
		},
		{
			name:           "пустой хеш",
			password:       password,
			hashedPassword: "",
			expected:       false,
		},
		{
			name:           "регистрозависимость",
			password:       "TESTPASSWORD123",
			hashedPassword: hashedPassword,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crypto.VerifyPassword(tt.password, tt.hashedPassword)
			if result != tt.expected {
				t.Errorf("VerifyPassword() = %v, ожидалось %v", result, tt.expected)
			}
		})
	}
}

// TestGenerateSalt проверяет генерацию соли.
func TestGenerateSalt(t *testing.T) {
	// Генерируем несколько солей
	salts := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		salt, err := crypto.GenerateSalt()
		if err != nil {
			t.Fatalf("GenerateSalt() ошибка: %v", err)
		}

		// Проверяем размер соли
		if len(salt) != crypto.SaltSize {
			t.Errorf("GenerateSalt() неправильный размер соли: получили %d, ожидали %d", len(salt), crypto.SaltSize)
		}

		salts[i] = salt
	}

	// Проверяем, что соли разные
	for i := 0; i < len(salts); i++ {
		for j := i + 1; j < len(salts); j++ {
			if bytes.Equal(salts[i], salts[j]) {
				t.Errorf("GenerateSalt() сгенерировал одинаковые соли")
			}
		}
	}
}

// TestDeriveKey проверяет создание ключей из паролей.
func TestDeriveKey(t *testing.T) {
	password := "testpassword"
	salt := []byte("testsalt12345678") // 16 байт

	tests := []struct {
		name   string
		keyLen int
	}{
		{
			name:   "ключ 16 байт",
			keyLen: 16,
		},
		{
			name:   "ключ 32 байта",
			keyLen: 32,
		},
		{
			name:   "ключ 64 байта",
			keyLen: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := crypto.DeriveKey(password, salt, tt.keyLen)

			// Проверяем длину ключа
			if len(key) != tt.keyLen {
				t.Errorf("DeriveKey() неправильная длина ключа: получили %d, ожидали %d", len(key), tt.keyLen)
			}

			// Проверяем, что ключ не состоит из одних нулей
			allZeros := true
			for _, b := range key {
				if b != 0 {
					allZeros = false
					break
				}
			}
			if allZeros {
				t.Errorf("DeriveKey() сгенерировал ключ из одних нулей")
			}
		})
	}

	// Проверяем детерминированность
	key1 := crypto.DeriveKey(password, salt, 32)
	key2 := crypto.DeriveKey(password, salt, 32)
	if !bytes.Equal(key1, key2) {
		t.Errorf("DeriveKey() должен быть детерминированным")
	}

	// Проверяем, что разные пароли дают разные ключи
	key3 := crypto.DeriveKey("differentpassword", salt, 32)
	if bytes.Equal(key1, key3) {
		t.Errorf("Разные пароли не должны давать одинаковые ключи")
	}

	// Проверяем, что разные соли дают разные ключи
	differentSalt := []byte("differentsalt123")
	key4 := crypto.DeriveKey(password, differentSalt, 32)
	if bytes.Equal(key1, key4) {
		t.Errorf("Разные соли не должны давать одинаковые ключи")
	}
}

// TestGenerateRandomBytes проверяет генерацию случайных байт.
func TestGenerateRandomBytes(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "16 байт",
			length:  16,
			wantErr: false,
		},
		{
			name:    "32 байта",
			length:  32,
			wantErr: false,
		},
		{
			name:    "100 байт",
			length:  100,
			wantErr: false,
		},
		{
			name:    "нулевая длина",
			length:  0,
			wantErr: true,
		},
		{
			name:    "отрицательная длина",
			length:  -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := crypto.GenerateRandomBytes(tt.length)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateRandomBytes() ожидалась ошибка")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateRandomBytes() неожиданная ошибка: %v", err)
				return
			}

			if len(bytes) != tt.length {
				t.Errorf("GenerateRandomBytes() неправильная длина: получили %d, ожидали %d", len(bytes), tt.length)
			}
		})
	}

	// Проверяем случайность
	sequences := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		seq, err := crypto.GenerateRandomBytes(32)
		if err != nil {
			t.Fatalf("Ошибка генерации случайных байт: %v", err)
		}
		sequences[i] = seq
	}

	// Проверяем, что последовательности разные
	for i := 0; i < len(sequences); i++ {
		for j := i + 1; j < len(sequences); j++ {
			if bytes.Equal(sequences[i], sequences[j]) {
				t.Errorf("GenerateRandomBytes() сгенерировал одинаковые последовательности")
			}
		}
	}
}

// TestCreateMasterKey проверяет создание мастер-ключа.
func TestCreateMasterKey(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "обычный пароль",
			password: "userpassword123",
			wantErr:  false,
		},
		{
			name:     "сложный пароль",
			password: "Very!Complex@Password#2024",
			wantErr:  false,
		},
		{
			name:     "пустой пароль",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, salt, err := crypto.CreateMasterKey(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateMasterKey() ожидалась ошибка для пустого пароля")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateMasterKey() неожиданная ошибка: %v", err)
				return
			}

			// Проверяем размеры
			if len(key) != 32 {
				t.Errorf("Неправильный размер ключа: получили %d, ожидали 32", len(key))
			}
			if len(salt) != crypto.SaltSize {
				t.Errorf("Неправильный размер соли: получили %d, ожидали %d", len(salt), crypto.SaltSize)
			}

			// Проверяем, что можно восстановить тот же ключ
			restoredKey := crypto.RestoreMasterKey(tt.password, salt)
			if !bytes.Equal(key, restoredKey) {
				t.Errorf("Восстановленный ключ не совпадает с исходным")
			}
		})
	}
}

// TestValidatePasswordStrength проверяет валидацию надежности пароля.
func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "надежный пароль",
			password: "password123",
			expected: true,
		},
		{
			name:     "сложный пароль",
			password: "MyStrongPassword2024",
			expected: true,
		},
		{
			name:     "слишком короткий",
			password: "pass1",
			expected: false,
		},
		{
			name:     "только буквы",
			password: "passwordwithoutdigits",
			expected: false,
		},
		{
			name:     "только цифры",
			password: "12345678",
			expected: false,
		},
		{
			name:     "пустой пароль",
			password: "",
			expected: false,
		},
		{
			name:     "минимально допустимый",
			password: "abcdefg1",
			expected: true,
		},
		{
			name:     "с заглавными буквами",
			password: "Password123",
			expected: true,
		},
		{
			name:     "русские буквы с цифрами",
			password: "пароль123",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crypto.ValidatePasswordStrength(tt.password)
			if result != tt.expected {
				t.Errorf("ValidatePasswordStrength(%q) = %v, ожидалось %v", tt.password, result, tt.expected)
			}
		})
	}
}