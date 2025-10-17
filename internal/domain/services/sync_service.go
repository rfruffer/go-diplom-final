package services

import (
	"context"
	"fmt"
	"time"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/fylgushev/go-diplom-final/internal/domain/interfaces"
)

// SyncAction представляет тип действия синхронизации
type SyncAction string

const (
	// SyncActionCreate создание нового элемента
	SyncActionCreate SyncAction = "create"
	// SyncActionUpdate обновление существующего элемента
	SyncActionUpdate SyncAction = "update"
	// SyncActionDelete удаление элемента
	SyncActionDelete SyncAction = "delete"
)

// SyncItem представляет элемент для синхронизации
type SyncItem struct {
	// Action тип действия для синхронизации
	Action SyncAction `json:"action"`
	// DataItem элемент данных
	DataItem *entities.DataItem `json:"data_item"`
	// ConflictResolution стратегия разрешения конфликта
	ConflictResolution ConflictResolutionStrategy `json:"conflict_resolution,omitempty"`
}

// ConflictResolutionStrategy стратегия разрешения конфликтов
type ConflictResolutionStrategy string

const (
	// ConflictResolveByVersion разрешить по версии (берем более новую)
	ConflictResolveByVersion ConflictResolutionStrategy = "by_version"
	// ConflictResolveByTimestamp разрешить по времени (берем более новую)
	ConflictResolveByTimestamp ConflictResolutionStrategy = "by_timestamp"
	// ConflictResolveServerWins сервер всегда выигрывает
	ConflictResolveServerWins ConflictResolutionStrategy = "server_wins"
	// ConflictResolveClientWins клиент всегда выигрывает
	ConflictResolveClientWins ConflictResolutionStrategy = "client_wins"
)

// SyncResult результат операции синхронизации
type SyncResult struct {
	// ItemID идентификатор элемента
	ItemID string `json:"item_id"`
	// Action выполненное действие
	Action SyncAction `json:"action"`
	// Success успешность операции
	Success bool `json:"success"`
	// Error ошибка (если есть)
	Error string `json:"error,omitempty"`
	// ConflictResolved был ли разрешен конфликт
	ConflictResolved bool `json:"conflict_resolved,omitempty"`
	// FinalVersion финальная версия элемента
	FinalVersion int64 `json:"final_version"`
}

// SyncRequest запрос синхронизации от клиента
type SyncRequest struct {
	// UserID идентификатор пользователя
	UserID string `json:"user_id"`
	// LastSyncTime время последней синхронизации
	LastSyncTime time.Time `json:"last_sync_time"`
	// Items элементы для синхронизации
	Items []SyncItem `json:"items"`
	// ConflictStrategy стратегия разрешения конфликтов по умолчанию
	ConflictStrategy ConflictResolutionStrategy `json:"conflict_strategy"`
}

// SyncResponse ответ сервера на синхронизацию
type SyncResponse struct {
	// Success успешность операции синхронизации
	Success bool `json:"success"`
	// Results результаты обработки каждого элемента
	Results []SyncResult `json:"results"`
	// ServerChanges изменения с сервера для клиента
	ServerChanges []SyncItem `json:"server_changes"`
	// SyncTimestamp время синхронизации
	SyncTimestamp time.Time `json:"sync_timestamp"`
	// TotalConflicts количество разрешенных конфликтов
	TotalConflicts int `json:"total_conflicts"`
}

// SyncService сервис для управления синхронизацией данных между клиентами
type SyncService struct {
	dataRepo interfaces.DataRepository
}

// NewSyncService создает новый сервис синхронизации
func NewSyncService(dataRepo interfaces.DataRepository) *SyncService {
	return &SyncService{
		dataRepo: dataRepo,
	}
}

// SynchronizeData выполняет синхронизацию данных пользователя
func (s *SyncService) SynchronizeData(ctx context.Context, request *SyncRequest) (*SyncResponse, error) {
	response := &SyncResponse{
		Success:       true,
		Results:       make([]SyncResult, 0, len(request.Items)),
		ServerChanges: make([]SyncItem, 0),
		SyncTimestamp: time.Now(),
	}

	// Обрабатываем изменения от клиента
	for _, item := range request.Items {
		result := s.processSyncItem(ctx, item, request.ConflictStrategy)
		response.Results = append(response.Results, result)
		
		if result.ConflictResolved {
			response.TotalConflicts++
		}
		
		if !result.Success {
			response.Success = false
		}
	}

	// Получаем изменения с сервера для клиента
	serverChanges, err := s.getServerChanges(ctx, request.UserID, request.LastSyncTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get server changes: %w", err)
	}
	response.ServerChanges = serverChanges

	return response, nil
}

// processSyncItem обрабатывает один элемент синхронизации
func (s *SyncService) processSyncItem(ctx context.Context, item SyncItem, defaultStrategy ConflictResolutionStrategy) SyncResult {
	result := SyncResult{
		ItemID:  item.DataItem.ID,
		Action:  item.Action,
		Success: true,
	}

	switch item.Action {
	case SyncActionCreate:
		result = s.handleCreate(ctx, item.DataItem)
	case SyncActionUpdate:
		result = s.handleUpdate(ctx, item.DataItem, item.ConflictResolution, defaultStrategy)
	case SyncActionDelete:
		result = s.handleDelete(ctx, item.DataItem, item.ConflictResolution, defaultStrategy)
	default:
		result.Success = false
		result.Error = fmt.Sprintf("unknown sync action: %s", item.Action)
	}

	return result
}

// handleCreate обрабатывает создание нового элемента
func (s *SyncService) handleCreate(ctx context.Context, item *entities.DataItem) SyncResult {
	result := SyncResult{
		ItemID:       item.ID,
		Action:       SyncActionCreate,
		Success:      true,
		FinalVersion: item.Version,
	}

	// Проверяем, существует ли элемент
	existing, err := s.dataRepo.GetDataItem(ctx, item.UserID, item.ID)
	if err == nil && existing != nil {
		// Элемент уже существует - это конфликт
		result.Success = false
		result.Error = "item already exists"
		result.ConflictResolved = true
		result.FinalVersion = existing.Version
		return result
	}

	// Создаем новый элемент
	if err := s.dataRepo.CreateDataItem(ctx, item); err != nil {
		result.Success = false
		result.Error = err.Error()
	}

	return result
}

// handleUpdate обрабатывает обновление элемента
func (s *SyncService) handleUpdate(ctx context.Context, item *entities.DataItem, strategy, defaultStrategy ConflictResolutionStrategy) SyncResult {
	result := SyncResult{
		ItemID:       item.ID,
		Action:       SyncActionUpdate,
		Success:      true,
		FinalVersion: item.Version,
	}

	// Получаем текущую версию элемента
	existing, err := s.dataRepo.GetDataItem(ctx, item.UserID, item.ID)
	if err != nil {
		result.Success = false
		result.Error = "item not found"
		return result
	}

	// Проверяем наличие конфликта версий
	if existing.Version != item.Version-1 {
		// Конфликт версий - нужно разрешить
		resolved, err := s.resolveConflict(ctx, existing, item, strategy, defaultStrategy)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			return result
		}
		
		result.ConflictResolved = true
		result.FinalVersion = resolved.Version
		
		// Обновляем элемент разрешенной версией
		if err := s.dataRepo.UpdateDataItem(ctx, resolved); err != nil {
			result.Success = false
			result.Error = err.Error()
		}
		
		return result
	}

	// Нет конфликта - просто обновляем
	if err := s.dataRepo.UpdateDataItem(ctx, item); err != nil {
		result.Success = false
		result.Error = err.Error()
	}

	return result
}

// handleDelete обрабатывает удаление элемента
func (s *SyncService) handleDelete(ctx context.Context, item *entities.DataItem, strategy, defaultStrategy ConflictResolutionStrategy) SyncResult {
	result := SyncResult{
		ItemID:       item.ID,
		Action:       SyncActionDelete,
		Success:      true,
		FinalVersion: item.Version,
	}

	// Получаем текущую версию элемента
	existing, err := s.dataRepo.GetDataItem(ctx, item.UserID, item.ID)
	if err != nil {
		// Элемент не найден - возможно уже удален
		result.Success = true
		return result
	}

	// Проверяем конфликт версий
	if existing.Version != item.Version-1 {
		// Конфликт - разрешаем
		resolved, err := s.resolveConflict(ctx, existing, item, strategy, defaultStrategy)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			return result
		}
		
		result.ConflictResolved = true
		result.FinalVersion = resolved.Version
		
		// Если после разрешения конфликта элемент должен быть удален
		if resolved.IsDeleted {
			if err := s.dataRepo.DeleteDataItem(ctx, item.UserID, item.ID); err != nil {
				result.Success = false
				result.Error = err.Error()
			}
		} else {
			if err := s.dataRepo.UpdateDataItem(ctx, resolved); err != nil {
				result.Success = false
				result.Error = err.Error()
			}
		}
		
		return result
	}

	// Выполняем мягкое удаление
	existing.MarkAsDeleted()
	if err := s.dataRepo.UpdateDataItem(ctx, existing); err != nil {
		result.Success = false
		result.Error = err.Error()
	} else {
		result.FinalVersion = existing.Version
	}

	return result
}

// resolveConflict разрешает конфликт между двумя версиями элемента
func (s *SyncService) resolveConflict(ctx context.Context, server, client *entities.DataItem, strategy, defaultStrategy ConflictResolutionStrategy) (*entities.DataItem, error) {
	// Используем указанную стратегию или стратегию по умолчанию
	resolveStrategy := strategy
	if resolveStrategy == "" {
		resolveStrategy = defaultStrategy
	}
	if resolveStrategy == "" {
		resolveStrategy = ConflictResolveByVersion
	}

	switch resolveStrategy {
	case ConflictResolveByVersion:
		// Берем элемент с большей версией
		if server.Version >= client.Version {
			return server, nil
		}
		return client, nil

	case ConflictResolveByTimestamp:
		// Берем элемент с более поздним временем обновления
		if server.UpdatedAt.After(client.UpdatedAt) || server.UpdatedAt.Equal(client.UpdatedAt) {
			return server, nil
		}
		return client, nil

	case ConflictResolveServerWins:
		// Сервер всегда выигрывает
		return server, nil

	case ConflictResolveClientWins:
		// Клиент всегда выигрывает, но увеличиваем версию
		client.Version = server.Version + 1
		return client, nil

	default:
		return nil, fmt.Errorf("unknown conflict resolution strategy: %s", resolveStrategy)
	}
}

// getServerChanges получает изменения с сервера для клиента
func (s *SyncService) getServerChanges(ctx context.Context, userID string, lastSyncTime time.Time) ([]SyncItem, error) {
	// Получаем все элементы пользователя, измененные после последней синхронизации
	items, _, err := s.dataRepo.GetUserDataItemsModifiedAfter(ctx, userID, lastSyncTime)
	if err != nil {
		return nil, err
	}

	changes := make([]SyncItem, 0, len(items))
	for _, item := range items {
		action := SyncActionUpdate
		if item.CreatedAt.After(lastSyncTime) {
			action = SyncActionCreate
		}
		if item.IsDeleted {
			action = SyncActionDelete
		}

		changes = append(changes, SyncItem{
			Action:   action,
			DataItem: item,
		})
	}

	return changes, nil
}

// GetLastSyncTime получает время последней синхронизации пользователя
func (s *SyncService) GetLastSyncTime(ctx context.Context, userID string) (time.Time, error) {
	return s.dataRepo.GetUserLastSyncTime(ctx, userID)
}

// UpdateLastSyncTime обновляет время последней синхронизации пользователя
func (s *SyncService) UpdateLastSyncTime(ctx context.Context, userID string, syncTime time.Time) error {
	return s.dataRepo.UpdateUserLastSyncTime(ctx, userID, syncTime)
}