package entities

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============= User Tests =============

func TestNewUser(t *testing.T) {
	username := "testuser"
	user := NewUser(username)

	assert.NotEmpty(t, user.ID, "ID должен быть сгенерирован")
	assert.Equal(t, username, user.Username)
	assert.True(t, user.IsActive, "Новый пользователь должен быть активным")
}

func TestUser_SetPassword(t *testing.T) {
	user := NewUser("testuser")
	passwordHash := "hashed_password"
	salt := "random_salt"

	user.SetPassword(passwordHash, salt)

	assert.Equal(t, passwordHash, user.PasswordHash)
	assert.Equal(t, salt, user.Salt)
}

func TestUser_UpdateLastLogin(t *testing.T) {
	user := NewUser("testuser")
	beforeUpdate := time.Now()

	user.UpdateLastLogin()

	assert.NotNil(t, user.LastLoginAt)
	assert.True(t, user.LastLoginAt.After(beforeUpdate) || user.LastLoginAt.Equal(beforeUpdate))
}

func TestUser_Deactivate(t *testing.T) {
	user := NewUser("testuser")
	assert.True(t, user.IsActive)

	user.Deactivate()

	assert.False(t, user.IsActive)
}

func TestUser_Activate(t *testing.T) {
	user := NewUser("testuser")
	user.Deactivate()

	user.Activate()

	assert.True(t, user.IsActive)
}

func TestUser_UpdateLastSync(t *testing.T) {
	user := NewUser("testuser")
	beforeUpdate := time.Now()

	user.UpdateLastSync()

	assert.NotNil(t, user.LastSyncAt)
	assert.True(t, user.LastSyncAt.After(beforeUpdate) || user.LastSyncAt.Equal(beforeUpdate))
}

func TestUser_TableName(t *testing.T) {
	user := User{}
	assert.Equal(t, "users", user.TableName())
}

// ============= DataItem Tests =============

func TestNewDataItem(t *testing.T) {
	userID := "user-123"
	dataType := DataTypeLoginPassword
	name := "GitHub Login"
	data := []byte("encrypted_data")

	item := NewDataItem(userID, dataType, name, data)

	assert.NotEmpty(t, item.ID, "ID должен быть сгенерирован")
	assert.Equal(t, userID, item.UserID)
	assert.Equal(t, dataType, item.Type)
	assert.Equal(t, name, item.Name)
	assert.Equal(t, data, item.Data)
	assert.Equal(t, int64(1), item.Version)
	assert.False(t, item.IsDeleted)
}

func TestDataItem_SetMetadata(t *testing.T) {
	item := NewDataItem("user-123", DataTypeLoginPassword, "Test", []byte("data"))

	metadata := map[string]string{"url": "https://github.com", "category": "work"}
	err := item.SetMetadata(metadata)

	require.NoError(t, err)
	assert.NotNil(t, item.Metadata)

	// Проверка, что метаданные корректно сериализованы
	var retrievedMetadata map[string]string
	err = json.Unmarshal(item.Metadata, &retrievedMetadata)
	require.NoError(t, err)
	assert.Equal(t, metadata, retrievedMetadata)
}

func TestDataItem_SetMetadata_Nil(t *testing.T) {
	item := NewDataItem("user-123", DataTypeLoginPassword, "Test", []byte("data"))

	err := item.SetMetadata(nil)

	require.NoError(t, err)
	assert.Nil(t, item.Metadata)
}

func TestDataItem_GetMetadata(t *testing.T) {
	item := NewDataItem("user-123", DataTypeLoginPassword, "Test", []byte("data"))

	metadata := map[string]string{"url": "https://github.com", "category": "work"}
	err := item.SetMetadata(metadata)
	require.NoError(t, err)

	var retrievedMetadata map[string]string
	err = item.GetMetadata(&retrievedMetadata)

	require.NoError(t, err)
	assert.Equal(t, metadata, retrievedMetadata)
}

func TestDataItem_GetMetadata_Nil(t *testing.T) {
	item := NewDataItem("user-123", DataTypeLoginPassword, "Test", []byte("data"))

	var retrievedMetadata map[string]string
	err := item.GetMetadata(&retrievedMetadata)

	require.NoError(t, err)
}

func TestDataItem_IncrementVersion(t *testing.T) {
	item := NewDataItem("user-123", DataTypeLoginPassword, "Test", []byte("data"))
	initialVersion := item.Version

	item.IncrementVersion()

	assert.Equal(t, initialVersion+1, item.Version)
}

func TestDataItem_MarkAsDeleted(t *testing.T) {
	item := NewDataItem("user-123", DataTypeLoginPassword, "Test", []byte("data"))
	initialVersion := item.Version

	item.MarkAsDeleted()

	assert.True(t, item.IsDeleted)
	assert.Equal(t, initialVersion+1, item.Version)
}

func TestDataItem_Restore(t *testing.T) {
	item := NewDataItem("user-123", DataTypeLoginPassword, "Test", []byte("data"))
	item.MarkAsDeleted()
	versionAfterDelete := item.Version

	item.Restore()

	assert.False(t, item.IsDeleted)
	assert.Equal(t, versionAfterDelete+1, item.Version)
}

func TestDataItem_TableName(t *testing.T) {
	item := DataItem{}
	assert.Equal(t, "data_items", item.TableName())
}

// ============= LoginPasswordData Tests =============

func TestNewLoginPasswordData(t *testing.T) {
	login := "user@example.com"
	password := "secret123"

	lpd := NewLoginPasswordData(login, password)

	assert.Equal(t, login, lpd.Login)
	assert.Equal(t, password, lpd.Password)
	assert.Empty(t, lpd.URL)
	assert.Empty(t, lpd.Notes)
}

func TestLoginPasswordData_SetURL(t *testing.T) {
	lpd := NewLoginPasswordData("user", "pass")
	url := "https://example.com"

	lpd.SetURL(url)

	assert.Equal(t, url, lpd.URL)
}

func TestLoginPasswordData_SetNotes(t *testing.T) {
	lpd := NewLoginPasswordData("user", "pass")
	notes := "Test notes"

	lpd.SetNotes(notes)

	assert.Equal(t, notes, lpd.Notes)
}

func TestLoginPasswordData_Validate_Success(t *testing.T) {
	lpd := NewLoginPasswordData("user", "pass")

	err := lpd.Validate()

	assert.NoError(t, err)
}

func TestLoginPasswordData_Validate_EmptyLogin(t *testing.T) {
	lpd := NewLoginPasswordData("", "pass")

	err := lpd.Validate()

	assert.ErrorIs(t, err, ErrEmptyLogin)
}

func TestLoginPasswordData_Validate_EmptyPassword(t *testing.T) {
	lpd := NewLoginPasswordData("user", "")

	err := lpd.Validate()

	assert.ErrorIs(t, err, ErrEmptyPassword)
}

// ============= TextData Tests =============

func TestNewTextData(t *testing.T) {
	content := "Test content"

	td := NewTextData(content)

	assert.Equal(t, content, td.Content)
	assert.Equal(t, "plain", td.Format)
}

func TestTextData_SetFormat(t *testing.T) {
	td := NewTextData("content")
	format := "markdown"

	td.SetFormat(format)

	assert.Equal(t, format, td.Format)
}

func TestTextData_AddTag(t *testing.T) {
	td := NewTextData("content")

	td.AddTag("important")
	td.AddTag("work")

	assert.Len(t, td.Tags, 2)
	assert.Contains(t, td.Tags, "important")
	assert.Contains(t, td.Tags, "work")
}

func TestTextData_AddTag_Duplicate(t *testing.T) {
	td := NewTextData("content")
	td.AddTag("important")

	td.AddTag("important") // Попытка добавить дубликат

	assert.Len(t, td.Tags, 1, "Дубликаты не должны добавляться")
}

func TestTextData_RemoveTag(t *testing.T) {
	td := NewTextData("content")
	td.AddTag("tag1")
	td.AddTag("tag2")
	td.AddTag("tag3")

	td.RemoveTag("tag2")

	assert.Len(t, td.Tags, 2)
	assert.NotContains(t, td.Tags, "tag2")
}

func TestTextData_RemoveTag_NotFound(t *testing.T) {
	td := NewTextData("content")
	td.AddTag("tag1")

	td.RemoveTag("nonexistent")

	assert.Len(t, td.Tags, 1, "Размер не должен измениться")
}

func TestTextData_HasTag(t *testing.T) {
	td := NewTextData("content")
	td.AddTag("important")

	assert.True(t, td.HasTag("important"))
	assert.False(t, td.HasTag("nonexistent"))
}

func TestTextData_SetNotes(t *testing.T) {
	td := NewTextData("content")
	notes := "Test notes"

	td.SetNotes(notes)

	assert.Equal(t, notes, td.Notes)
}

func TestTextData_Validate_Success(t *testing.T) {
	td := NewTextData("content")

	err := td.Validate()

	assert.NoError(t, err)
}

func TestTextData_Validate_EmptyContent(t *testing.T) {
	td := &TextData{Content: ""}

	err := td.Validate()

	assert.ErrorIs(t, err, ErrEmptyContent)
}

func TestTextData_GetWordCount(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{"Empty", "", 0},
		{"Single word", "Hello", 1},
		{"Multiple words", "Hello world from Go", 4},
		{"With newlines", "Hello\nworld\ntest", 3},
		{"With tabs", "Hello\tworld\ttest", 3},
		{"Multiple spaces", "Hello   world", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := NewTextData(tt.content)
			count := td.GetWordCount()
			assert.Equal(t, tt.expected, count)
		})
	}
}

func TestTextData_GetCharCount(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{"Empty", "", 0},
		{"ASCII", "Hello", 5},
		{"Unicode", "Привет", 6},
		{"Mixed", "Hello мир", 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := NewTextData(tt.content)
			count := td.GetCharCount()
			assert.Equal(t, tt.expected, count)
		})
	}
}

// ============= BinaryData Tests =============

func TestNewBinaryData(t *testing.T) {
	content := []byte("binary content")
	fileName := "test.bin"

	bd := NewBinaryData(content, fileName)

	assert.Equal(t, content, bd.Content)
	assert.Equal(t, fileName, bd.FileName)
	assert.Equal(t, int64(len(content)), bd.Size)
}

func TestBinaryData_SetMimeType(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.png")
	mimeType := "image/png"

	bd.SetMimeType(mimeType)

	assert.Equal(t, mimeType, bd.MimeType)
}

func TestBinaryData_SetChecksum(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")
	checksum := "abc123"

	bd.SetChecksum(checksum)

	assert.Equal(t, checksum, bd.Checksum)
}

func TestBinaryData_AddTag(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")

	bd.AddTag("documents")
	bd.AddTag("important")

	assert.Len(t, bd.Tags, 2)
	assert.Contains(t, bd.Tags, "documents")
	assert.Contains(t, bd.Tags, "important")
}

func TestBinaryData_AddTag_Duplicate(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")
	bd.AddTag("documents")

	bd.AddTag("documents")

	assert.Len(t, bd.Tags, 1)
}

func TestBinaryData_RemoveTag(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")
	bd.AddTag("tag1")
	bd.AddTag("tag2")

	bd.RemoveTag("tag1")

	assert.Len(t, bd.Tags, 1)
	assert.NotContains(t, bd.Tags, "tag1")
}

func TestBinaryData_HasTag(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")
	bd.AddTag("important")

	assert.True(t, bd.HasTag("important"))
	assert.False(t, bd.HasTag("nonexistent"))
}

func TestBinaryData_SetNotes(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")
	notes := "Important file"

	bd.SetNotes(notes)

	assert.Equal(t, notes, bd.Notes)
}

func TestBinaryData_Validate_Success(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")

	err := bd.Validate()

	assert.NoError(t, err)
}

func TestBinaryData_Validate_EmptyContent(t *testing.T) {
	bd := NewBinaryData([]byte{}, "test.bin")

	err := bd.Validate()

	assert.ErrorIs(t, err, ErrEmptyContent)
}

func TestBinaryData_Validate_EmptyFileName(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "")

	err := bd.Validate()

	assert.ErrorIs(t, err, ErrEmptyFileName)
}

func TestBinaryData_Validate_InvalidSize(t *testing.T) {
	bd := NewBinaryData([]byte("content"), "test.bin")
	bd.Size = 999 // Неправильный размер

	err := bd.Validate()

	assert.ErrorIs(t, err, ErrInvalidSize)
}

func TestBinaryData_GetFileExtension(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"test.txt", ".txt"},
		{"document.PDF", ".pdf"},
		{"archive.tar.gz", ".gz"},
		{"noextension", ""},
	}

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			bd := NewBinaryData([]byte("content"), tt.fileName)
			ext := bd.GetFileExtension()
			assert.Equal(t, tt.expected, ext)
		})
	}
}

func TestBinaryData_GetHumanReadableSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			bd := NewBinaryData(make([]byte, tt.size), "test.bin")
			size := bd.GetHumanReadableSize()
			assert.Equal(t, tt.expected, size)
		})
	}
}

func TestBinaryData_IsImage(t *testing.T) {
	tests := []struct {
		fileName string
		mimeType string
		expected bool
	}{
		{"photo.jpg", "image/jpeg", true},
		{"photo.png", "", true},
		{"photo.gif", "", true},
		{"document.pdf", "", false},
		{"test.txt", "text/plain", false},
		{"unknown.xyz", "image/xyz", true}, // По MIME-типу
	}

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			bd := NewBinaryData([]byte("content"), tt.fileName)
			if tt.mimeType != "" {
				bd.SetMimeType(tt.mimeType)
			}
			isImage := bd.IsImage()
			assert.Equal(t, tt.expected, isImage)
		})
	}
}

func TestBinaryData_IsDocument(t *testing.T) {
	tests := []struct {
		fileName string
		mimeType string
		expected bool
	}{
		{"doc.pdf", "application/pdf", true},
		{"file.doc", "", true},
		{"file.docx", "", true},
		{"note.txt", "", true},
		{"photo.jpg", "", false},
		{"archive.zip", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			bd := NewBinaryData([]byte("content"), tt.fileName)
			if tt.mimeType != "" {
				bd.SetMimeType(tt.mimeType)
			}
			isDoc := bd.IsDocument()
			assert.Equal(t, tt.expected, isDoc)
		})
	}
}

// ============= CreditCardData Tests =============

func TestNewCreditCardData(t *testing.T) {
	number := "4532015112830366"
	expiryMonth := 12
	expiryYear := 2025
	cvv := "123"
	cardholderName := "John Doe"

	ccd := NewCreditCardData(number, expiryMonth, expiryYear, cvv, cardholderName)

	assert.Equal(t, number, ccd.Number)
	assert.Equal(t, expiryMonth, ccd.ExpiryMonth)
	assert.Equal(t, expiryYear, ccd.ExpiryYear)
	assert.Equal(t, cvv, ccd.CVV)
	assert.Equal(t, cardholderName, ccd.CardholderName)
}

func TestCreditCardData_SetBankName(t *testing.T) {
	ccd := NewCreditCardData("4532015112830366", 12, 2025, "123", "John Doe")
	bankName := "Test Bank"

	ccd.SetBankName(bankName)

	assert.Equal(t, bankName, ccd.BankName)
}

func TestCreditCardData_SetNotes(t *testing.T) {
	ccd := NewCreditCardData("4532015112830366", 12, 2025, "123", "John Doe")
	notes := "Personal card"

	ccd.SetNotes(notes)

	assert.Equal(t, notes, ccd.Notes)
}

func TestCreditCardData_Validate_Success(t *testing.T) {
	currentYear := time.Now().Year()
	ccd := NewCreditCardData("4532015112830366", 12, currentYear+1, "123", "John Doe")

	err := ccd.Validate()

	assert.NoError(t, err)
}

func TestCreditCardData_Validate_InvalidNumber(t *testing.T) {
	ccd := NewCreditCardData("1234567890123456", 12, 2025, "123", "John Doe")

	err := ccd.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "алгоритма Луна")
}

func TestCreditCardData_Validate_InvalidMonth(t *testing.T) {
	ccd := NewCreditCardData("4532015112830366", 13, 2025, "123", "John Doe")

	err := ccd.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "месяц")
}

func TestCreditCardData_Validate_PastYear(t *testing.T) {
	currentYear := time.Now().Year()
	ccd := NewCreditCardData("4532015112830366", 12, currentYear-1, "123", "John Doe")

	err := ccd.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "прошлом")
}

func TestCreditCardData_Validate_InvalidCVV(t *testing.T) {
	ccd := NewCreditCardData("4532015112830366", 12, 2025, "12", "John Doe")

	err := ccd.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CVV")
}

func TestCreditCardData_Validate_EmptyCardholderName(t *testing.T) {
	ccd := NewCreditCardData("4532015112830366", 12, 2025, "123", "")

	err := ccd.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "держателя")
}

func TestCreditCardData_GetMaskedNumber(t *testing.T) {
	ccd := NewCreditCardData("4532015112830366", 12, 2025, "123", "John Doe")

	masked := ccd.GetMaskedNumber()

	assert.Equal(t, "****-****-****-0366", masked)
}

func TestCreditCardData_GetMaskedNumber_ShortNumber(t *testing.T) {
	ccd := NewCreditCardData("123", 12, 2025, "123", "John Doe")

	masked := ccd.GetMaskedNumber()

	assert.Equal(t, "****", masked)
}

func TestCreditCardData_IsExpired(t *testing.T) {
	currentYear := time.Now().Year()
	currentMonth := int(time.Now().Month())

	tests := []struct {
		name     string
		month    int
		year     int
		expected bool
	}{
		{"Future year", 12, currentYear + 1, false},
		{"Past year", 12, currentYear - 1, true},
	}

	// Добавляем тест для текущего года с будущим месяцем только если мы не в декабре
	if currentMonth < 12 {
		tests = append(tests, struct {
			name     string
			month    int
			year     int
			expected bool
		}{"Current year future month", 12, currentYear, false})
	}

	// Добавляем тест для текущего года с прошлым месяцем только если мы не в январе
	if currentMonth > 1 {
		tests = append(tests, struct {
			name     string
			month    int
			year     int
			expected bool
		}{"Current year past month", currentMonth - 1, currentYear, true})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccd := NewCreditCardData("4532015112830366", tt.month, tt.year, "123", "John Doe")
			isExpired := ccd.IsExpired()
			assert.Equal(t, tt.expected, isExpired)
		})
	}
}

func TestCreditCardData_LuhnCheck(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{"Valid Visa", "4532015112830366", true},
		{"Valid Mastercard", "5425233430109903", true},
		{"Invalid", "1234567890123456", false},
		{"Valid with spaces", "4532 0151 1283 0366", true},
		{"Valid with dashes", "4532-0151-1283-0366", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentYear := time.Now().Year()
			ccd := NewCreditCardData(tt.number, 12, currentYear+1, "123", "John Doe")
			err := ccd.Validate()
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
