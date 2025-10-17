package services

import (
	"encoding/json"
	"fmt"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
)

// DataSerializer сервис для сериализации и десериализации данных
type DataSerializer struct{}

// NewDataSerializer создает новый сервис сериализации данных
func NewDataSerializer() *DataSerializer {
	return &DataSerializer{}
}

// SerializeLoginPassword сериализует данные логин/пароль
func (ds *DataSerializer) SerializeLoginPassword(data *entities.LoginPasswordData) ([]byte, error) {
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize login/password data: %w", err)
	}

	return jsonData, nil
}

// DeserializeLoginPassword десериализует данные логин/пароль
func (ds *DataSerializer) DeserializeLoginPassword(data []byte) (*entities.LoginPasswordData, error) {
	var loginPassword entities.LoginPasswordData
	if err := json.Unmarshal(data, &loginPassword); err != nil {
		return nil, fmt.Errorf("failed to deserialize login/password data: %w", err)
	}

	if err := loginPassword.Validate(); err != nil {
		return nil, fmt.Errorf("deserialized data validation failed: %w", err)
	}

	return &loginPassword, nil
}

// SerializeTextData сериализует текстовые данные
func (ds *DataSerializer) SerializeTextData(data *entities.TextData) ([]byte, error) {
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize text data: %w", err)
	}

	return jsonData, nil
}

// DeserializeTextData десериализует текстовые данные
func (ds *DataSerializer) DeserializeTextData(data []byte) (*entities.TextData, error) {
	var textData entities.TextData
	if err := json.Unmarshal(data, &textData); err != nil {
		return nil, fmt.Errorf("failed to deserialize text data: %w", err)
	}

	if err := textData.Validate(); err != nil {
		return nil, fmt.Errorf("deserialized data validation failed: %w", err)
	}

	return &textData, nil
}

// SerializeCreditCard сериализует данные кредитной карты
func (ds *DataSerializer) SerializeCreditCard(data *entities.CreditCardData) ([]byte, error) {
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize credit card data: %w", err)
	}

	return jsonData, nil
}

// DeserializeCreditCard десериализует данные кредитной карты
func (ds *DataSerializer) DeserializeCreditCard(data []byte) (*entities.CreditCardData, error) {
	var cardData entities.CreditCardData
	if err := json.Unmarshal(data, &cardData); err != nil {
		return nil, fmt.Errorf("failed to deserialize credit card data: %w", err)
	}

	if err := cardData.Validate(); err != nil {
		return nil, fmt.Errorf("deserialized data validation failed: %w", err)
	}

	return &cardData, nil
}

// SerializeBinaryData сериализует бинарные данные
func (ds *DataSerializer) SerializeBinaryData(data *entities.BinaryData) ([]byte, error) {
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize binary data: %w", err)
	}

	return jsonData, nil
}

// DeserializeBinaryData десериализует бинарные данные
func (ds *DataSerializer) DeserializeBinaryData(data []byte) (*entities.BinaryData, error) {
	var binaryData entities.BinaryData
	if err := json.Unmarshal(data, &binaryData); err != nil {
		return nil, fmt.Errorf("failed to deserialize binary data: %w", err)
	}

	if err := binaryData.Validate(); err != nil {
		return nil, fmt.Errorf("deserialized data validation failed: %w", err)
	}

	return &binaryData, nil
}

// SerializeByType сериализует данные в зависимости от типа
func (ds *DataSerializer) SerializeByType(dataType entities.DataType, data interface{}) ([]byte, error) {
	switch dataType {
	case entities.DataTypeLoginPassword:
		loginPassword, ok := data.(*entities.LoginPasswordData)
		if !ok {
			return nil, fmt.Errorf("invalid data type for login/password")
		}
		return ds.SerializeLoginPassword(loginPassword)

	case entities.DataTypeText:
		textData, ok := data.(*entities.TextData)
		if !ok {
			return nil, fmt.Errorf("invalid data type for text")
		}
		return ds.SerializeTextData(textData)

	case entities.DataTypeCreditCard:
		cardData, ok := data.(*entities.CreditCardData)
		if !ok {
			return nil, fmt.Errorf("invalid data type for credit card")
		}
		return ds.SerializeCreditCard(cardData)

	case entities.DataTypeBinary:
		binaryData, ok := data.(*entities.BinaryData)
		if !ok {
			return nil, fmt.Errorf("invalid data type for binary")
		}
		return ds.SerializeBinaryData(binaryData)

	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType.String())
	}
}

// DeserializeByType десериализует данные в зависимости от типа
func (ds *DataSerializer) DeserializeByType(dataType entities.DataType, data []byte) (interface{}, error) {
	switch dataType {
	case entities.DataTypeLoginPassword:
		return ds.DeserializeLoginPassword(data)

	case entities.DataTypeText:
		return ds.DeserializeTextData(data)

	case entities.DataTypeCreditCard:
		return ds.DeserializeCreditCard(data)

	case entities.DataTypeBinary:
		return ds.DeserializeBinaryData(data)

	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType.String())
	}
}