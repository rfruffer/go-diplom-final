// Package repository предоставляет реализации репозиториев для работы с базой данных.
//
// Этот пакет содержит реализации интерфейсов репозиториев для хранения
// данных пользователей и их элементов данных. Использует GORM для
// взаимодействия с PostgreSQL базой данных.
//
// Основные компоненты:
//   - UserRepository: управление пользователями (создание, поиск, обновление)
//   - DataRepository: управление элементами данных пользователей (CRUD операции, синхронизация)
//
// Пример использования UserRepository:
//
//	repo := repository.NewUserRepository(db)
//	user := entities.NewUser("username")
//	err := repo.CreateUser(ctx, user)
//
// Пример использования DataRepository:
//
//	repo := repository.NewDataRepository(db)
//	item := entities.NewDataItem(userID, dataType, name, encryptedData)
//	err := repo.CreateDataItem(ctx, item)
package repository
