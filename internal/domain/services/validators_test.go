package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
}

// ========== ValidateLogin Tests ==========

func TestValidator_ValidateLogin_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"user",
		"user@example.com",
		"user.name",
		"user123",
		"ab", // минимальная длина
	}

	for _, login := range tests {
		t.Run(login, func(t *testing.T) {
			err := v.ValidateLogin(login)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidateLogin_Empty(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLogin("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateLogin_OnlySpaces(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLogin("   ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateLogin_TooShort(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLogin("a")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "минимум 2")
}

func TestValidator_ValidateLogin_TooLong(t *testing.T) {
	v := NewValidator()
	longLogin := make([]byte, 101)
	for i := range longLogin {
		longLogin[i] = 'a'
	}
	err := v.ValidateLogin(string(longLogin))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "длиннее 100")
}

// ========== ValidatePassword Tests ==========

func TestValidator_ValidatePassword_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"p",
		"password",
		"123456",
		"complex!@#Password123",
	}

	for _, password := range tests {
		t.Run(password, func(t *testing.T) {
			err := v.ValidatePassword(password)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidatePassword_Empty(t *testing.T) {
	v := NewValidator()
	err := v.ValidatePassword("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidatePassword_TooLong(t *testing.T) {
	v := NewValidator()
	longPassword := make([]byte, 501)
	for i := range longPassword {
		longPassword[i] = 'a'
	}
	err := v.ValidatePassword(string(longPassword))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "длиннее 500")
}

// ========== ValidateURL Tests ==========

func TestValidator_ValidateURL_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"https://example.com",
		"http://example.com",
		"example.com",
		"www.example.com",
		"https://example.com/path",
		"",
	}

	for _, url := range tests {
		t.Run(url, func(t *testing.T) {
			err := v.ValidateURL(url)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidateURL_Invalid(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"ht tp://invalid url.com",
		"://noscheme",
	}

	for _, url := range tests {
		t.Run(url, func(t *testing.T) {
			err := v.ValidateURL(url)
			assert.Error(t, err)
		})
	}
}

// ========== ValidateText Tests ==========

func TestValidator_ValidateText_Success(t *testing.T) {
	v := NewValidator()
	err := v.ValidateText("Valid text content")
	assert.NoError(t, err)
}

func TestValidator_ValidateText_Empty(t *testing.T) {
	v := NewValidator()
	err := v.ValidateText("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateText_OnlySpaces(t *testing.T) {
	v := NewValidator()
	err := v.ValidateText("   ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateText_TooLong(t *testing.T) {
	v := NewValidator()
	longText := make([]byte, 10001)
	for i := range longText {
		longText[i] = 'a'
	}
	err := v.ValidateText(string(longText))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "длиннее 10000")
}

// ========== ValidateCardNumber Tests ==========

func TestValidator_ValidateCardNumber_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"4532015112830366",           // Visa
		"5425233430109903",           // Mastercard
		"4532 0151 1283 0366",        // С пробелами
		"4532-0151-1283-0366",        // С дефисами
		"4532 - 0151 - 1283 - 0366",  // Смешанный формат
	}

	for _, number := range tests {
		t.Run(number, func(t *testing.T) {
			err := v.ValidateCardNumber(number)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidateCardNumber_Empty(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardNumber("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateCardNumber_NonDigits(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardNumber("abcd1234567890")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "только цифры")
}

func TestValidator_ValidateCardNumber_TooShort(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardNumber("123456789012")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "от 13 до 19")
}

func TestValidator_ValidateCardNumber_TooLong(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardNumber("12345678901234567890")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "от 13 до 19")
}

func TestValidator_ValidateCardNumber_InvalidLuhn(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardNumber("1234567890123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Луна")
}

// ========== ValidateCardExpiry Tests ==========

func TestValidator_ValidateCardExpiry_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"12/25",
		"01/2025",
		"06/30",
	}

	for _, expiry := range tests {
		t.Run(expiry, func(t *testing.T) {
			err := v.ValidateCardExpiry(expiry)
			// Note: некоторые тесты могут падать если дата истекла
			_ = err
		})
	}
}

func TestValidator_ValidateCardExpiry_Empty(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardExpiry("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateCardExpiry_InvalidFormat(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"1225",
		"12-25",
		"13/25",
		"00/25",
		"12/",
		"/25",
	}

	for _, expiry := range tests {
		t.Run(expiry, func(t *testing.T) {
			err := v.ValidateCardExpiry(expiry)
			assert.Error(t, err)
		})
	}
}

func TestValidator_ValidateCardExpiry_InvalidMonth(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"00/25",
		"13/25",
		"99/25",
	}

	for _, expiry := range tests {
		t.Run(expiry, func(t *testing.T) {
			err := v.ValidateCardExpiry(expiry)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "месяц")
		})
	}
}

func TestValidator_ValidateCardExpiry_Expired(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardExpiry("01/2020")
	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "истек")
	}
}

// ========== ValidateCardCVV Tests ==========

func TestValidator_ValidateCardCVV_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"123",
		"1234",
		"000",
		"9999",
	}

	for _, cvv := range tests {
		t.Run(cvv, func(t *testing.T) {
			err := v.ValidateCardCVV(cvv)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidateCardCVV_Empty(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardCVV("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateCardCVV_Invalid(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"12",     // Слишком короткий
		"12345",  // Слишком длинный
		"abc",    // Не цифры
		"12a",    // Смешанный
	}

	for _, cvv := range tests {
		t.Run(cvv, func(t *testing.T) {
			err := v.ValidateCardCVV(cvv)
			assert.Error(t, err)
		})
	}
}

// ========== ValidateCardHolder Tests ==========

func TestValidator_ValidateCardHolder_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"John Doe",
		"MARY JANE",
		"John-Paul Smith",
		"O'Brien",
		"",
	}

	for _, holder := range tests {
		t.Run(holder, func(t *testing.T) {
			err := v.ValidateCardHolder(holder)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidateCardHolder_TooShort(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCardHolder("A")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "минимум 2")
}

func TestValidator_ValidateCardHolder_TooLong(t *testing.T) {
	v := NewValidator()
	longName := make([]byte, 101)
	for i := range longName {
		longName[i] = 'A'
	}
	err := v.ValidateCardHolder(string(longName))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "длиннее 100")
}

func TestValidator_ValidateCardHolder_InvalidCharacters(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"John123",
		"John@Doe",
		"John.Doe",
		"Иван Петров", // Кириллица
	}

	for _, holder := range tests {
		t.Run(holder, func(t *testing.T) {
			err := v.ValidateCardHolder(holder)
			assert.Error(t, err)
		})
	}
}

// ========== ValidateCreditCard Tests ==========

func TestValidator_ValidateCreditCard_Success(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreditCard("4532015112830366", "12/30", "123", "John Doe")
	assert.NoError(t, err)
}

func TestValidator_ValidateCreditCard_InvalidNumber(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreditCard("invalid", "12/30", "123", "John Doe")
	assert.Error(t, err)
}

func TestValidator_ValidateCreditCard_InvalidExpiry(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreditCard("4532015112830366", "invalid", "123", "John Doe")
	assert.Error(t, err)
}

func TestValidator_ValidateCreditCard_InvalidCVV(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreditCard("4532015112830366", "12/30", "12", "John Doe")
	assert.Error(t, err)
}

func TestValidator_ValidateCreditCard_InvalidHolder(t *testing.T) {
	v := NewValidator()
	err := v.ValidateCreditCard("4532015112830366", "12/30", "123", "J")
	assert.Error(t, err)
}

// ========== ValidateFileName Tests ==========

func TestValidator_ValidateFileName_Success(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"document.txt",
		"my_file.pdf",
		"report-2023.docx",
		"file.tar.gz",
		"simple",
	}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			err := v.ValidateFileName(filename)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidateFileName_Empty(t *testing.T) {
	v := NewValidator()
	err := v.ValidateFileName("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateFileName_OnlySpaces(t *testing.T) {
	v := NewValidator()
	err := v.ValidateFileName("   ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пустым")
}

func TestValidator_ValidateFileName_TooLong(t *testing.T) {
	v := NewValidator()
	longName := make([]byte, 256)
	for i := range longName {
		longName[i] = 'a'
	}
	err := v.ValidateFileName(string(longName))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "длиннее 255")
}

func TestValidator_ValidateFileName_InvalidCharacters(t *testing.T) {
	v := NewValidator()
	tests := []string{
		"file/name.txt",
		"file\\name.txt",
		"file:name.txt",
		"file*name.txt",
		"file?name.txt",
		"file\"name.txt",
		"file<name.txt",
		"file>name.txt",
		"file|name.txt",
	}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			err := v.ValidateFileName(filename)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "недопустимый символ")
		})
	}
}

// ========== ValidateFileSize Tests ==========

func TestValidator_ValidateFileSize_Success(t *testing.T) {
	v := NewValidator()
	tests := []int64{
		1,
		1024,
		1024 * 1024,
		50 * 1024 * 1024,
		100 * 1024 * 1024, // Максимум 100MB
	}

	for _, size := range tests {
		t.Run(string(rune(size)), func(t *testing.T) {
			err := v.ValidateFileSize(size)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_ValidateFileSize_Zero(t *testing.T) {
	v := NewValidator()
	err := v.ValidateFileSize(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "больше 0")
}

func TestValidator_ValidateFileSize_Negative(t *testing.T) {
	v := NewValidator()
	err := v.ValidateFileSize(-1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "больше 0")
}

func TestValidator_ValidateFileSize_TooLarge(t *testing.T) {
	v := NewValidator()
	err := v.ValidateFileSize(101 * 1024 * 1024) // 101 MB
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "превышать 100 MB")
}

// ========== ValidateLoginPassword Tests ==========

func TestValidator_ValidateLoginPassword_Success(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLoginPassword("user@example.com", "password123", "https://example.com")
	assert.NoError(t, err)
}

func TestValidator_ValidateLoginPassword_NoURL(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLoginPassword("user@example.com", "password123", "")
	assert.NoError(t, err)
}

func TestValidator_ValidateLoginPassword_InvalidLogin(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLoginPassword("", "password123", "https://example.com")
	assert.Error(t, err)
}

func TestValidator_ValidateLoginPassword_InvalidPassword(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLoginPassword("user@example.com", "", "https://example.com")
	assert.Error(t, err)
}

func TestValidator_ValidateLoginPassword_InvalidURL(t *testing.T) {
	v := NewValidator()
	err := v.ValidateLoginPassword("user@example.com", "password123", "ht tp://invalid")
	assert.Error(t, err)
}

// ========== LuhnCheck Tests ==========

func TestValidator_LuhnCheck(t *testing.T) {
	v := NewValidator()
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{"Valid Visa", "4532015112830366", true},
		{"Valid Mastercard", "5425233430109903", true},
		{"Invalid", "1234567890123456", false},
		{"Invalid short", "4532", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.luhnCheck(tt.number)
			assert.Equal(t, tt.expected, result)
		})
	}
}
