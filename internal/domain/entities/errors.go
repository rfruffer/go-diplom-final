package entities

import "errors"

// Ошибки валидации для модели LoginPasswordData
var (
	// ErrEmptyLogin возникает при попытке создать данные логин/пароль с пустым логином.
	ErrEmptyLogin = errors.New("логин не может быть пустым")

	// ErrEmptyPassword возникает при попытке создать данные логин/пароль с пустым паролем.
	ErrEmptyPassword = errors.New("пароль не может быть пустым")
)

// Ошибки валидации для модели TextData
var (
	// ErrEmptyContent возникает при попытке создать текстовые данные с пустым содержимым.
	ErrEmptyContent = errors.New("содержимое не может быть пустым")
)

// Ошибки валидации для модели BinaryData
var (
	// ErrEmptyFileName возникает при попытке создать бинарные данные с пустым именем файла.
	ErrEmptyFileName = errors.New("имя файла не может быть пустым")

	// ErrInvalidSize возникает когда размер файла не соответствует длине содержимого.
	ErrInvalidSize = errors.New("размер файла не соответствует содержимому")
)

// Ошибки для модели User
var (
	// ErrUserNotFound возникает при попытке найти несуществующего пользователя.
	ErrUserNotFound = errors.New("пользователь не найден")

	// ErrUserAlreadyExists возникает при попытке создать пользователя с уже существующим именем.
	ErrUserAlreadyExists = errors.New("пользователь с таким именем уже существует")

	// ErrInvalidCredentials возникает при неверных учетных данных при аутентификации.
	ErrInvalidCredentials = errors.New("неверные учетные данные")

	// ErrUserDeactivated возникает при попытке аутентификации деактивированного пользователя.
	ErrUserDeactivated = errors.New("пользователь деактивирован")
)

// Ошибки для модели DataItem
var (
	// ErrDataItemNotFound возникает при попытке найти несуществующий элемент данных.
	ErrDataItemNotFound = errors.New("элемент данных не найден")

	// ErrDataItemAlreadyDeleted возникает при попытке операции с уже удаленным элементом.
	ErrDataItemAlreadyDeleted = errors.New("элемент данных уже удален")

	// ErrInvalidDataType возникает при использовании неподдерживаемого типа данных.
	ErrInvalidDataType = errors.New("неподдерживаемый тип данных")

	// ErrVersionConflict возникает при конфликте версий во время синхронизации.
	ErrVersionConflict = errors.New("конфликт версий элемента данных")

	// ErrAccessDenied возникает при попытке доступа к чужим данным.
	ErrAccessDenied = errors.New("доступ запрещен")
)