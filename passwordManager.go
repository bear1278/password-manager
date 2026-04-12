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
)

type PasswordManager struct {
	passwords     map[string]Password `json:"passwords"`
	masterKey     []byte              `json:"-"`
	filePath      string              `json:"-"`
	isInitialized bool                `json:"-"`
}

func NewPasswordManager(filePath string) *PasswordManager {
	return &PasswordManager{filePath: filePath, passwords: make(map[string]Password), isInitialized: false}
}

/* Важно: Использование пароля напрямую как ключа AES — упрощение для учебных целей.
В production используют KDF (Key Derivation Function),
например scrypt или argon2, для безопасного преобразования пароля в ключ.*/

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

// TODO: После реализации CheckPasswordStrength (Этап 10) добавьте в SavePassword проверку надёжности пароля перед сохранением.
func (pm *PasswordManager) SavePassword(name, value, category string) error {
	if !pm.isInitialized {
		return errors.New("password manager is not initialized")
	}
	if _, ok := pm.passwords[name]; ok {
		return errors.New("password name already exists")
	}
	password := NewPassword(name, value, category)
	pm.passwords[name] = password
	return nil
}

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

func (pm *PasswordManager) ListPasswords() []Password {
	result := make([]Password, 0)
	for _, v := range pm.passwords {
		result = append(result, v)
	}
	return result
}

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
