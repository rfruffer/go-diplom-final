package interfaces

import (
	"context"
	"time"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
)

// UserRepository defines the interface for user storage operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByLogin(ctx context.Context, login string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, userID string) error
}

// DataRepository defines the interface for data storage operations
type DataRepository interface {
	CreateDataItem(ctx context.Context, item *entities.DataItem) error
	GetDataItem(ctx context.Context, userID, itemID string) (*entities.DataItem, error)
	GetUserDataItems(ctx context.Context, userID string, dataType entities.DataType, limit, offset int32) ([]*entities.DataItem, int32, error)
	UpdateDataItem(ctx context.Context, item *entities.DataItem) error
	DeleteDataItem(ctx context.Context, userID, itemID string) error
	GetDataItemsSince(ctx context.Context, userID string, since int64) ([]*entities.DataItem, error)
	
	// Sync-related methods
	GetUserDataItemsModifiedAfter(ctx context.Context, userID string, after time.Time) ([]*entities.DataItem, int32, error)
	GetUserLastSyncTime(ctx context.Context, userID string) (time.Time, error)
	UpdateUserLastSyncTime(ctx context.Context, userID string, syncTime time.Time) error
}