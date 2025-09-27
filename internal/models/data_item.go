package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DataItem представляет элемент данных, хранимый пользователем в системе GophKeeper.
// Содержит зашифрованные пользовательские данные любого из поддерживаемых типов
// вместе с метаинформацией и данными для синхронизации.
type DataItem struct {
	// ID уникальный идентификатор элемента данных.
	// Генерируется автоматически при создании элемента с использованием UUID.
	ID string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// UserID идентификатор пользователя-владельца данных.
	// Связывает элемент данных с конкретным пользователем.
	UserID string `json:"user_id" gorm:"type:uuid;not null;index"`

	// Type тип хранимых данных.
	// Определяет, как интерпретировать содержимое поля Data.
	Type DataType `json:"type" gorm:"type:int;not null"`

	// Name пользовательское название элемента данных.
	// Позволяет пользователю идентифицировать элемент (например, "GitHub аккаунт").
	Name string `json:"name" gorm:"type:varchar(255);not null"`

	// Data зашифрованные данные пользователя.
	// Содержимое всегда зашифровано перед сохранением и расшифровывается при получении.
	Data []byte `json:"data" gorm:"type:bytea;not null"`

	// Metadata произвольная текстовая метаинформация.
	// Хранится в формате JSON и может содержать любые дополнительные данные
	// (URL сайта, описание, теги и т.д.).
	Metadata json.RawMessage `json:"metadata" gorm:"type:jsonb"`

	// CreatedAt время создания элемента данных.
	// Автоматически устанавливается при создании элемента.
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt время последнего обновления элемента данных.
	// Автоматически обновляется при любых изменениях элемента.
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Version версия элемента данных для управления конфликтами.
	// Инкрементируется при каждом обновлении и используется для разрешения
	// конфликтов при синхронизации между клиентами.
	Version int64 `json:"version" gorm:"default:1"`

	// IsDeleted флаг мягкого удаления.
	// Позволяет помечать элементы как удаленные без физического удаления
	// для корректной синхронизации между клиентами.
	IsDeleted bool `json:"is_deleted" gorm:"default:false"`

	// User связь с моделью пользователя.
	// Используется ORM для автоматической загрузки связанных данных.
	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// NewDataItem создает новый элемент данных с базовыми значениями.
// Генерирует уникальный ID и устанавливает версию 1.
// Время создания и обновления будут установлены автоматически ORM.
func NewDataItem(userID string, dataType DataType, name string, data []byte) *DataItem {
	return &DataItem{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      dataType,
		Name:      name,
		Data:      data,
		Version:   1,
		IsDeleted: false,
	}
}

// SetMetadata устанавливает метаинформацию для элемента данных.
// Принимает любую структуру, которая может быть сериализована в JSON.
// Возвращает ошибку, если сериализация не удалась.
func (di *DataItem) SetMetadata(metadata interface{}) error {
	if metadata == nil {
		di.Metadata = nil
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	di.Metadata = data
	return nil
}

// GetMetadata извлекает метаинформацию элемента данных.
// Десериализует JSON метаданные в предоставленную структуру.
// Возвращает ошибку, если десериализация не удалась.
func (di *DataItem) GetMetadata(target interface{}) error {
	if di.Metadata == nil {
		return nil
	}

	return json.Unmarshal(di.Metadata, target)
}

// IncrementVersion увеличивает версию элемента данных на 1.
// Вызывается при каждом обновлении элемента для отслеживания изменений.
func (di *DataItem) IncrementVersion() {
	di.Version++
}

// MarkAsDeleted помечает элемент как удаленный без физического удаления.
// Используется для мягкого удаления, которое необходимо для корректной
// синхронизации удалений между клиентами.
func (di *DataItem) MarkAsDeleted() {
	di.IsDeleted = true
	di.IncrementVersion()
}

// Restore восстанавливает ранее удаленный элемент.
// Снимает флаг удаления и увеличивает версию.
func (di *DataItem) Restore() {
	di.IsDeleted = false
	di.IncrementVersion()
}

// TableName возвращает название таблицы в базе данных для модели DataItem.
// Используется ORM GORM для корректного маппинга модели на таблицу.
func (DataItem) TableName() string {
	return "data_items"
}