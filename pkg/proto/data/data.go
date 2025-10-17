package data

import (
	"context"
	"time"
)

// DataType represents the type of stored data
type DataType int32

const (
	DataType_DATA_TYPE_UNSPECIFIED DataType = 0
	DataType_DATA_TYPE_CREDENTIALS DataType = 1
	DataType_DATA_TYPE_TEXT        DataType = 2
	DataType_DATA_TYPE_BINARY      DataType = 3
	DataType_DATA_TYPE_CARD        DataType = 4
	
	// Backward compatibility
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
	SynchronizeData(ctx context.Context, req *SynchronizeDataRequest) (*SynchronizeDataResponse, error)
}

// DataServiceClient is the interface for data service client
type DataServiceClient interface {
	StoreData(ctx context.Context, req *StoreDataRequest) (*StoreDataResponse, error)
	GetData(ctx context.Context, req *GetDataRequest) (*GetDataResponse, error)
	GetUserData(ctx context.Context, req *GetUserDataRequest) (*GetUserDataResponse, error)
	UpdateData(ctx context.Context, req *UpdateDataRequest) (*UpdateDataResponse, error)
	DeleteData(ctx context.Context, req *DeleteDataRequest) (*DeleteDataResponse, error)
	SyncData(ctx context.Context, req *SyncDataRequest) (*SyncDataResponse, error)
	SynchronizeData(ctx context.Context, req *SynchronizeDataRequest) (*SynchronizeDataResponse, error)
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

// SyncAction represents the type of sync action
type SyncAction int32

const (
	SyncAction_SYNC_ACTION_UNSPECIFIED SyncAction = 0
	SyncAction_SYNC_ACTION_CREATE      SyncAction = 1
	SyncAction_SYNC_ACTION_UPDATE      SyncAction = 2
	SyncAction_SYNC_ACTION_DELETE      SyncAction = 3
)

// ConflictResolutionStrategy represents conflict resolution strategy
type ConflictResolutionStrategy int32

const (
	ConflictResolutionStrategy_CONFLICT_RESOLVE_UNSPECIFIED  ConflictResolutionStrategy = 0
	ConflictResolutionStrategy_CONFLICT_RESOLVE_BY_VERSION   ConflictResolutionStrategy = 1
	ConflictResolutionStrategy_CONFLICT_RESOLVE_BY_TIMESTAMP ConflictResolutionStrategy = 2
	ConflictResolutionStrategy_CONFLICT_RESOLVE_SERVER_WINS  ConflictResolutionStrategy = 3
	ConflictResolutionStrategy_CONFLICT_RESOLVE_CLIENT_WINS  ConflictResolutionStrategy = 4
)

// SyncItem represents an item for synchronization
type SyncItem struct {
	Action             SyncAction                 `json:"action"`
	DataItem           *DataItem                  `json:"data_item"`
	ConflictResolution ConflictResolutionStrategy `json:"conflict_resolution"`
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	ItemId           string     `json:"item_id"`
	Action           SyncAction `json:"action"`
	Success          bool       `json:"success"`
	Error            string     `json:"error"`
	ConflictResolved bool       `json:"conflict_resolved"`
	FinalVersion     int64      `json:"final_version"`
}

// SynchronizeDataRequest represents advanced synchronization request
type SynchronizeDataRequest struct {
	UserId           string                     `json:"user_id"`
	LastSyncTime     time.Time                  `json:"last_sync_time"`
	Items            []*SyncItem                `json:"items"`
	ConflictStrategy ConflictResolutionStrategy `json:"conflict_strategy"`
}

// SynchronizeDataResponse represents advanced synchronization response
type SynchronizeDataResponse struct {
	Success        bool          `json:"success"`
	Message        string        `json:"message"`
	Results        []*SyncResult `json:"results"`
	ServerChanges  []*SyncItem   `json:"server_changes"`
	SyncTimestamp  time.Time     `json:"sync_timestamp"`
	TotalConflicts int32         `json:"total_conflicts"`
}