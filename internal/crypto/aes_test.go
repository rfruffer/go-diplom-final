package crypto

import (
	"bytes"
	"testing"
)

// TestNewAESCipher проверяет создание нового AES шифра.
func TestNewAESCipher(t *testing.T) {
	tests := []struct {
		name    string
		keyLen  int
		wantErr bool
	}{
		{
			name:    "валидный ключ 256 бит",
			keyLen:  32,
			wantErr: false,
		},
		{
			name:    "невалидный ключ 128 бит",
			keyLen:  16,
			wantErr: true,
		},
		{
			name:    "невалидный ключ 0 байт",
			keyLen:  0,
			wantErr: true,
		},
		{
			name:    "невалидный ключ 64 байта",
			keyLen:  64,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keyLen)
			cipher, err := NewAESCipher(key)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewAESCipher() ожидалась ошибка, но получили nil")
				}
				if cipher != nil {
					t.Errorf("NewAESCipher() ожидался nil, но получили cipher")
				}
			} else {
				if err != nil {
					t.Errorf("NewAESCipher() неожиданная ошибка: %v", err)
				}
				if cipher == nil {
					t.Errorf("NewAESCipher() ожидался cipher, но получили nil")
				}
			}
		})
	}
}

// TestAESCipher_EncryptDecrypt проверяет шифрование и расшифровку.
func TestAESCipher_EncryptDecrypt(t *testing.T) {
	// Создаем тестовый ключ
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	cipher, err := NewAESCipher(key)
	if err != nil {
		t.Fatalf("Не удалось создать cipher: %v", err)
	}

	tests := []struct {
		name      string
		plaintext []byte
		wantErr   bool
	}{
		{
			name:      "обычный текст",
			plaintext: []byte("Hello, World!"),
			wantErr:   false,
		},
		{
			name:      "пустые данные",
			plaintext: []byte{},
			wantErr:   true,
		},
		{
			name:      "длинный текст",
			plaintext: []byte("Это очень длинный текст для проверки шифрования больших данных в системе GophKeeper"),
			wantErr:   false,
		},
		{
			name:      "бинарные данные",
			plaintext: []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Шифруем
			ciphertext, err := cipher.Encrypt(tt.plaintext)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Encrypt() ожидалась ошибка, но получили nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Encrypt() неожиданная ошибка: %v", err)
				return
			}

			// Проверяем, что зашифрованные данные отличаются от исходных
			if bytes.Equal(ciphertext, tt.plaintext) {
				t.Errorf("Зашифрованные данные не должны совпадать с исходными")
			}

			// Расшифровываем
			decrypted, err := cipher.Decrypt(ciphertext)
			if err != nil {
				t.Errorf("Decrypt() неожиданная ошибка: %v", err)
				return
			}

			// Проверяем, что расшифрованные данные совпадают с исходными
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Расшифрованные данные не совпадают с исходными")
				t.Errorf("Ожидалось: %v", tt.plaintext)
				t.Errorf("Получено: %v", decrypted)
			}
		})
	}
}

// TestAESCipher_DecryptInvalidData проверяет обработку некорректных данных при расшифровке.
func TestAESCipher_DecryptInvalidData(t *testing.T) {
	key := make([]byte, 32)
	cipher, err := NewAESCipher(key)
	if err != nil {
		t.Fatalf("Не удалось создать cipher: %v", err)
	}

	tests := []struct {
		name       string
		ciphertext []byte
	}{
		{
			name:       "слишком короткие данные",
			ciphertext: []byte{0x01, 0x02},
		},
		{
			name:       "пустые данные",
			ciphertext: []byte{},
		},
		{
			name:       "поврежденные данные",
			ciphertext: make([]byte, 50), // нули
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cipher.Decrypt(tt.ciphertext)
			if err == nil {
				t.Errorf("Decrypt() ожидалась ошибка для некорректных данных")
			}
		})
	}
}

// TestGenerateKey проверяет генерацию ключей.
func TestGenerateKey(t *testing.T) {
	// Генерируем несколько ключей
	keys := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		key, err := GenerateKey()
		if err != nil {
			t.Fatalf("GenerateKey() ошибка: %v", err)
		}

		// Проверяем длину ключа
		if len(key) != 32 {
			t.Errorf("GenerateKey() неправильная длина ключа: получили %d, ожидали 32", len(key))
		}

		keys[i] = key
	}

	// Проверяем, что ключи разные
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if bytes.Equal(keys[i], keys[j]) {
				t.Errorf("GenerateKey() сгенерировал одинаковые ключи")
			}
		}
	}
}

// TestEncryptDecryptWithPassword проверяет шифрование с паролем.
func TestEncryptDecryptWithPassword(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
		password  string
		wantErr   bool
	}{
		{
			name:      "обычные данные с паролем",
			plaintext: []byte("Секретные данные пользователя"),
			password:  "strong_password_123",
			wantErr:   false,
		},
		{
			name:      "пустой пароль",
			plaintext: []byte("данные"),
			password:  "",
			wantErr:   true,
		},
		{
			name:      "пустые данные",
			plaintext: []byte{},
			password:  "password",
			wantErr:   true,
		},
		{
			name:      "русский пароль",
			plaintext: []byte("тестовые данные"),
			password:  "русский_пароль_тест",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Шифруем
			ciphertext, err := EncryptWithPassword(tt.plaintext, tt.password)
			if tt.wantErr {
				if err == nil {
					t.Errorf("EncryptWithPassword() ожидалась ошибка")
				}
				return
			}

			if err != nil {
				t.Errorf("EncryptWithPassword() неожиданная ошибка: %v", err)
				return
			}

			// Расшифровываем с правильным паролем
			decrypted, err := DecryptWithPassword(ciphertext, tt.password)
			if err != nil {
				t.Errorf("DecryptWithPassword() неожиданная ошибка: %v", err)
				return
			}

			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Расшифрованные данные не совпадают с исходными")
			}

			// Проверяем с неправильным паролем
			_, err = DecryptWithPassword(ciphertext, "wrong_password")
			if err == nil {
				t.Errorf("DecryptWithPassword() должен возвращать ошибку для неправильного пароля")
			}
		})
	}
}

// TestEncryptionSecurity проверяет безопасность шифрования.
func TestEncryptionSecurity(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Не удалось сгенерировать ключ: %v", err)
	}

	cipher, err := NewAESCipher(key)
	if err != nil {
		t.Fatalf("Не удалось создать cipher: %v", err)
	}

	plaintext := []byte("одинаковые данные")

	// Шифруем одни и те же данные несколько раз
	ciphertexts := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		ciphertext, err := cipher.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Ошибка шифрования: %v", err)
		}
		ciphertexts[i] = ciphertext
	}

	// Проверяем, что результаты шифрования разные (благодаря случайному nonce)
	for i := 0; i < len(ciphertexts); i++ {
		for j := i + 1; j < len(ciphertexts); j++ {
			if bytes.Equal(ciphertexts[i], ciphertexts[j]) {
				t.Errorf("Шифрование одинаковых данных дало одинаковый результат (nonce не работает)")
			}
		}
	}

	// Проверяем, что все результаты корректно расшифровываются
	for i, ciphertext := range ciphertexts {
		decrypted, err := cipher.Decrypt(ciphertext)
		if err != nil {
			t.Errorf("Ошибка расшифровки %d: %v", i, err)
		}
		if !bytes.Equal(decrypted, plaintext) {
			t.Errorf("Неправильная расшифровка %d", i)
		}
	}
}