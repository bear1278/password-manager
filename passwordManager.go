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
