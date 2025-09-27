package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CreditCardData представляет структуру данных для типа CreditCard.
// Содержит информацию о банковской карте с валидацией основных полей.
// Эта структура сериализуется в JSON и шифруется перед сохранением в DataItem.Data.
type CreditCardData struct {
	// Number номер банковской карты.
	// Должен содержать от 13 до 19 цифр согласно стандартам платежных систем.
	Number string `json:"number"`

	// ExpiryMonth месяц истечения срока действия карты (1-12).
	// Используется совместно с ExpiryYear для определения срока действия.
	ExpiryMonth int `json:"expiry_month"`

	// ExpiryYear год истечения срока действия карты.
	// Должен быть больше или равен текущему году.
	ExpiryYear int `json:"expiry_year"`

	// CVV код безопасности карты.
	// Обычно 3 цифры на обратной стороне карты или 4 цифры на лицевой стороне (American Express).
	CVV string `json:"cvv"`

	// CardholderName имя держателя карты.
	// Имя, указанное на карте, обычно в формате "FIRSTNAME LASTNAME".
	CardholderName string `json:"cardholder_name"`

	// BankName название банка-эмитента карты.
	// Опциональное поле для удобства идентификации карты.
	BankName string `json:"bank_name,omitempty"`

	// Notes дополнительные заметки к карте.
	// Может содержать информацию о лимитах, PIN-коде и другую полезную информацию.
	Notes string `json:"notes,omitempty"`
}

// NewCreditCardData создает новый экземпляр данных банковской карты.
// Принимает обязательные параметры: номер карты, месяц и год истечения, CVV и имя держателя.
func NewCreditCardData(number string, expiryMonth, expiryYear int, cvv, cardholderName string) *CreditCardData {
	return &CreditCardData{
		Number:         number,
		ExpiryMonth:    expiryMonth,
		ExpiryYear:     expiryYear,
		CVV:            cvv,
		CardholderName: cardholderName,
	}
}

// SetBankName устанавливает название банка для карты.
// Используется для дополнительной идентификации карты пользователем.
func (ccd *CreditCardData) SetBankName(bankName string) {
	ccd.BankName = bankName
}

// SetNotes устанавливает заметки для карты.
// Может содержать любую дополнительную информацию о карте.
func (ccd *CreditCardData) SetNotes(notes string) {
	ccd.Notes = notes
}

// Validate проверяет корректность данных банковской карты.
// Проверяет формат номера карты, срок действия, CVV и имя держателя.
// Возвращает ошибку с описанием проблемы, если валидация не прошла.
func (ccd *CreditCardData) Validate() error {
	if err := ccd.validateNumber(); err != nil {
		return err
	}
	if err := ccd.validateExpiry(); err != nil {
		return err
	}
	if err := ccd.validateCVV(); err != nil {
		return err
	}
	if err := ccd.validateCardholderName(); err != nil {
		return err
	}
	return nil
}

// validateNumber проверяет корректность номера карты.
// Проверяет длину номера и соответствие алгоритму Луна.
func (ccd *CreditCardData) validateNumber() error {
	// Убираем все пробелы и дефисы
	number := strings.ReplaceAll(strings.ReplaceAll(ccd.Number, " ", ""), "-", "")
	
	// Проверяем, что номер состоит только из цифр
	if !regexp.MustCompile(`^\d+$`).MatchString(number) {
		return fmt.Errorf("номер карты должен содержать только цифры")
	}

	// Проверяем длину номера
	if len(number) < 13 || len(number) > 19 {
		return fmt.Errorf("номер карты должен содержать от 13 до 19 цифр")
	}

	// Проверяем алгоритм Луна
	if !ccd.luhnCheck(number) {
		return fmt.Errorf("некорректный номер карты (не прошел проверку алгоритма Луна)")
	}

	return nil
}

// validateExpiry проверяет корректность срока действия карты.
func (ccd *CreditCardData) validateExpiry() error {
	if ccd.ExpiryMonth < 1 || ccd.ExpiryMonth > 12 {
		return fmt.Errorf("месяц истечения должен быть от 1 до 12")
	}

	currentYear := time.Now().Year()
	if ccd.ExpiryYear < currentYear {
		return fmt.Errorf("год истечения не может быть в прошлом")
	}

	// Проверяем, что карта не истекла в текущем году
	if ccd.ExpiryYear == currentYear && ccd.ExpiryMonth < int(time.Now().Month()) {
		return fmt.Errorf("карта уже истекла")
	}

	return nil
}

// validateCVV проверяет корректность CVV кода.
func (ccd *CreditCardData) validateCVV() error {
	if !regexp.MustCompile(`^\d{3,4}$`).MatchString(ccd.CVV) {
		return fmt.Errorf("CVV должен содержать 3 или 4 цифры")
	}
	return nil
}

// validateCardholderName проверяет корректность имени держателя карты.
func (ccd *CreditCardData) validateCardholderName() error {
	if strings.TrimSpace(ccd.CardholderName) == "" {
		return fmt.Errorf("имя держателя карты не может быть пустым")
	}
	return nil
}

// luhnCheck проверяет номер карты по алгоритму Луна.
// Возвращает true, если номер корректен.
func (ccd *CreditCardData) luhnCheck(number string) bool {
	var sum int
	isEven := false

	// Проходим номер справа налево
	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if isEven {
			digit *= 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}

		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}

// GetMaskedNumber возвращает замаскированный номер карты.
// Показывает только последние 4 цифры, остальные заменяет звездочками.
// Используется для безопасного отображения номера карты.
func (ccd *CreditCardData) GetMaskedNumber() string {
	if len(ccd.Number) < 4 {
		return "****"
	}
	return "****-****-****-" + ccd.Number[len(ccd.Number)-4:]
}

// IsExpired проверяет, истек ли срок действия карты.
// Возвращает true, если карта просрочена.
func (ccd *CreditCardData) IsExpired() bool {
	now := time.Now()
	return ccd.ExpiryYear < now.Year() || 
		   (ccd.ExpiryYear == now.Year() && ccd.ExpiryMonth < int(now.Month()))
}