package services

import (
	"testing"
	"time"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataSerializer(t *testing.T) {
	ds := NewDataSerializer()
	assert.NotNil(t, ds)
}

// ========== LoginPassword Serialization Tests ==========

func TestDataSerializer_SerializeLoginPassword_Success(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewLoginPasswordData("user@example.com", "password123")

	serialized, err := ds.SerializeLoginPassword(data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeLoginPassword_ValidationError(t *testing.T) {
	ds := NewDataSerializer()
	data := &entities.LoginPasswordData{Login: "", Password: "pass"}

	serialized, err := ds.SerializeLoginPassword(data)

	assert.Error(t, err)
	assert.Nil(t, serialized)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestDataSerializer_DeserializeLoginPassword_Success(t *testing.T) {
	ds := NewDataSerializer()
	original := entities.NewLoginPasswordData("user@example.com", "password123")
	original.SetURL("https://example.com")

	serialized, err := ds.SerializeLoginPassword(original)
	require.NoError(t, err)

	deserialized, err := ds.DeserializeLoginPassword(serialized)

	require.NoError(t, err)
	assert.Equal(t, original.Login, deserialized.Login)
	assert.Equal(t, original.Password, deserialized.Password)
	assert.Equal(t, original.URL, deserialized.URL)
}

func TestDataSerializer_DeserializeLoginPassword_InvalidJSON(t *testing.T) {
	ds := NewDataSerializer()
	invalidData := []byte("invalid json")

	deserialized, err := ds.DeserializeLoginPassword(invalidData)

	assert.Error(t, err)
	assert.Nil(t, deserialized)
}

func TestDataSerializer_DeserializeLoginPassword_ValidationError(t *testing.T) {
	ds := NewDataSerializer()
	invalidData := []byte(`{"login":"","password":""}`)

	deserialized, err := ds.DeserializeLoginPassword(invalidData)

	assert.Error(t, err)
	assert.Nil(t, deserialized)
	assert.Contains(t, err.Error(), "validation failed")
}

// ========== TextData Serialization Tests ==========

func TestDataSerializer_SerializeTextData_Success(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewTextData("This is a test content")
	data.SetFormat("markdown")
	data.AddTag("important")

	serialized, err := ds.SerializeTextData(data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeTextData_ValidationError(t *testing.T) {
	ds := NewDataSerializer()
	data := &entities.TextData{Content: ""}

	serialized, err := ds.SerializeTextData(data)

	assert.Error(t, err)
	assert.Nil(t, serialized)
}

func TestDataSerializer_DeserializeTextData_Success(t *testing.T) {
	ds := NewDataSerializer()
	original := entities.NewTextData("This is a test content")
	original.SetFormat("markdown")
	original.AddTag("important")

	serialized, err := ds.SerializeTextData(original)
	require.NoError(t, err)

	deserialized, err := ds.DeserializeTextData(serialized)

	require.NoError(t, err)
	assert.Equal(t, original.Content, deserialized.Content)
	assert.Equal(t, original.Format, deserialized.Format)
	assert.Equal(t, original.Tags, deserialized.Tags)
}

func TestDataSerializer_DeserializeTextData_InvalidJSON(t *testing.T) {
	ds := NewDataSerializer()
	invalidData := []byte("invalid json")

	deserialized, err := ds.DeserializeTextData(invalidData)

	assert.Error(t, err)
	assert.Nil(t, deserialized)
}

// ========== CreditCard Serialization Tests ==========

func TestDataSerializer_SerializeCreditCard_Success(t *testing.T) {
	ds := NewDataSerializer()
	currentYear := time.Now().Year()
	data := entities.NewCreditCardData("4532015112830366", 12, currentYear+1, "123", "John Doe")

	serialized, err := ds.SerializeCreditCard(data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeCreditCard_ValidationError(t *testing.T) {
	ds := NewDataSerializer()
	data := &entities.CreditCardData{
		Number:         "",
		ExpiryMonth:    12,
		ExpiryYear:     2025,
		CVV:            "123",
		CardholderName: "John Doe",
	}

	serialized, err := ds.SerializeCreditCard(data)

	assert.Error(t, err)
	assert.Nil(t, serialized)
}

func TestDataSerializer_DeserializeCreditCard_Success(t *testing.T) {
	ds := NewDataSerializer()
	currentYear := time.Now().Year()
	original := entities.NewCreditCardData("4532015112830366", 12, currentYear+1, "123", "John Doe")
	original.SetBankName("Test Bank")

	serialized, err := ds.SerializeCreditCard(original)
	require.NoError(t, err)

	deserialized, err := ds.DeserializeCreditCard(serialized)

	require.NoError(t, err)
	assert.Equal(t, original.Number, deserialized.Number)
	assert.Equal(t, original.ExpiryMonth, deserialized.ExpiryMonth)
	assert.Equal(t, original.ExpiryYear, deserialized.ExpiryYear)
	assert.Equal(t, original.CVV, deserialized.CVV)
	assert.Equal(t, original.CardholderName, deserialized.CardholderName)
	assert.Equal(t, original.BankName, deserialized.BankName)
}

func TestDataSerializer_DeserializeCreditCard_InvalidJSON(t *testing.T) {
	ds := NewDataSerializer()
	invalidData := []byte("invalid json")

	deserialized, err := ds.DeserializeCreditCard(invalidData)

	assert.Error(t, err)
	assert.Nil(t, deserialized)
}

// ========== BinaryData Serialization Tests ==========

func TestDataSerializer_SerializeBinaryData_Success(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewBinaryData([]byte("binary content"), "test.bin")
	data.SetMimeType("application/octet-stream")

	serialized, err := ds.SerializeBinaryData(data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeBinaryData_ValidationError(t *testing.T) {
	ds := NewDataSerializer()
	data := &entities.BinaryData{
		Content:  []byte{},
		FileName: "",
	}

	serialized, err := ds.SerializeBinaryData(data)

	assert.Error(t, err)
	assert.Nil(t, serialized)
}

func TestDataSerializer_DeserializeBinaryData_Success(t *testing.T) {
	ds := NewDataSerializer()
	originalContent := []byte("binary content")
	original := entities.NewBinaryData(originalContent, "test.bin")
	original.SetMimeType("application/octet-stream")

	serialized, err := ds.SerializeBinaryData(original)
	require.NoError(t, err)

	deserialized, err := ds.DeserializeBinaryData(serialized)

	require.NoError(t, err)
	assert.Equal(t, original.Content, deserialized.Content)
	assert.Equal(t, original.FileName, deserialized.FileName)
	assert.Equal(t, original.MimeType, deserialized.MimeType)
	assert.Equal(t, original.Size, deserialized.Size)
}

func TestDataSerializer_DeserializeBinaryData_InvalidJSON(t *testing.T) {
	ds := NewDataSerializer()
	invalidData := []byte("invalid json")

	deserialized, err := ds.DeserializeBinaryData(invalidData)

	assert.Error(t, err)
	assert.Nil(t, deserialized)
}

// ========== SerializeByType Tests ==========

func TestDataSerializer_SerializeByType_LoginPassword(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewLoginPasswordData("user", "pass")

	serialized, err := ds.SerializeByType(entities.DataTypeLoginPassword, data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeByType_Text(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewTextData("content")

	serialized, err := ds.SerializeByType(entities.DataTypeText, data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeByType_CreditCard(t *testing.T) {
	ds := NewDataSerializer()
	currentYear := time.Now().Year()
	data := entities.NewCreditCardData("4532015112830366", 12, currentYear+1, "123", "John Doe")

	serialized, err := ds.SerializeByType(entities.DataTypeCreditCard, data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeByType_Binary(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewBinaryData([]byte("content"), "file.bin")

	serialized, err := ds.SerializeByType(entities.DataTypeBinary, data)

	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestDataSerializer_SerializeByType_InvalidType(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewTextData("content")

	serialized, err := ds.SerializeByType(entities.DataType(999), data)

	assert.Error(t, err)
	assert.Nil(t, serialized)
	assert.Contains(t, err.Error(), "unsupported data type")
}

func TestDataSerializer_SerializeByType_WrongDataType(t *testing.T) {
	ds := NewDataSerializer()
	data := entities.NewTextData("content")

	// Пытаемся сериализовать TextData как LoginPassword
	serialized, err := ds.SerializeByType(entities.DataTypeLoginPassword, data)

	assert.Error(t, err)
	assert.Nil(t, serialized)
	assert.Contains(t, err.Error(), "invalid data type")
}

// ========== DeserializeByType Tests ==========

func TestDataSerializer_DeserializeByType_LoginPassword(t *testing.T) {
	ds := NewDataSerializer()
	original := entities.NewLoginPasswordData("user", "pass")
	serialized, _ := ds.SerializeLoginPassword(original)

	deserialized, err := ds.DeserializeByType(entities.DataTypeLoginPassword, serialized)

	require.NoError(t, err)
	lpd, ok := deserialized.(*entities.LoginPasswordData)
	require.True(t, ok)
	assert.Equal(t, original.Login, lpd.Login)
}

func TestDataSerializer_DeserializeByType_Text(t *testing.T) {
	ds := NewDataSerializer()
	original := entities.NewTextData("content")
	serialized, _ := ds.SerializeTextData(original)

	deserialized, err := ds.DeserializeByType(entities.DataTypeText, serialized)

	require.NoError(t, err)
	td, ok := deserialized.(*entities.TextData)
	require.True(t, ok)
	assert.Equal(t, original.Content, td.Content)
}

func TestDataSerializer_DeserializeByType_CreditCard(t *testing.T) {
	ds := NewDataSerializer()
	currentYear := time.Now().Year()
	original := entities.NewCreditCardData("4532015112830366", 12, currentYear+1, "123", "John Doe")
	serialized, _ := ds.SerializeCreditCard(original)

	deserialized, err := ds.DeserializeByType(entities.DataTypeCreditCard, serialized)

	require.NoError(t, err)
	ccd, ok := deserialized.(*entities.CreditCardData)
	require.True(t, ok)
	assert.Equal(t, original.Number, ccd.Number)
}

func TestDataSerializer_DeserializeByType_Binary(t *testing.T) {
	ds := NewDataSerializer()
	original := entities.NewBinaryData([]byte("content"), "file.bin")
	serialized, _ := ds.SerializeBinaryData(original)

	deserialized, err := ds.DeserializeByType(entities.DataTypeBinary, serialized)

	require.NoError(t, err)
	bd, ok := deserialized.(*entities.BinaryData)
	require.True(t, ok)
	assert.Equal(t, original.FileName, bd.FileName)
}

func TestDataSerializer_DeserializeByType_InvalidType(t *testing.T) {
	ds := NewDataSerializer()
	data := []byte(`{"content":"test"}`)

	deserialized, err := ds.DeserializeByType(entities.DataType(999), data)

	assert.Error(t, err)
	assert.Nil(t, deserialized)
	assert.Contains(t, err.Error(), "unsupported data type")
}
