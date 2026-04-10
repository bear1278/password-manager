package main

import "errors"

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
