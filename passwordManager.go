package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

// PasswordManager manages encrypted password storage
type PasswordManager struct {
	passwords     map[string]Password `json:"passwords"`
	masterKey     []byte              `json:"-"`
	filePath      string              `json:"-"`
	isInitialized bool                `json:"-"`
}

// NewPasswordManager creates a new PasswordManager instance
func NewPasswordManager(filePath string) *PasswordManager {
	return &PasswordManager{filePath: filePath, passwords: make(map[string]Password), isInitialized: false}
}

/* Важно: Использование пароля напрямую как ключа AES — упрощение для учебных целей.
В production используют KDF (Key Derivation Function),
например scrypt или argon2, для безопасного преобразования пароля в ключ.*/

// SetMasterPassword sets and validates the master password
// Note: Using password directly as AES key is simplified for educational purposes.
// Production code should use KDF like scrypt or argon2.
func (pm *PasswordManager) SetMasterPassword(masterPassword string) error {
	if len(masterPassword) < 8 {
		return errors.New("master password must be at least 8 characters")
	}
	buffer := make([]byte, 32)
	copy(buffer, []byte(masterPassword))
	pm.masterKey = buffer
	pm.isInitialized = true
	return nil
}

// SavePassword saves a new password after validation
func (pm *PasswordManager) SavePassword(name, value, category string) error {
	if !pm.isInitialized {
		return errors.New("password manager is not initialized")
	}
	if _, ok := pm.passwords[name]; ok {
		return errors.New("password name already exists")
	}
	if err := pm.CheckPasswordStrength(value); err != nil {
		return err
	}
	password := NewPassword(name, value, strings.ToLower(category))
	pm.passwords[name] = password
	return nil
}

// GetPassword retrieves a password by name
func (pm *PasswordManager) GetPassword(name string) (Password, error) {
	if !pm.isInitialized {
		return Password{}, errors.New("password manager is not initialized")
	}
	p, ok := pm.passwords[name]
	if !ok {
		return Password{}, errors.New("password does not exist")
	}
	return p, nil
}

// ListPasswords returns all passwords as a slice
func (pm *PasswordManager) ListPasswords() []Password {
	result := make([]Password, 0)
	for _, v := range pm.passwords {
		result = append(result, v)
	}
	return result
}

// CheckPasswordStrength validates password complexity
func (pm *PasswordManager) CheckPasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	hasBig, hasDigit, hasSmall, hasSpec := false, false, false, false
	for _, v := range password {
		if strings.ContainsRune("0123456789", v) {
			hasDigit = true
		}
		if strings.ContainsRune("!@#$%^&*()_+=-/?.,<>';:]}[{~", v) {
			hasSpec = true
		}
		if strings.ContainsRune("qwertyuiopasdfghjklzxcvbnm", v) {
			hasSmall = true
		}
		if strings.ContainsRune("QWERTYUIOPASDFGHJKLZXCVBNM", v) {
			hasBig = true
		}
	}
	if !hasBig {
		return errors.New("password must contain capital characters")
	}
	if !hasSmall {
		return errors.New("password must contain lowercase characters")
	}
	if !hasSpec {
		return errors.New("password must contain special characters")
	}
	if !hasDigit {
		return errors.New("password must contain digits")
	}
	return nil
}

// GeneratePassword generates a cryptographically secure random password
func (pm *PasswordManager) GeneratePassword(length int) (string, error) {

	if length < 8 {
		return "", errors.New("password must be at least 8 characters")
	}
	symbols := `qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890!@#$%^&*()_+=-/?.,<>';:]}[{~`
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	stringBuilder := strings.Builder{}
	var key int
	for _, v := range buffer {
		key = int(v) % len(symbols)
		stringBuilder.WriteString(string(symbols[key]))
	}
	return stringBuilder.String(), nil
}

// SaveToFile encrypts and saves passwords to file
func (pm *PasswordManager) SaveToFile() error {
	if !pm.isInitialized {
		return errors.New("password manager is not initialized")
	}
	json, err := json.Marshal(pm.passwords)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(pm.masterKey)
	if err != nil {
		return err
	}
	gsm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gsm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return err
	}
	result := gsm.Seal(nil, nonce, json, nil)
	file, err := os.Create(pm.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(nonce)
	if err != nil {
		return err
	}
	_, err = file.Write(result)
	if err != nil {
		return err
	}
	return nil
}

// LoadFromFile loads and decrypts passwords from file
func (pm *PasswordManager) LoadFromFile() error {
	if !pm.isInitialized {
		return errors.New("password manager is not initialized")
	}
	file, err := os.Open(pm.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	block, err := aes.NewCipher(pm.masterKey)
	if err != nil {
		return err
	}
	gsm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gsm.NonceSize())
	_, err = io.ReadFull(file, nonce)
	if err != nil {
		return err
	}
	encryptedData, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	data, err := gsm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pm.passwords)
	if err != nil {
		return err
	}
	return nil
}

// GetPasswordsByCategory returns all passwords in a specific category
func (pm *PasswordManager) GetPasswordsByCategory(category string) []Password {
	result := make([]Password, 0)
	for _, v := range pm.passwords {
		if strings.ToLower(category) == v.Category {
			result = append(result, v)
		}
	}
	return result
}

// FindDuplicatePasswords finds passwords with duplicate values
func (pm *PasswordManager) FindDuplicatePasswords() map[string][]string {
	duplicates := make(map[string][]string)
	for _, v := range pm.passwords {
		if duplicates[v.Value] == nil {
			duplicates[v.Value] = make([]string, 0)
		}
		duplicates[v.Value] = append(duplicates[v.Value], v.Name)
	}
	for k, v := range duplicates {
		if len(v) == 1 {
			delete(duplicates, k)
		}
	}
	return duplicates
}

// UpdatePassword updates an existing password
func (pm *PasswordManager) UpdatePassword(name, newValue string) error {
	if !pm.isInitialized {
		return errors.New("password manager is not initialized")
	}
	password, ok := pm.passwords[name]
	if !ok {
		return errors.New("password does not exist")
	}
	if err := pm.CheckPasswordStrength(newValue); err != nil {
		return err
	}
	password.Value = newValue
	password.LastModified = time.Now()
	pm.passwords[name] = password
	return nil
}

// DeletePassword removes a password from storage
func (pm *PasswordManager) DeletePassword(name string) error {
	if !pm.isInitialized {
		return errors.New("password manager is not initialized")
	}
	if _, ok := pm.passwords[name]; !ok {
		return errors.New("password does not exist")
	}
	delete(pm.passwords, name)
	return nil
}

// ListCategories returns all unique categories
func (pm *PasswordManager) ListCategories() []string {
	categories := make([]string, 0)
	var multitude = make(map[string]bool)
	for _, v := range pm.passwords {
		if !multitude[v.Category] {
			multitude[v.Category] = true
		}
	}
	for k, _ := range multitude {
		categories = append(categories, k)
	}
	return categories
}

// GetPasswordStats returns statistics about stored passwords
func (pm *PasswordManager) GetPasswordStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["total"] = len(pm.passwords)
	categories := pm.ListCategories()
	for _, v := range categories {
		stats[v] = len(pm.GetPasswordsByCategory(v))
	}
	firstIteration := true
	var oldest, newest Password
	for _, v := range pm.passwords {
		if firstIteration {
			firstIteration = false
			oldest, newest = v, v
		}
		if v.CreatedAt.Before(oldest.CreatedAt) {
			oldest = v
		}
		if v.CreatedAt.After(newest.CreatedAt) {
			newest = v
		}
	}
	if len(pm.passwords) > 0 {
		stats["oldest"] = oldest.CreatedAt.Format(time.DateTime)
		stats["newest"] = newest.CreatedAt.Format(time.DateTime)
	} else {
		stats["oldest"], stats["newest"] = "", ""
	}
	return stats
}
