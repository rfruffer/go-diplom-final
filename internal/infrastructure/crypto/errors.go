package crypto

import "errors"

// Ошибки валидации и криптографических операций
var (
	// ErrEmptyPassword возникает при попытке обработать пустой пароль.
	ErrEmptyPassword = errors.New("пароль не может быть пустым")

	// ErrWeakPassword возникает при попытке использовать слабый пароль.
	ErrWeakPassword = errors.New("пароль не соответствует требованиям безопасности")

	// ErrInvalidLength возникает при запросе некорректной длины данных.
	ErrInvalidLength = errors.New("некорректная длина данных")

	// ErrKeyDerivationFailed возникает при ошибке создания ключа из пароля.
	ErrKeyDerivationFailed = errors.New("не удалось создать ключ из пароля")

	// ErrEncryptionFailed возникает при ошибке шифрования данных.
	ErrEncryptionFailed = errors.New("не удалось зашифровать данные")

	// ErrDecryptionFailed возникает при ошибке расшифровки данных.
	ErrDecryptionFailed = errors.New("не удалось расшифровать данные")

	// ErrInvalidFormat возникает при некорректном формате зашифрованных данных.
	ErrInvalidFormat = errors.New("некорректный формат данных")
)