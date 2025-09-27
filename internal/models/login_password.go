package models

// LoginPasswordData представляет структуру данных для типа LoginPassword.
// Содержит пары логин/пароль с дополнительной информацией.
// Эта структура сериализуется в JSON и шифруется перед сохранением в DataItem.Data.
type LoginPasswordData struct {
	// Login логин пользователя для доступа к сервису.
	// Может быть именем пользователя, email адресом или другим идентификатором.
	Login string `json:"login"`

	// Password пароль для доступа к сервису.
	// В зашифрованном виде хранится в базе данных.
	Password string `json:"password"`

	// URL адрес сервиса, к которому относятся логин и пароль.
	// Опциональное поле для удобства пользователя.
	URL string `json:"url,omitempty"`

	// Notes дополнительные заметки к данным логина/пароля.
	// Может содержать информацию о секретных вопросах, двухфакторной аутентификации и т.д.
	Notes string `json:"notes,omitempty"`
}

// NewLoginPasswordData создает новый экземпляр данных логин/пароль.
// Принимает обязательные параметры логин и пароль.
// URL и заметки могут быть установлены позже через соответствующие методы.
func NewLoginPasswordData(login, password string) *LoginPasswordData {
	return &LoginPasswordData{
		Login:    login,
		Password: password,
	}
}

// SetURL устанавливает URL сервиса для данных логин/пароль.
// Используется для связывания учетных данных с конкретным веб-сайтом или сервисом.
func (lpd *LoginPasswordData) SetURL(url string) {
	lpd.URL = url
}

// SetNotes устанавливает заметки для данных логин/пароль.
// Может содержать любую дополнительную информацию, полезную пользователю.
func (lpd *LoginPasswordData) SetNotes(notes string) {
	lpd.Notes = notes
}

// Validate проверяет корректность данных логин/пароль.
// Возвращает ошибку, если обязательные поля пустые или данные некорректны.
func (lpd *LoginPasswordData) Validate() error {
	if lpd.Login == "" {
		return ErrEmptyLogin
	}
	if lpd.Password == "" {
		return ErrEmptyPassword
	}
	return nil
}