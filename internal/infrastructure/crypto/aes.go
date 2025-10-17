// Package crypto предоставляет криптографические функции для GophKeeper.
// Этот пакет реализует шифрование данных пользователей, хеширование паролей
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Cipher интерфейс для шифрования и расшифровки данных
type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// Ошибки шифрования
var (
	// ErrInvalidKeySize возникает при неподдерживаемом размере ключа шифрования.
	ErrInvalidKeySize = errors.New("неподдерживаемый размер ключа")

	// ErrInvalidCiphertext возникает при попытке расшифровать некорректные данные.
	ErrInvalidCiphertext = errors.New("некорректные зашифрованные данные")

	// ErrEmptyPlaintext возникает при попытке зашифровать пустые данные.
	ErrEmptyPlaintext = errors.New("нельзя шифровать пустые данные")
)

// AESCipher представляет шифр AES для шифрования и расшифровки данных.
// Использует AES-256 в режиме GCM для обеспечения конфиденциальности
type AESCipher struct {
	// key ключ шифрования AES-256 (32 байта).
	key []byte
}

// NewAESCipher создает новый экземпляр AES шифра.
// Принимает ключ размером 32 байта для AES-256.
// Возвращает ошибку, если размер ключа некорректен.
func NewAESCipher(key []byte) (*AESCipher, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	return &AESCipher{
		key: key,
	}, nil
}

// Encrypt шифрует plaintext с использованием AES-256-GCM.
// Возвращает зашифрованные данные с встроенным nonce.
// Данные включают: [nonce][ciphertext+tag]
// Обеспечивает как конфиденциальность, так и целостность данных.
func (c *AESCipher) Encrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, ErrEmptyPlaintext
	}

	// Создаем AES блок-шифр
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	// Создаем GCM режим для аутентифицированного шифрования
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Генерируем случайный nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Шифруем данные
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Возвращаем nonce + ciphertext
	result := make([]byte, len(nonce)+len(ciphertext))
	copy(result[:len(nonce)], nonce)
	copy(result[len(nonce):], ciphertext)

	return result, nil
}

// Decrypt расшифровывает ciphertext, зашифрованный с помощью Encrypt.
// Ожидает данные в формате: [nonce][ciphertext+tag]
// Проверяет целостность данных и возвращает оригинальный plaintext.
func (c *AESCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	// Создаем AES блок-шифр
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	// Создаем GCM режим
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Проверяем минимальную длину (nonce + tag)
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize+gcm.Overhead() {
		return nil, ErrInvalidCiphertext
	}

	// Извлекаем nonce и зашифрованные данные
	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	// Расшифровываем и проверяем целостность
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	return plaintext, nil
}

// GenerateKey генерирует случайный 256-битный ключ для AES-256.
// Использует криптографически стойкий генератор случайных чисел.
// Возвращает ключ размером 32 байта.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // 256 бит = 32 байта
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptWithPassword шифрует данные с использованием пароля.
// Пароль преобразуется в ключ с помощью PBKDF2.
// Возвращает зашифрованные данные с встроенной солью.
func EncryptWithPassword(plaintext []byte, password string) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, ErrEmptyPlaintext
	}
	if password == "" {
		return nil, errors.New("пароль не может быть пустым")
	}

	// Генерируем соль
	salt, err := GenerateSalt()
	if err != nil {
		return nil, err
	}

	// Создаем ключ из пароля
	key := DeriveKey(password, salt, 32)

	// Создаем AES шифр
	cipher, err := NewAESCipher(key)
	if err != nil {
		return nil, err
	}

	// Шифруем данные
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	// Возвращаем соль + зашифрованные данные
	result := make([]byte, len(salt)+len(ciphertext))
	copy(result[:len(salt)], salt)
	copy(result[len(salt):], ciphertext)

	return result, nil
}

// DecryptWithPassword расшифровывает данные, зашифрованные с помощью EncryptWithPassword.
// Извлекает соль из данных и восстанавливает ключ из пароля.
func DecryptWithPassword(ciphertext []byte, password string) ([]byte, error) {
	if password == "" {
		return nil, errors.New("пароль не может быть пустым")
	}

	// Проверяем минимальную длину (соль + nonce + tag)
	saltSize := 16                        // размер соли
	if len(ciphertext) < saltSize+12+16 { // 12 - nonce, 16 - GCM tag
		return nil, ErrInvalidCiphertext
	}

	// Извлекаем соль и зашифрованные данные
	salt := ciphertext[:saltSize]
	encryptedData := ciphertext[saltSize:]

	// Восстанавливаем ключ из пароля
	key := DeriveKey(password, salt, 32)

	// Создаем AES шифр
	cipher, err := NewAESCipher(key)
	if err != nil {
		return nil, err
	}

	// Расшифровываем данные
	return cipher.Decrypt(encryptedData)
}
