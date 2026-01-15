package infrastructure_test

import (
	"testing"

	"github.com/fylgushev/go-diplom-final/internal/infrastructure/crypto"
)

// TestEncodeDecodeBase64 проверяет кодирование и декодирование base64.
func TestEncodeDecodeBase64(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "простые данные",
			data: []byte("Hello, World!"),
		},
		{
			name: "пустые данные",
			data: []byte{},
		},
		{
			name: "бинарные данные",
			data: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
		},
		{
			name: "русский текст",
			data: []byte("Привет, мир!"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Кодируем
			encoded := crypto.EncodeBase64(tt.data)

			// Декодируем
			decoded, err := crypto.DecodeBase64(encoded)
			if err != nil {
				t.Errorf("DecodeBase64() ошибка: %v", err)
				return
			}

			// Проверяем, что данные восстановились
			if len(decoded) != len(tt.data) {
				t.Errorf("Неправильная длина после декодирования: получили %d, ожидали %d", len(decoded), len(tt.data))
			}

			for i, b := range decoded {
				if b != tt.data[i] {
					t.Errorf("Данные не совпадают на позиции %d: получили %d, ожидали %d", i, b, tt.data[i])
				}
			}
		})
	}
}

// TestDecodeBase64InvalidData проверяет обработку некорректных base64 данных.
func TestDecodeBase64InvalidData(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
	}{
		{
			name:    "некорректные символы",
			encoded: "Hello@World!",
		},
		{
			name:    "неправильная длина",
			encoded: "SGVsbG8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := crypto.DecodeBase64(tt.encoded)
			if err == nil {
				t.Errorf("DecodeBase64() ожидалась ошибка для некорректных данных")
			}
		})
	}
}

// TestCalculateVerifyChecksum проверяет вычисление и верификацию контрольных сумм.
func TestCalculateVerifyChecksum(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "простые данные",
			data: []byte("test data"),
		},
		{
			name: "пустые данные",
			data: []byte{},
		},
		{
			name: "бинарные данные",
			data: []byte{0x00, 0xFF, 0x12, 0x34, 0x56, 0x78},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Вычисляем контрольную сумму
			checksum := crypto.CalculateChecksum(tt.data)

			// Проверяем, что контрольная сумма не пустая
			if checksum == "" {
				t.Errorf("CalculateChecksum() вернул пустую строку")
			}

			// Проверяем, что контрольная сумма имеет правильную длину (SHA-256 в hex = 64 символа)
			if len(checksum) != 64 {
				t.Errorf("Неправильная длина контрольной суммы: получили %d, ожидали 64", len(checksum))
			}

			// Проверяем верификацию с правильной контрольной суммой
			if !crypto.VerifyChecksum(tt.data, checksum) {
				t.Errorf("VerifyChecksum() не подтвердил правильную контрольную сумму")
			}

			// Проверяем верификацию с неправильной контрольной суммой
			wrongChecksum := "0000000000000000000000000000000000000000000000000000000000000000"
			if crypto.VerifyChecksum(tt.data, wrongChecksum) {
				t.Errorf("VerifyChecksum() подтвердил неправильную контрольную сумму")
			}
		})
	}

	// Проверяем, что разные данные дают разные контрольные суммы
	checksum1 := crypto.CalculateChecksum([]byte("data1"))
	checksum2 := crypto.CalculateChecksum([]byte("data2"))
	if checksum1 == checksum2 {
		t.Errorf("Разные данные дали одинаковые контрольные суммы")
	}

	// Проверяем детерминированность
	data := []byte("consistent data")
	checksum1 = crypto.CalculateChecksum(data)
	checksum2 = crypto.CalculateChecksum(data)
	if checksum1 != checksum2 {
		t.Errorf("CalculateChecksum() не детерминирован")
	}
}

// TestSecureWipe проверяет безопасное затирание данных.
func TestSecureWipe(t *testing.T) {
	// Создаем тестовые данные
	data := []byte("sensitive data that should be wiped")
	originalLen := len(data)

	// Сохраняем копию для проверки
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	// Затираем данные
	crypto.SecureWipe(data)

	// Проверяем, что длина не изменилась
	if len(data) != originalLen {
		t.Errorf("SecureWipe() изменил длину среза")
	}

	// Проверяем, что все байты стали нулевыми
	for i, b := range data {
		if b != 0 {
			t.Errorf("SecureWipe() не затер байт на позиции %d: получили %d, ожидали 0", i, b)
		}
	}

	// Проверяем, что исходные данные действительно были другими
	allZeros := true
	for _, b := range dataCopy {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Errorf("Тестовые данные уже состояли из нулей")
	}
}

// TestCompareSecure проверяет безопасное сравнение данных.
func TestCompareSecure(t *testing.T) {
	tests := []struct {
		name     string
		a        []byte
		b        []byte
		expected bool
	}{
		{
			name:     "одинаковые данные",
			a:        []byte("test data"),
			b:        []byte("test data"),
			expected: true,
		},
		{
			name:     "разные данные",
			a:        []byte("test data"),
			b:        []byte("different"),
			expected: false,
		},
		{
			name:     "разные длины",
			a:        []byte("short"),
			b:        []byte("much longer data"),
			expected: false,
		},
		{
			name:     "пустые срезы",
			a:        []byte{},
			b:        []byte{},
			expected: true,
		},
		{
			name:     "один пустой",
			a:        []byte("data"),
			b:        []byte{},
			expected: false,
		},
		{
			name:     "бинарные данные одинаковые",
			a:        []byte{0x00, 0x01, 0xFF, 0xFE},
			b:        []byte{0x00, 0x01, 0xFF, 0xFE},
			expected: true,
		},
		{
			name:     "бинарные данные разные",
			a:        []byte{0x00, 0x01, 0xFF, 0xFE},
			b:        []byte{0x00, 0x01, 0xFF, 0xFD},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crypto.CompareSecure(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareSecure() = %v, ожидалось %v", result, tt.expected)
			}
		})
	}
}

// TestGenerateRandomString проверяет генерацию случайных строк.
func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "короткая строка",
			length: 8,
		},
		{
			name:   "средняя строка",
			length: 16,
		},
		{
			name:   "длинная строка",
			length: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str, err := crypto.GenerateRandomString(tt.length)
			if err != nil {
				t.Errorf("GenerateRandomString() ошибка: %v", err)
				return
			}

			// Проверяем длину (может быть чуть меньше из-за обрезки base64)
			if len(str) > tt.length {
				t.Errorf("Строка слишком длинная: получили %d, максимум %d", len(str), tt.length)
			}
			if len(str) < tt.length-3 { // допускаем небольшую погрешность
				t.Errorf("Строка слишком короткая: получили %d, ожидали около %d", len(str), tt.length)
			}

			// Проверяем, что строка не пустая
			if str == "" {
				t.Errorf("GenerateRandomString() вернул пустую строку")
			}
		})
	}

	// Проверяем уникальность
	strings := make([]string, 10)
	for i := 0; i < 10; i++ {
		str, err := crypto.GenerateRandomString(20)
		if err != nil {
			t.Fatalf("Ошибка генерации строки: %v", err)
		}
		strings[i] = str
	}

	// Проверяем, что все строки разные
	for i := 0; i < len(strings); i++ {
		for j := i + 1; j < len(strings); j++ {
			if strings[i] == strings[j] {
				t.Errorf("GenerateRandomString() сгенерировал одинаковые строки")
			}
		}
	}
}

// TestIsEncrypted проверяет определение зашифрованных данных.
func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "пустые данные",
			data:     []byte{},
			expected: false,
		},
		{
			name:     "обычный текст",
			data:     []byte("Hello, World! This is plain text."),
			expected: false,
		},
		{
			name:     "повторяющиеся символы",
			data:     []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			expected: false,
		},
		{
			name:     "псевдослучайные данные",
			data:     []byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x70, 0x81, 0x92, 0xa3, 0xb4, 0xc5, 0xd6, 0xe7, 0xf8, 0x09},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crypto.IsEncrypted(tt.data)
			if result != tt.expected {
				t.Errorf("IsEncrypted() = %v, ожидалось %v", result, tt.expected)
			}
		})
	}

	// Проверяем на реально зашифрованных данных
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Не удалось сгенерировать ключ: %v", err)
	}

	cipher, err := crypto.NewAESCipher(key)
	if err != nil {
		t.Fatalf("Не удалось создать cipher: %v", err)
	}

	plaintext := []byte("This is a test message that will be encrypted")
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Не удалось зашифровать данные: %v", err)
	}

	// Зашифрованные данные должны определяться как зашифрованные
	if !crypto.IsEncrypted(ciphertext) {
		t.Errorf("IsEncrypted() не определил зашифрованные данные")
	}

	// Исходные данные не должны определяться как зашифрованные
	if crypto.IsEncrypted(plaintext) {
		t.Errorf("IsEncrypted() определил обычный текст как зашифрованный")
	}
}

// TestValidateKey проверяет валидацию ключей.
func TestValidateKey(t *testing.T) {
	tests := []struct {
		name     string
		keyLen   int
		expected bool
	}{
		{
			name:     "AES-128 ключ",
			keyLen:   16,
			expected: true,
		},
		{
			name:     "AES-192 ключ",
			keyLen:   24,
			expected: true,
		},
		{
			name:     "AES-256 ключ",
			keyLen:   32,
			expected: true,
		},
		{
			name:     "слишком короткий ключ",
			keyLen:   8,
			expected: false,
		},
		{
			name:     "слишком длинный ключ",
			keyLen:   64,
			expected: false,
		},
		{
			name:     "пустой ключ",
			keyLen:   0,
			expected: false,
		},
		{
			name:     "нестандартный размер",
			keyLen:   20,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keyLen)
			result := crypto.ValidateKey(key)
			if result != tt.expected {
				t.Errorf("ValidateKey() для ключа размером %d байт = %v, ожидалось %v", tt.keyLen, result, tt.expected)
			}
		})
	}
}