// Package models содержит основные модели данных для системы GophKeeper.
// Этот пакет определяет структуры для пользователей, элементов данных,
// типов данных и связанных с ними операций согласно техническому заданию.
package entities

import "fmt"

// DataType представляет тип хранимых данных в системе GophKeeper.
// Поддерживаются четыре основных типа данных согласно техническому заданию:
// пары логин/пароль, произвольные текстовые данные, бинарные данные и данные банковских карт.
type DataType int

const (
	// DataTypeLoginPassword представляет тип данных для пар логин/пароль.
	// Используется для хранения учетных записей пользователей на различных сервисах.
	DataTypeLoginPassword DataType = iota

	// DataTypeText представляет тип для произвольных текстовых данных.
	// Может содержать заметки, коды активации, секретные вопросы и другую текстовую информацию.
	DataTypeText

	// DataTypeBinary представляет тип для произвольных бинарных данных.
	// Используется для хранения файлов, изображений, документов и других бинарных объектов.
	DataTypeBinary

	// DataTypeCreditCard представляет тип для данных банковских карт.
	// Содержит номер карты, срок действия, CVV и другую связанную информацию.
	DataTypeCreditCard
)

// String возвращает строковое представление типа данных.
// Используется для логирования, отображения пользователю и сериализации.
// Возвращает понятное человеку название типа данных.
func (dt DataType) String() string {
	switch dt {
	case DataTypeLoginPassword:
		return "login_password"
	case DataTypeText:
		return "text_data"
	case DataTypeBinary:
		return "binary_data"
	case DataTypeCreditCard:
		return "credit_card"
	default:
		return fmt.Sprintf("unknown_data_type_%d", int(dt))
	}
}

// IsValid проверяет, является ли тип данных валидным.
// Возвращает true, если тип данных определен в системе, иначе false.
// Используется для валидации входных данных.
func (dt DataType) IsValid() bool {
	return dt >= DataTypeLoginPassword && dt <= DataTypeCreditCard
}

// ParseDataType преобразует строковое представление в тип DataType.
// Принимает строку с названием типа данных и возвращает соответствующий DataType.
// Если строка не соответствует ни одному типу, возвращает ошибку.
func ParseDataType(s string) (DataType, error) {
	switch s {
	case "login_password":
		return DataTypeLoginPassword, nil
	case "text_data":
		return DataTypeText, nil
	case "binary_data":
		return DataTypeBinary, nil
	case "credit_card":
		return DataTypeCreditCard, nil
	default:
		return 0, fmt.Errorf("неизвестный тип данных: %s", s)
	}
}

// AllDataTypes возвращает срез всех поддерживаемых типов данных.
// Используется для итерации по всем возможным типам и валидации.
func AllDataTypes() []DataType {
	return []DataType{DataTypeLoginPassword, DataTypeText, DataTypeBinary, DataTypeCreditCard}
}