package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// EncodeBase64 кодирует данные в base64.
// Используется для безопасного представления бинарных данных в текстовом виде.
// Применяется для хранения зашифрованных данных в JSON или базе данных.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 декодирует данные из base64.
// Используется для восстановления бинарных данных из текстового представления.
// Возвращает ошибку, если данные имеют некорректный формат base64.
func DecodeBase64(encoded string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	return data, nil
}

// CalculateChecksum вычисляет SHA-256 хеш данных.
// Используется для проверки целостности данных и создания контрольных сумм.
// Возвращает хеш в виде hex-строки для удобного хранения и сравнения.
func CalculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// VerifyChecksum проверяет соответствие данных их контрольной сумме.
// Вычисляет SHA-256 хеш данных и сравнивает с ожидаемой контрольной суммой.
// Возвращает true, если данные не были изменены.
func VerifyChecksum(data []byte, expectedChecksum string) bool {
	actualChecksum := CalculateChecksum(data)
	return actualChecksum == expectedChecksum
}

// SecureWipe затирает содержимое среза байт нулями.
// Используется для безопасного удаления конфиденциальных данных из памяти
// (паролей, ключей шифрования) после их использования.
func SecureWipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

// CompareSecure выполняет безопасное сравнение двух срезов байт.
// Использует constant-time сравнение для защиты от timing attacks.
// Возвращает true, если срезы идентичны.
func CompareSecure(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// GenerateRandomString генерирует случайную строку заданной длины.
// Использует криптографически стойкий генератор для создания случайных данных.
// Результат кодируется в base64 для получения читаемой строки.
func GenerateRandomString(length int) (string, error) {
	// Вычисляем количество байт для получения строки нужной длины
	// base64 увеличивает размер примерно в 4/3 раза
	byteLength := (length * 3) / 4
	if byteLength == 0 {
		byteLength = 1
	}

	bytes, err := GenerateRandomBytes(byteLength)
	if err != nil {
		return "", err
	}

	// Кодируем в base64 и обрезаем до нужной длины
	encoded := EncodeBase64(bytes)
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}

// IsEncrypted проверяет, выглядят ли данные как зашифрованные.
// Выполняет эвристическую проверку на основе энтропии данных.
// Возвращает true, если данные, вероятно, зашифрованы.
func IsEncrypted(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Проверяем, есть ли явные признаки текста
	textBytes := 0
	for _, b := range data {
		// ASCII текст, пробелы, переводы строк
		if (b >= 32 && b <= 126) || b == '\n' || b == '\r' || b == '\t' {
			textBytes++
		}
	}

	// Если более 80% байтов выглядят как текст, это вероятно не зашифровано
	if float64(textBytes)/float64(len(data)) > 0.8 {
		return false
	}

	// Простая эвристика: проверяем распределение байтов
	// Зашифрованные данные должны иметь более равномерное распределение
	frequency := make(map[byte]int)
	for _, b := range data {
		frequency[b]++
	}

	// Если какой-то байт встречается слишком часто, данные, вероятно, не зашифрованы
	maxFrequency := len(data) / 3 // 33% от общего количества
	for _, count := range frequency {
		if count > maxFrequency {
			return false
		}
	}

	return true
}

// ValidateKey проверяет, что ключ имеет корректный размер для AES.
// Поддерживает ключи размером 16, 24 или 32 байта (AES-128, AES-192, AES-256).
// Возвращает true, если размер ключа валиден.
func ValidateKey(key []byte) bool {
	keyLen := len(key)
	return keyLen == 16 || keyLen == 24 || keyLen == 32
}