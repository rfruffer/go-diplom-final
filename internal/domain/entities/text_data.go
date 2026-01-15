package entities

// TextData представляет структуру данных для типа TextData.
// Содержит произвольные текстовые данные с дополнительной информацией.
// Эта структура сериализуется в JSON и шифруется перед сохранением в DataItem.Data.
type TextData struct {
	// Content основное содержимое текстовых данных.
	// Может содержать заметки, коды активации, секретные вопросы,
	// списки одноразовых паролей и любую другую текстовую информацию.
	Content string `json:"content"`

	// Format формат текстовых данных.
	// Может указывать на тип содержимого: "plain", "markdown", "code", "json" и т.д.
	// Используется клиентом для корректного отображения данных.
	Format string `json:"format,omitempty"`

	// Tags теги для категоризации и поиска текстовых данных.
	// Позволяет пользователю организовывать данные по категориям.
	Tags []string `json:"tags,omitempty"`

	// Notes дополнительные заметки к текстовым данным.
	// Может содержать описание, инструкции по использованию и другую метаинформацию.
	Notes string `json:"notes,omitempty"`
}

// NewTextData создает новый экземпляр текстовых данных.
// Принимает обязательный параметр content с основным содержимым.
// Дополнительные поля могут быть установлены через соответствующие методы.
func NewTextData(content string) *TextData {
	return &TextData{
		Content: content,
		Format:  "plain", // значение по умолчанию
	}
}

// SetFormat устанавливает формат текстовых данных.
// Поддерживаемые форматы: "plain", "markdown", "code", "json", "xml", "yaml".
func (td *TextData) SetFormat(format string) {
	td.Format = format
}

// AddTag добавляет тег к текстовым данным.
// Проверяет, что тег еще не существует, чтобы избежать дублирования.
func (td *TextData) AddTag(tag string) {
	for _, existingTag := range td.Tags {
		if existingTag == tag {
			return // тег уже существует
		}
	}
	td.Tags = append(td.Tags, tag)
}

// RemoveTag удаляет тег из текстовых данных.
// Если тег не найден, операция игнорируется.
func (td *TextData) RemoveTag(tag string) {
	for i, existingTag := range td.Tags {
		if existingTag == tag {
			td.Tags = append(td.Tags[:i], td.Tags[i+1:]...)
			return
		}
	}
}

// HasTag проверяет наличие тега в текстовых данных.
// Возвращает true, если тег найден.
func (td *TextData) HasTag(tag string) bool {
	for _, existingTag := range td.Tags {
		if existingTag == tag {
			return true
		}
	}
	return false
}

// SetNotes устанавливает заметки для текстовых данных.
// Может содержать любую дополнительную информацию о данных.
func (td *TextData) SetNotes(notes string) {
	td.Notes = notes
}

// Validate проверяет корректность текстовых данных.
// Возвращает ошибку, если обязательные поля пустые или данные некорректны.
func (td *TextData) Validate() error {
	if td.Content == "" {
		return ErrEmptyContent
	}
	return nil
}

// GetWordCount возвращает количество слов в содержимом.
// Используется для статистики и отображения информации о данных.
func (td *TextData) GetWordCount() int {
	if td.Content == "" {
		return 0
	}
	
	// Простой подсчет слов по пробелам
	// Можно улучшить с помощью регулярных выражений для более точного подсчета
	words := 0
	inWord := false
	
	for _, char := range td.Content {
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	
	return words
}

// GetCharCount возвращает количество символов в содержимом.
// Используется для статистики и проверки лимитов.
func (td *TextData) GetCharCount() int {
	return len([]rune(td.Content))
}