package handlers

import (
	"context"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/fylgushev/go-diplom-final/internal/domain/services"
	"github.com/fylgushev/go-diplom-final/pkg/proto/data"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	// Конвертируем protobuf items в доменные объекты
	syncItems, err := h.convertProtoSyncItems(req.Items)
	if err != nil {
		return &data.SynchronizeDataResponse{
			Success: false,
			Message: "failed to convert sync items: " + err.Error(),
		}, nil
	}

	// Выполняем синхронизацию
	syncResponse, err := h.performSync(ctx, req, syncItems)
	if err != nil {
		return &data.SynchronizeDataResponse{
			Success: false,
			Message: "synchronization failed: " + err.Error(),
		}, nil
	}

	// Конвертируем результат обратно в protobuf
	return h.buildSyncResponse(syncResponse), nil
}

// convertProtoSyncItems конвертирует список protobuf sync items в доменные объекты
func (h *SyncHandler) convertProtoSyncItems(items []*data.SyncItem) ([]services.SyncItem, error) {
	syncItems := make([]services.SyncItem, 0, len(items))
	for _, item := range items {
		syncItem, err := h.convertProtoSyncItemToDomain(item)
		if err != nil {
			return nil, err
		}
		syncItems = append(syncItems, syncItem)
	}
	return syncItems, nil
}

// performSync выполняет синхронизацию и обновляет время
func (h *SyncHandler) performSync(ctx context.Context, req *data.SynchronizeDataRequest, syncItems []services.SyncItem) (*services.SyncResponse, error) {
	syncRequest := &services.SyncRequest{
		UserID:           req.UserId,
		LastSyncTime:     req.LastSyncTime.AsTime(),
		Items:            syncItems,
		ConflictStrategy: h.convertProtoConflictStrategy(req.ConflictStrategy),
	}

	syncResponse, err := h.syncService.SynchronizeData(ctx, syncRequest)
	if err != nil {
		return nil, err
	}

	// Обновляем время последней синхронизации (ошибки игнорируем)
	_ = h.syncService.UpdateLastSyncTime(ctx, req.UserId, syncResponse.SyncTimestamp)

	return syncResponse, nil
}

// buildSyncResponse строит protobuf ответ из доменного объекта
func (h *SyncHandler) buildSyncResponse(syncResponse *services.SyncResponse) *data.SynchronizeDataResponse {
	protoResponse := &data.SynchronizeDataResponse{
		Success:        syncResponse.Success,
		SyncTimestamp:  timestamppb.New(syncResponse.SyncTimestamp),
		TotalConflicts: int32(syncResponse.TotalConflicts),
		Results:        h.convertSyncResults(syncResponse.Results),
		ServerChanges:  h.convertServerChanges(syncResponse.ServerChanges),
	}
	return protoResponse
}

// convertSyncResults конвертирует результаты синхронизации
func (h *SyncHandler) convertSyncResults(results []services.SyncResult) []*data.SyncResult {
	protoResults := make([]*data.SyncResult, 0, len(results))
	for _, result := range results {
		if protoResult, err := h.convertDomainSyncResultToProto(result); err == nil {
			protoResults = append(protoResults, protoResult)
		}
	}
	return protoResults
}

// convertServerChanges конвертирует изменения с сервера
func (h *SyncHandler) convertServerChanges(changes []services.SyncItem) []*data.SyncItem {
	protoChanges := make([]*data.SyncItem, 0, len(changes))
	for _, change := range changes {
		if protoChange, err := h.convertDomainSyncItemToProto(change); err == nil {
			protoChanges = append(protoChanges, protoChange)
		}
	}
	return protoChanges
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
		ID:        item.Id,
		UserID:    item.UserId,
		Type:      h.convertProtoDataType(item.Type),
		Name:      "", // Имя извлекается из метаданных
		Data:      item.Data,
		Version:   item.Version,
		IsDeleted: item.IsDeleted,
		CreatedAt: item.CreatedAt.AsTime(),
		UpdatedAt: item.UpdatedAt.AsTime(),
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
		Id:        item.ID,
		UserId:    item.UserID,
		Type:      h.convertDomainDataTypeToProto(item.Type),
		Data:      item.Data,
		Version:   item.Version,
		IsDeleted: item.IsDeleted,
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
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