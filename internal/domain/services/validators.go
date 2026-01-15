package services

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Validator содержит методы валидации для различных типов данных
type Validator struct{}

// NewValidator создает новый валидатор
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateLoginPassword валидирует данные логин/пароль
func (v *Validator) ValidateLoginPassword(login, password, urlStr string) error {
	if err := v.ValidateLogin(login); err != nil {
		return err
	}
	
	if err := v.ValidatePassword(password); err != nil {
		return err
	}
	
	if urlStr != "" {
		if err := v.ValidateURL(urlStr); err != nil {
			return err
		}
	}
	
	return nil
}

// ValidateLogin валидирует логин
func (v *Validator) ValidateLogin(login string) error {
	if strings.TrimSpace(login) == "" {
		return fmt.Errorf("логин не может быть пустым")
	}
	
	if len(login) < 2 {
		return fmt.Errorf("логин должен содержать минимум 2 символа")
	}
	
	if len(login) > 100 {
		return fmt.Errorf("логин не может быть длиннее 100 символов")
	}
	
	return nil
}

// ValidatePassword валидирует пароль
func (v *Validator) ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("пароль не может быть пустым")
	}
	
	if len(password) < 1 {
		return fmt.Errorf("пароль должен содержать минимум 1 символ")
	}
	
	if len(password) > 500 {
		return fmt.Errorf("пароль не может быть длиннее 500 символов")
	}
	
	return nil
}

// ValidateURL валидирует URL
func (v *Validator) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return nil // URL необязательный
	}
	
	// Попробуем добавить схему, если её нет
	if !strings.Contains(urlStr, "://") {
		urlStr = "https://" + urlStr
	}
	
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("некорректный URL: %w", err)
	}
	
	if parsedURL.Host == "" {
		return fmt.Errorf("URL должен содержать домен")
	}
	
	return nil
}

// ValidateText валидирует текстовые данные
func (v *Validator) ValidateText(text string) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("текст не может быть пустым")
	}
	
	if len(text) > 10000 {
		return fmt.Errorf("текст не может быть длиннее 10000 символов")
	}
	
	return nil
}

// ValidateCreditCard валидирует данные кредитной карты
func (v *Validator) ValidateCreditCard(number, expiry, cvv, holder string) error {
	if err := v.ValidateCardNumber(number); err != nil {
		return err
	}
	
	if err := v.ValidateCardExpiry(expiry); err != nil {
		return err
	}
	
	if err := v.ValidateCardCVV(cvv); err != nil {
		return err
	}
	
	if err := v.ValidateCardHolder(holder); err != nil {
		return err
	}
	
	return nil
}

// ValidateCardNumber валидирует номер карты
func (v *Validator) ValidateCardNumber(number string) error {
	// Убираем пробелы и дефисы
	cleanNumber := strings.ReplaceAll(strings.ReplaceAll(number, " ", ""), "-", "")
	
	if cleanNumber == "" {
		return fmt.Errorf("номер карты не может быть пустым")
	}
	
	// Проверяем, что содержит только цифры
	if !regexp.MustCompile(`^\d+$`).MatchString(cleanNumber) {
		return fmt.Errorf("номер карты должен содержать только цифры")
	}
	
	// Проверяем длину (13-19 цифр для большинства карт)
	if len(cleanNumber) < 13 || len(cleanNumber) > 19 {
		return fmt.Errorf("номер карты должен содержать от 13 до 19 цифр")
	}
	
	// Алгоритм Луна для проверки контрольной суммы
	if !v.luhnCheck(cleanNumber) {
		return fmt.Errorf("некорректный номер карты (не прошел проверку Луна)")
	}
	
	return nil
}

// ValidateCardExpiry валидирует срок действия карты
func (v *Validator) ValidateCardExpiry(expiry string) error {
	if expiry == "" {
		return fmt.Errorf("срок действия карты не может быть пустым")
	}
	
	// Проверяем формат MM/YY или MM/YYYY
	parts := strings.Split(expiry, "/")
	if len(parts) != 2 {
		return fmt.Errorf("срок действия должен быть в формате MM/YY или MM/YYYY")
	}
	
	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return fmt.Errorf("некорректный месяц (должен быть от 01 до 12)")
	}
	
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("некорректный год")
	}
	
	// Если год двузначный, преобразуем в четырехзначный
	if year < 100 {
		currentYear := time.Now().Year()
		century := (currentYear / 100) * 100
		year += century
		
		// Если год получился в прошлом, добавляем столетие
		if year < currentYear {
			year += 100
		}
	}
	
	// Проверяем, что карта не истекла
	expiryDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	expiryDate = expiryDate.AddDate(0, 1, -1) // Последний день месяца
	
	if expiryDate.Before(time.Now()) {
		return fmt.Errorf("срок действия карты истек")
	}
	
	return nil
}

// ValidateCardCVV валидирует CVV код
func (v *Validator) ValidateCardCVV(cvv string) error {
	if cvv == "" {
		return fmt.Errorf("CVV код не может быть пустым")
	}
	
	if !regexp.MustCompile(`^\d{3,4}$`).MatchString(cvv) {
		return fmt.Errorf("CVV код должен содержать 3 или 4 цифры")
	}
	
	return nil
}

// ValidateCardHolder валидирует имя держателя карты
func (v *Validator) ValidateCardHolder(holder string) error {
	if holder == "" {
		return nil // Имя держателя необязательно
	}
	
	holder = strings.TrimSpace(holder)
	if len(holder) < 2 {
		return fmt.Errorf("имя держателя карты должно содержать минимум 2 символа")
	}
	
	if len(holder) > 100 {
		return fmt.Errorf("имя держателя карты не может быть длиннее 100 символов")
	}
	
	// Проверяем, что содержит только буквы, пробелы, дефисы и апострофы
	if !regexp.MustCompile(`^[a-zA-Z\s\-']+$`).MatchString(holder) {
		return fmt.Errorf("имя держателя карты должно содержать только буквы, пробелы, дефисы и апострофы")
	}
	
	return nil
}

// ValidateFileName валидирует имя файла
func (v *Validator) ValidateFileName(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("имя файла не может быть пустым")
	}
	
	if len(filename) > 255 {
		return fmt.Errorf("имя файла не может быть длиннее 255 символов")
	}
	
	// Проверяем на недопустимые символы для имени файла
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("имя файла содержит недопустимый символ: %s", char)
		}
	}
	
	return nil
}

// ValidateFileSize валидирует размер файла
func (v *Validator) ValidateFileSize(size int64) error {
	const maxSize = 100 * 1024 * 1024 // 100 MB
	
	if size <= 0 {
		return fmt.Errorf("размер файла должен быть больше 0")
	}
	
	if size > maxSize {
		return fmt.Errorf("размер файла не может превышать 100 MB")
	}
	
	return nil
}

// luhnCheck проверяет номер карты по алгоритму Луна
func (v *Validator) luhnCheck(number string) bool {
	var sum int
	var alternate bool
	
	// Проходим справа налево
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		
		sum += digit
		alternate = !alternate
	}
	
	return sum%10 == 0
}