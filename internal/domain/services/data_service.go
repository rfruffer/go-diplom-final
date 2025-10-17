package services

import (
	"context"
	"fmt"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/fylgushev/go-diplom-final/internal/domain/interfaces"
	"github.com/fylgushev/go-diplom-final/internal/infrastructure/crypto"
)

// DataService сервис для работы с данными пользователя
type DataService struct {
	dataRepo   interfaces.DataRepository
	serializer *DataSerializer
	cipher     crypto.Cipher
}

// NewDataService создает новый сервис данных
func NewDataService(dataRepo interfaces.DataRepository, masterKey []byte) (*DataService, error) {
	cipher, err := crypto.NewAESCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	return &DataService{
		dataRepo:   dataRepo,
		serializer: NewDataSerializer(),
		cipher:     cipher,
	}, nil
}

// StoreLoginPassword сохраняет данные логин/пароль
func (ds *DataService) StoreLoginPassword(ctx context.Context, userID, name string, data *entities.LoginPasswordData, metadata map[string]interface{}) (*entities.DataItem, error) {
	// Сериализуем данные
	serializedData, err := ds.serializer.SerializeLoginPassword(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize login/password data: %w", err)
	}

	// Шифруем данные
	encryptedData, err := ds.cipher.Encrypt(serializedData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Создаем элемент данных
	item := entities.NewDataItem(userID, entities.DataTypeLoginPassword, name, encryptedData)
	
	// Устанавливаем метаданные
	if metadata != nil {
		if err := item.SetMetadata(metadata); err != nil {
			return nil, fmt.Errorf("failed to set metadata: %w", err)
		}
	}

	// Сохраняем в репозитории
	if err := ds.dataRepo.CreateDataItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to store data item: %w", err)
	}

	return item, nil
}

// GetLoginPassword получает данные логин/пароль
func (ds *DataService) GetLoginPassword(ctx context.Context, userID, itemID string) (*entities.LoginPasswordData, *entities.DataItem, error) {
	// Получаем элемент данных
	item, err := ds.dataRepo.GetDataItem(ctx, userID, itemID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get data item: %w", err)
	}

	// Проверяем тип данных
	if item.Type != entities.DataTypeLoginPassword {
		return nil, nil, fmt.Errorf("item is not login/password type")
	}

	// Расшифровываем данные
	decryptedData, err := ds.cipher.Decrypt(item.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Десериализуем данные
	loginPassword, err := ds.serializer.DeserializeLoginPassword(decryptedData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize login/password data: %w", err)
	}

	return loginPassword, item, nil
}

// StoreTextData сохраняет текстовые данные
func (ds *DataService) StoreTextData(ctx context.Context, userID, name string, data *entities.TextData, metadata map[string]interface{}) (*entities.DataItem, error) {
	serializedData, err := ds.serializer.SerializeTextData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize text data: %w", err)
	}

	encryptedData, err := ds.cipher.Encrypt(serializedData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	item := entities.NewDataItem(userID, entities.DataTypeText, name, encryptedData)
	
	if metadata != nil {
		if err := item.SetMetadata(metadata); err != nil {
			return nil, fmt.Errorf("failed to set metadata: %w", err)
		}
	}

	if err := ds.dataRepo.CreateDataItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to store data item: %w", err)
	}

	return item, nil
}

// GetTextData получает текстовые данные
func (ds *DataService) GetTextData(ctx context.Context, userID, itemID string) (*entities.TextData, *entities.DataItem, error) {
	item, err := ds.dataRepo.GetDataItem(ctx, userID, itemID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get data item: %w", err)
	}

	if item.Type != entities.DataTypeText {
		return nil, nil, fmt.Errorf("item is not text type")
	}

	decryptedData, err := ds.cipher.Decrypt(item.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	textData, err := ds.serializer.DeserializeTextData(decryptedData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize text data: %w", err)
	}

	return textData, item, nil
}

// StoreCreditCard сохраняет данные кредитной карты
func (ds *DataService) StoreCreditCard(ctx context.Context, userID, name string, data *entities.CreditCardData, metadata map[string]interface{}) (*entities.DataItem, error) {
	serializedData, err := ds.serializer.SerializeCreditCard(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize credit card data: %w", err)
	}

	encryptedData, err := ds.cipher.Encrypt(serializedData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	item := entities.NewDataItem(userID, entities.DataTypeCreditCard, name, encryptedData)
	
	if metadata != nil {
		if err := item.SetMetadata(metadata); err != nil {
			return nil, fmt.Errorf("failed to set metadata: %w", err)
		}
	}

	if err := ds.dataRepo.CreateDataItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to store data item: %w", err)
	}

	return item, nil
}

// GetCreditCard получает данные кредитной карты
func (ds *DataService) GetCreditCard(ctx context.Context, userID, itemID string) (*entities.CreditCardData, *entities.DataItem, error) {
	item, err := ds.dataRepo.GetDataItem(ctx, userID, itemID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get data item: %w", err)
	}

	if item.Type != entities.DataTypeCreditCard {
		return nil, nil, fmt.Errorf("item is not credit card type")
	}

	decryptedData, err := ds.cipher.Decrypt(item.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	cardData, err := ds.serializer.DeserializeCreditCard(decryptedData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize credit card data: %w", err)
	}

	return cardData, item, nil
}

// StoreBinaryData сохраняет бинарные данные
func (ds *DataService) StoreBinaryData(ctx context.Context, userID, name string, data *entities.BinaryData, metadata map[string]interface{}) (*entities.DataItem, error) {
	serializedData, err := ds.serializer.SerializeBinaryData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize binary data: %w", err)
	}

	encryptedData, err := ds.cipher.Encrypt(serializedData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	item := entities.NewDataItem(userID, entities.DataTypeBinary, name, encryptedData)
	
	if metadata != nil {
		if err := item.SetMetadata(metadata); err != nil {
			return nil, fmt.Errorf("failed to set metadata: %w", err)
		}
	}

	if err := ds.dataRepo.CreateDataItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to store data item: %w", err)
	}

	return item, nil
}

// GetBinaryData получает бинарные данные
func (ds *DataService) GetBinaryData(ctx context.Context, userID, itemID string) (*entities.BinaryData, *entities.DataItem, error) {
	item, err := ds.dataRepo.GetDataItem(ctx, userID, itemID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get data item: %w", err)
	}

	if item.Type != entities.DataTypeBinary {
		return nil, nil, fmt.Errorf("item is not binary type")
	}

	decryptedData, err := ds.cipher.Decrypt(item.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	binaryData, err := ds.serializer.DeserializeBinaryData(decryptedData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize binary data: %w", err)
	}

	return binaryData, item, nil
}

// ListUserData получает список данных пользователя
func (ds *DataService) ListUserData(ctx context.Context, userID string, dataType entities.DataType, limit, offset int32) ([]*entities.DataItem, int32, error) {
	return ds.dataRepo.GetUserDataItems(ctx, userID, dataType, limit, offset)
}

// DeleteDataItem удаляет элемент данных
func (ds *DataService) DeleteDataItem(ctx context.Context, userID, itemID string) error {
	return ds.dataRepo.DeleteDataItem(ctx, userID, itemID)
}

// UpdateDataItem обновляет элемент данных
func (ds *DataService) UpdateDataItem(ctx context.Context, item *entities.DataItem) error {
	item.IncrementVersion()
	return ds.dataRepo.UpdateDataItem(ctx, item)
}