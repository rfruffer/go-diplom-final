package data

import (
	"context"
	"time"
)

// DataType represents the type of stored data
type DataType int32

const (
	DataTypeUnspecified DataType = 0
	DataTypeCredentials DataType = 1
	DataTypeText        DataType = 2
	DataTypeBinary      DataType = 3
	DataTypeCard        DataType = 4
)

// DataServiceServer is the interface for data service implementation
type DataServiceServer interface {
	StoreData(ctx context.Context, req *StoreDataRequest) (*StoreDataResponse, error)
	GetData(ctx context.Context, req *GetDataRequest) (*GetDataResponse, error)
	GetUserData(ctx context.Context, req *GetUserDataRequest) (*GetUserDataResponse, error)
	UpdateData(ctx context.Context, req *UpdateDataRequest) (*UpdateDataResponse, error)
	DeleteData(ctx context.Context, req *DeleteDataRequest) (*DeleteDataResponse, error)
	SyncData(ctx context.Context, req *SyncDataRequest) (*SyncDataResponse, error)
}

// DataServiceClient is the interface for data service client
type DataServiceClient interface {
	StoreData(ctx context.Context, req *StoreDataRequest) (*StoreDataResponse, error)
	GetData(ctx context.Context, req *GetDataRequest) (*GetDataResponse, error)
	GetUserData(ctx context.Context, req *GetUserDataRequest) (*GetUserDataResponse, error)
	UpdateData(ctx context.Context, req *UpdateDataRequest) (*UpdateDataResponse, error)
	DeleteData(ctx context.Context, req *DeleteDataRequest) (*DeleteDataResponse, error)
	SyncData(ctx context.Context, req *SyncDataRequest) (*SyncDataResponse, error)
}

// StoreDataRequest represents request to store new data
type StoreDataRequest struct {
	UserID   string    `json:"user_id"`
	Type     DataType  `json:"type"`
	Metadata string    `json:"metadata"`
	Data     []byte    `json:"data"`
	Checksum string    `json:"checksum"`
}

// StoreDataResponse represents response after storing data
type StoreDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	DataID  string `json:"data_id"`
	Version int64  `json:"version"`
}

// GetDataRequest represents request to get data by ID
type GetDataRequest struct {
	UserID string `json:"user_id"`
	DataID string `json:"data_id"`
}

// GetDataResponse represents response with requested data
type GetDataResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	DataItem *DataItem `json:"data_item"`
}

// GetUserDataRequest represents request to get all user data
type GetUserDataRequest struct {
	UserID string   `json:"user_id"`
	Type   DataType `json:"type"`
	Limit  int32    `json:"limit"`
	Offset int32    `json:"offset"`
}

// GetUserDataResponse represents response with user data list
type GetUserDataResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	DataItems  []*DataItem `json:"data_items"`
	TotalCount int32       `json:"total_count"`
}

// UpdateDataRequest represents request to update existing data
type UpdateDataRequest struct {
	UserID         string   `json:"user_id"`
	DataID         string   `json:"data_id"`
	Metadata       string   `json:"metadata"`
	Data           []byte   `json:"data"`
	Checksum       string   `json:"checksum"`
	CurrentVersion int64    `json:"current_version"`
}

// UpdateDataResponse represents response after updating data
type UpdateDataResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	NewVersion int64  `json:"new_version"`
}

// DeleteDataRequest represents request to delete data
type DeleteDataRequest struct {
	UserID         string `json:"user_id"`
	DataID         string `json:"data_id"`
	CurrentVersion int64  `json:"current_version"`
}

// DeleteDataResponse represents response after deleting data
type DeleteDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SyncDataRequest represents request to sync data
type SyncDataRequest struct {
	UserID     string    `json:"user_id"`
	LastSync   time.Time `json:"last_sync"`
	ExcludeIDs []string  `json:"exclude_ids"`
}

// SyncDataResponse represents response with sync data
type SyncDataResponse struct {
	Success       bool        `json:"success"`
	Message       string      `json:"message"`
	UpdatedItems  []*DataItem `json:"updated_items"`
	DeletedIDs    []string    `json:"deleted_ids"`
	SyncTimestamp time.Time   `json:"sync_timestamp"`
}

// DataItem represents a stored data item
type DataItem struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      DataType  `json:"type"`
	Metadata  string    `json:"metadata"`
	Data      []byte    `json:"data"`
	Checksum  string    `json:"checksum"`
	Version   int64     `json:"version"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Metadata structures for different data types

// CredentialsMetadata represents metadata for credentials
type CredentialsMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// TextMetadata represents metadata for text notes
type TextMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// BinaryMetadata represents metadata for binary files
type BinaryMetadata struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
	Description string `json:"description"`
}

// CardMetadata represents metadata for bank cards
type CardMetadata struct {
	Name        string `json:"name"`
	BankName    string `json:"bank_name"`
	CardType    string `json:"card_type"`
	Description string `json:"description"`
}