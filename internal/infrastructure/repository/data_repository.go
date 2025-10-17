package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/fylgushev/go-diplom-final/internal/domain/interfaces"
)

// DataRepository реализует интерфейс DataRepository с использованием GORM
type DataRepository struct {
	db *gorm.DB
}

// NewDataRepository создает новый экземпляр DataRepository
func NewDataRepository(db *gorm.DB) interfaces.DataRepository {
	return &DataRepository{db: db}
}

// CreateDataItem создает новый элемент данных в базе данных
func (r *DataRepository) CreateDataItem(ctx context.Context, item *entities.DataItem) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("failed to create data item: %w", err)
	}
	return nil
}

// GetDataItem получает элемент данных по пользователю и ID
func (r *DataRepository) GetDataItem(ctx context.Context, userID, itemID string) (*entities.DataItem, error) {
	var item entities.DataItem
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, itemID).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrDataItemNotFound
		}
		return nil, fmt.Errorf("failed to get data item: %w", err)
	}
	return &item, nil
}

// GetUserDataItems получает элементы данных пользователя с пагинацией
func (r *DataRepository) GetUserDataItems(ctx context.Context, userID string, dataType entities.DataType, limit, offset int32) ([]*entities.DataItem, int32, error) {
	var items []*entities.DataItem
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	// Фильтр по типу данных, если указан
	if dataType != entities.DataTypeUnknown {
		query = query.Where("type = ?", dataType)
	}

	// Подсчет общего количества
	if err := query.Model(&entities.DataItem{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count data items: %w", err)
	}

	// Получение данных с пагинацией
	err := query.Order("updated_at DESC").Limit(int(limit)).Offset(int(offset)).Find(&items).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user data items: %w", err)
	}

	return items, int32(total), nil
}

// UpdateDataItem обновляет элемент данных
func (r *DataRepository) UpdateDataItem(ctx context.Context, item *entities.DataItem) error {
	err := r.db.WithContext(ctx).Save(item).Error
	if err != nil {
		return fmt.Errorf("failed to update data item: %w", err)
	}
	return nil
}

// DeleteDataItem удаляет элемент данных
func (r *DataRepository) DeleteDataItem(ctx context.Context, userID, itemID string) error {
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, itemID).Delete(&entities.DataItem{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete data item: %w", err)
	}
	return nil
}

// GetDataItemsSince получает элементы данных, обновленные после указанного времени
func (r *DataRepository) GetDataItemsSince(ctx context.Context, userID string, since int64) ([]*entities.DataItem, error) {
	var items []*entities.DataItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND updated_at > ?", userID, since).
		Order("updated_at ASC").
		Find(&items).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get data items since: %w", err)
	}
	
	return items, nil
}

// GetUserDataItemsModifiedAfter получает элементы данных пользователя, измененные после указанного времени
func (r *DataRepository) GetUserDataItemsModifiedAfter(ctx context.Context, userID string, after time.Time) ([]*entities.DataItem, int32, error) {
	var items []*entities.DataItem
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ? AND updated_at > ?", userID, after)
	
	// Подсчет общего количества
	if err := query.Model(&entities.DataItem{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count modified data items: %w", err)
	}

	// Получение данных
	err := query.Order("updated_at ASC").Find(&items).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get modified data items: %w", err)
	}

	return items, int32(total), nil
}

// GetUserLastSyncTime получает время последней синхронизации пользователя
func (r *DataRepository) GetUserLastSyncTime(ctx context.Context, userID string) (time.Time, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Select("last_sync_at").Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return time.Time{}, entities.ErrUserNotFound
		}
		return time.Time{}, fmt.Errorf("failed to get user last sync time: %w", err)
	}
	
	if user.LastSyncAt == nil {
		return time.Time{}, nil
	}
	
	return *user.LastSyncAt, nil
}

// UpdateUserLastSyncTime обновляет время последней синхронизации пользователя
func (r *DataRepository) UpdateUserLastSyncTime(ctx context.Context, userID string, syncTime time.Time) error {
	err := r.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Update("last_sync_at", syncTime).Error
	
	if err != nil {
		return fmt.Errorf("failed to update user last sync time: %w", err)
	}
	
	return nil
}