# 🔐 Password Manager

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Консольный менеджер паролей с шифрованием AES-GCM, генерацией надежных паролей и управлением категориями.

## ✨ Функциональность

- **Безопасное хранение** - все пароли шифруются с использованием AES-GCM
- **Генерация паролей** - создание криптостойких паролей заданной длины
- **Управление категориями** - организация паролей по категориям
- **Поиск дубликатов** - обнаружение повторяющихся паролей
- **Статистика** - анализ хранилища паролей
- **Мастер-пароль** - единый ключ для доступа ко всем данным

## 📸 Демонстрация работы
```bash
==========================================
Password Manager
==========================================
1. Generate new password
2. Add new password
3. Get password
4. List all passwords
5. Update password
6. Delete password
7. List categories
8. Show password statistics
9. Find duplicate passwords
0. Exit
==========================================
Enter choice:

```

## 💻 Требования к системе

- Go 1.21 или выше
- Терминал с поддержкой цветового вывода (рекомендуется)
- Операционная система: Linux, macOS, Windows (WSL)

## 🔧 Установка

### Из исходного кода

```bash
# Клонирование репозитория
git clone https://github.com/yourusername/password-manager.git
cd password-manager

# Сборка проекта
go build .

# Запуск
./bin/password-manager
```

# 📖 Использование

### Инициализация

При первом запуске программа запросит мастер-пароль (минимум 8 символов):

```bash
Enter master password: 
✓ Success: Password saved successfully
```

### Добавление пароля
```bash
=== Add New Password ===
Enter service name: github.com
Enter password (or press Enter to generate): 
→ Info: Password generated successfully: K9#mP2$vL8@qR5
Enter category: development
✓ Success: Password saved successfully
```

### Поиск пароля
```bash
=== Search Password ===
Enter service name: github.com
✓ Success: Password search successfully
=== Password details ===
Service github.com
Category development
Password K9#mP2$vL8@qR5
Created 2024-01-15 14:30:25
Last Modified 2024-01-15 14:30:25
```

# 📁 Структура проекта

```text
password-manager/
├── main.go          # Точка входа, обработка пользовательского ввода
├── passwordManager.go       # Основная логика, шифрование, работа с файлами
├── password.go      # Модель данных пароля
├── screen.go            # Функции интерфейса пользователя
├── go.mod           # Управление зависимостями
└── README.md        # Документация
```

# Ключевые компоненты

- PasswordManager - основной класс, управляющий хранилищем
- AES-GCM шифрование - защита данных на диске
- Валидация паролей - проверка сложности паролей
- Генератор паролей - криптостойкая генерация

# 🗺️ Планы по развитию

- Поддержка экспорта/импорта в JSON
- Двухфакторная аутентификация
- Облачная синхронизация
- Графический интерфейс (TUI)
- Автоматическое заполнение паролей

