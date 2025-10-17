package handlers

import (
	"context"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/fylgushev/go-diplom-final/internal/domain/services"
	"github.com/fylgushev/go-diplom-final/pkg/proto/data"
)

// SyncHandler implements sync gRPC service
type SyncHandler struct {
	syncService *services.SyncService
	dataService *services.DataService
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler(syncService *services.SyncService, dataService *services.DataService) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
		dataService: dataService,
	}
}

// SynchronizeData handles advanced bidirectional synchronization with conflict resolution
func (h *SyncHandler) SynchronizeData(ctx context.Context, req *data.SynchronizeDataRequest) (*data.SynchronizeDataResponse, error) {
	// Получаем время последней синхронизации
	lastSyncTime := req.LastSyncTime

	// Конвертируем protobuf items в доменные объекты
	syncItems := make([]services.SyncItem, 0, len(req.Items))
	for _, item := range req.Items {
		syncItem, err := h.convertProtoSyncItemToDomain(item)
		if err != nil {
			return &data.SynchronizeDataResponse{
				Success: false,
				Message: "failed to convert sync item: " + err.Error(),
			}, nil
		}
		syncItems = append(syncItems, syncItem)
	}

	// Создаем запрос синхронизации
	syncRequest := &services.SyncRequest{
		UserID:           req.UserId,
		LastSyncTime:     lastSyncTime,
		Items:            syncItems,
		ConflictStrategy: h.convertProtoConflictStrategy(req.ConflictStrategy),
	}

	// Выполняем синхронизацию
	syncResponse, err := h.syncService.SynchronizeData(ctx, syncRequest)
	if err != nil {
		return &data.SynchronizeDataResponse{
			Success: false,
			Message: "synchronization failed: " + err.Error(),
		}, nil
	}

	// Обновляем время последней синхронизации
	if err := h.syncService.UpdateLastSyncTime(ctx, req.UserId, syncResponse.SyncTimestamp); err != nil {
		// Логируем ошибку, но не прерываем операцию
	}

	// Конвертируем результат обратно в protobuf
	protoResponse := &data.SynchronizeDataResponse{
		Success:        syncResponse.Success,
		SyncTimestamp:  syncResponse.SyncTimestamp,
		TotalConflicts: int32(syncResponse.TotalConflicts),
	}

	// Конвертируем результаты
	protoResponse.Results = make([]*data.SyncResult, 0, len(syncResponse.Results))
	for _, result := range syncResponse.Results {
		protoResult, err := h.convertDomainSyncResultToProto(result)
		if err != nil {
			continue // Пропускаем элементы с ошибками
		}
		protoResponse.Results = append(protoResponse.Results, protoResult)
	}

	// Конвертируем изменения с сервера
	protoResponse.ServerChanges = make([]*data.SyncItem, 0, len(syncResponse.ServerChanges))
	for _, change := range syncResponse.ServerChanges {
		protoChange, err := h.convertDomainSyncItemToProto(change)
		if err != nil {
			continue // Пропускаем элементы с ошибками
		}
		protoResponse.ServerChanges = append(protoResponse.ServerChanges, protoChange)
	}

	return protoResponse, nil
}

// convertProtoSyncItemToDomain конвертирует protobuf SyncItem в доменный объект
func (h *SyncHandler) convertProtoSyncItemToDomain(item *data.SyncItem) (services.SyncItem, error) {
	dataItem, err := h.convertProtoDataItemToDomain(item.DataItem)
	if err != nil {
		return services.SyncItem{}, err
	}

	return services.SyncItem{
		Action:             h.convertProtoSyncAction(item.Action),
		DataItem:           dataItem,
		ConflictResolution: h.convertProtoConflictStrategy(item.ConflictResolution),
	}, nil
}

// convertDomainSyncItemToProto конвертирует доменный SyncItem в protobuf
func (h *SyncHandler) convertDomainSyncItemToProto(item services.SyncItem) (*data.SyncItem, error) {
	protoDataItem, err := h.convertDomainDataItemToProto(item.DataItem)
	if err != nil {
		return nil, err
	}

	return &data.SyncItem{
		Action:             h.convertDomainSyncActionToProto(item.Action),
		DataItem:           protoDataItem,
		ConflictResolution: h.convertDomainConflictStrategyToProto(item.ConflictResolution),
	}, nil
}

// convertDomainSyncResultToProto конвертирует доменный SyncResult в protobuf
func (h *SyncHandler) convertDomainSyncResultToProto(result services.SyncResult) (*data.SyncResult, error) {
	return &data.SyncResult{
		ItemId:           result.ItemID,
		Action:           h.convertDomainSyncActionToProto(result.Action),
		Success:          result.Success,
		Error:            result.Error,
		ConflictResolved: result.ConflictResolved,
		FinalVersion:     result.FinalVersion,
	}, nil
}

// convertProtoDataItemToDomain конвертирует protobuf DataItem в доменный объект
func (h *SyncHandler) convertProtoDataItemToDomain(item *data.DataItem) (*entities.DataItem, error) {
	dataItem := &entities.DataItem{
		ID:        item.ID,
		UserID:    item.UserID,
		Type:      h.convertProtoDataType(item.Type),
		Name:      "", // Имя извлекается из метаданных
		Data:      item.Data,
		Version:   item.Version,
		IsDeleted: item.IsDeleted,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	// Устанавливаем метаданные
	if item.Metadata != "" {
		if err := dataItem.SetMetadata(item.Metadata); err != nil {
			return nil, err
		}
	}

	return dataItem, nil
}

// convertDomainDataItemToProto конвертирует доменный DataItem в protobuf
func (h *SyncHandler) convertDomainDataItemToProto(item *entities.DataItem) (*data.DataItem, error) {
	protoItem := &data.DataItem{
		ID:        item.ID,
		UserID:    item.UserID,
		Type:      h.convertDomainDataTypeToProto(item.Type),
		Data:      item.Data,
		Version:   item.Version,
		IsDeleted: item.IsDeleted,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	// Конвертируем метаданные в строку
	if item.Metadata != nil {
		protoItem.Metadata = string(item.Metadata)
	}

	return protoItem, nil
}

// Conversion helpers for sync actions
func (h *SyncHandler) convertProtoSyncAction(action data.SyncAction) services.SyncAction {
	switch action {
	case data.SyncAction_SYNC_ACTION_CREATE:
		return services.SyncActionCreate
	case data.SyncAction_SYNC_ACTION_UPDATE:
		return services.SyncActionUpdate
	case data.SyncAction_SYNC_ACTION_DELETE:
		return services.SyncActionDelete
	default:
		return services.SyncActionUpdate
	}
}

func (h *SyncHandler) convertDomainSyncActionToProto(action services.SyncAction) data.SyncAction {
	switch action {
	case services.SyncActionCreate:
		return data.SyncAction_SYNC_ACTION_CREATE
	case services.SyncActionUpdate:
		return data.SyncAction_SYNC_ACTION_UPDATE
	case services.SyncActionDelete:
		return data.SyncAction_SYNC_ACTION_DELETE
	default:
		return data.SyncAction_SYNC_ACTION_UPDATE
	}
}

// Conversion helpers for conflict resolution strategies
func (h *SyncHandler) convertProtoConflictStrategy(strategy data.ConflictResolutionStrategy) services.ConflictResolutionStrategy {
	switch strategy {
	case data.ConflictResolutionStrategy_CONFLICT_RESOLVE_BY_VERSION:
		return services.ConflictResolveByVersion
	case data.ConflictResolutionStrategy_CONFLICT_RESOLVE_BY_TIMESTAMP:
		return services.ConflictResolveByTimestamp
	case data.ConflictResolutionStrategy_CONFLICT_RESOLVE_SERVER_WINS:
		return services.ConflictResolveServerWins
	case data.ConflictResolutionStrategy_CONFLICT_RESOLVE_CLIENT_WINS:
		return services.ConflictResolveClientWins
	default:
		return services.ConflictResolveByVersion
	}
}

func (h *SyncHandler) convertDomainConflictStrategyToProto(strategy services.ConflictResolutionStrategy) data.ConflictResolutionStrategy {
	switch strategy {
	case services.ConflictResolveByVersion:
		return data.ConflictResolutionStrategy_CONFLICT_RESOLVE_BY_VERSION
	case services.ConflictResolveByTimestamp:
		return data.ConflictResolutionStrategy_CONFLICT_RESOLVE_BY_TIMESTAMP
	case services.ConflictResolveServerWins:
		return data.ConflictResolutionStrategy_CONFLICT_RESOLVE_SERVER_WINS
	case services.ConflictResolveClientWins:
		return data.ConflictResolutionStrategy_CONFLICT_RESOLVE_CLIENT_WINS
	default:
		return data.ConflictResolutionStrategy_CONFLICT_RESOLVE_BY_VERSION
	}
}

// Conversion helpers for data types
func (h *SyncHandler) convertProtoDataType(dataType data.DataType) entities.DataType {
	switch dataType {
	case data.DataType_DATA_TYPE_CREDENTIALS:
		return entities.DataTypeLoginPassword
	case data.DataType_DATA_TYPE_TEXT:
		return entities.DataTypeText
	case data.DataType_DATA_TYPE_BINARY:
		return entities.DataTypeBinary
	case data.DataType_DATA_TYPE_CARD:
		return entities.DataTypeCreditCard
	default:
		return entities.DataTypeUnknown
	}
}

func (h *SyncHandler) convertDomainDataTypeToProto(dataType entities.DataType) data.DataType {
	switch dataType {
	case entities.DataTypeLoginPassword:
		return data.DataType_DATA_TYPE_CREDENTIALS
	case entities.DataTypeText:
		return data.DataType_DATA_TYPE_TEXT
	case entities.DataTypeBinary:
		return data.DataType_DATA_TYPE_BINARY
	case entities.DataTypeCreditCard:
		return data.DataType_DATA_TYPE_CARD
	default:
		return data.DataType_DATA_TYPE_UNSPECIFIED
	}
}