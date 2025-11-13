package services

// Cipher интерфейс для шифрования и расшифровки данных
type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}
