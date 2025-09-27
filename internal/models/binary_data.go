package models

import (
	"fmt"
	"path/filepath"
	"strings"
)

// BinaryData представляет структуру данных для типа BinaryData.
// Содержит произвольные бинарные данные с метаинформацией о файле.
// Эта структура сериализуется в JSON и шифруется перед сохранением в DataItem.Data.
type BinaryData struct {
	// Content бинарное содержимое файла в формате base64.
	// Все бинарные данные кодируются в base64 для безопасного хранения в JSON.
	Content []byte `json:"content"`

	// FileName оригинальное имя файла.
	// Сохраняется для удобства пользователя и корректного восстановления файла.
	FileName string `json:"file_name"`

	// MimeType MIME-тип файла.
	// Определяет тип содержимого для корректной обработки и отображения.
	MimeType string `json:"mime_type,omitempty"`

	// Size размер файла в байтах.
	// Рассчитывается автоматически на основе содержимого.
	Size int64 `json:"size"`

	// Checksum контрольная сумма файла для проверки целостности.
	// Используется для верификации данных после расшифровки.
	Checksum string `json:"checksum,omitempty"`

	// Tags теги для категоризации и поиска бинарных данных.
	// Позволяет пользователю организовывать файлы по категориям.
	Tags []string `json:"tags,omitempty"`

	// Notes дополнительные заметки к файлу.
	// Может содержать описание, источник файла и другую метаинформацию.
	Notes string `json:"notes,omitempty"`
}

// NewBinaryData создает новый экземпляр бинарных данных.
// Принимает обязательные параметры: содержимое файла и имя файла.
// Автоматически рассчитывает размер файла.
func NewBinaryData(content []byte, fileName string) *BinaryData {
	return &BinaryData{
		Content:  content,
		FileName: fileName,
		Size:     int64(len(content)),
	}
}

// SetMimeType устанавливает MIME-тип для бинарных данных.
// Может быть определен автоматически на основе расширения файла.
func (bd *BinaryData) SetMimeType(mimeType string) {
	bd.MimeType = mimeType
}

// SetChecksum устанавливает контрольную сумму для проверки целостности.
// Рекомендуется использовать SHA-256 или другой криптографически стойкий алгоритм.
func (bd *BinaryData) SetChecksum(checksum string) {
	bd.Checksum = checksum
}

// AddTag добавляет тег к бинарным данным.
// Проверяет, что тег еще не существует, чтобы избежать дублирования.
func (bd *BinaryData) AddTag(tag string) {
	for _, existingTag := range bd.Tags {
		if existingTag == tag {
			return // тег уже существует
		}
	}
	bd.Tags = append(bd.Tags, tag)
}

// RemoveTag удаляет тег из бинарных данных.
// Если тег не найден, операция игнорируется.
func (bd *BinaryData) RemoveTag(tag string) {
	for i, existingTag := range bd.Tags {
		if existingTag == tag {
			bd.Tags = append(bd.Tags[:i], bd.Tags[i+1:]...)
			return
		}
	}
}

// HasTag проверяет наличие тега в бинарных данных.
// Возвращает true, если тег найден.
func (bd *BinaryData) HasTag(tag string) bool {
	for _, existingTag := range bd.Tags {
		if existingTag == tag {
			return true
		}
	}
	return false
}

// SetNotes устанавливает заметки для бинарных данных.
// Может содержать любую дополнительную информацию о файле.
func (bd *BinaryData) SetNotes(notes string) {
	bd.Notes = notes
}

// Validate проверяет корректность бинарных данных.
// Возвращает ошибку, если обязательные поля пустые или данные некорректны.
func (bd *BinaryData) Validate() error {
	if len(bd.Content) == 0 {
		return ErrEmptyContent
	}
	if bd.FileName == "" {
		return ErrEmptyFileName
	}
	if bd.Size != int64(len(bd.Content)) {
		return ErrInvalidSize
	}
	return nil
}

// GetFileExtension возвращает расширение файла.
// Извлекает расширение из имени файла для определения типа.
func (bd *BinaryData) GetFileExtension() string {
	return strings.ToLower(filepath.Ext(bd.FileName))
}

// GetHumanReadableSize возвращает размер файла в удобочитаемом формате.
// Преобразует байты в KB, MB, GB и т.д.
func (bd *BinaryData) GetHumanReadableSize() string {
	const unit = 1024
	size := float64(bd.Size)
	
	if size < unit {
		return fmt.Sprintf("%.0f B", size)
	}
	
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", size/float64(div), units[exp])
}

// IsImage проверяет, является ли файл изображением.
// Основывается на MIME-типе и расширении файла.
func (bd *BinaryData) IsImage() bool {
	if strings.HasPrefix(bd.MimeType, "image/") {
		return true
	}
	
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	ext := bd.GetFileExtension()
	
	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	
	return false
}

// IsDocument проверяет, является ли файл документом.
// Основывается на MIME-типе и расширении файла.
func (bd *BinaryData) IsDocument() bool {
	documentMimes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/plain",
		"text/rtf",
	}
	
	for _, mime := range documentMimes {
		if bd.MimeType == mime {
			return true
		}
	}
	
	documentExtensions := []string{".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt"}
	ext := bd.GetFileExtension()
	
	for _, docExt := range documentExtensions {
		if ext == docExt {
			return true
		}
	}
	
	return false
}