package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func clearScreen() { // Очистка экрана
	fmt.Print("\033[H\033[2J")
}
func showSuccess(message string) { // Вывод сообщения об успехе
	fmt.Println(colorGreen, "✓ Success: ", colorReset, message)
}
func showError(message string) { // Вывод сообщения об ошибке
	fmt.Println(colorRed, "✗ Error: ", colorReset, message)
}
func showInfo(message string) { // Вывод информационного сообщения
	fmt.Println(colorYellow, "→ Info: ", colorReset, message)
}
func waitForEnter() { // Ожидание нажатия Enter
	fmt.Println("Press Enter to continue...")
	_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	clearScreen()
}

func ReadUserInput(prompt string) string {
	fmt.Println(prompt)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(input)
}

func readPassword() (string, error) {
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func ShowMainMenu() {
	clearScreen()
	fmt.Println(strings.Repeat("=", 42))
	fmt.Println("Password Manager")
	fmt.Println(strings.Repeat("=", 42))
	fmt.Println(`1. Generate new password
	2. Add new password
	3. Get password
	4. List all passwords
	5. Update password
	6. Delete password
	7. List categories
	8. Show password statistics
	9. Find duplicate passwords
	0. Exit`)
	fmt.Println(strings.Repeat("=", 42))
}

func PrintPasswordList(passwords []Password) {
	fmt.Println("=== Password list ===")
	fmt.Printf("%-25s %-25s %-25s %-25s\n", "Name", "Category", "Created", "Last Modified")
	fmt.Println(strings.Repeat("-", 80) + "\n")
	for _, password := range passwords {
		fmt.Printf("%-25s %-25s %-25s %-25s\n", password.Name, password.Category, password.CreatedAt.Format(time.DateTime), password.LastModified.Format(time.DateTime))
	}
}

func ShowPasswordDetails(password Password) {
	fmt.Println("=== Password details ===")
	fmt.Println("Service", password.Name)
	fmt.Println("Category", password.Category)
	fmt.Println("Password", password.Value)
	fmt.Println("Created", password.CreatedAt.Format(time.DateTime))
	fmt.Println("Last Modified", password.LastModified.Format(time.DateTime))
}

func HandlePasswordGeneration(pm *PasswordManager) error {
	clearScreen()
	fmt.Println("=== Password Generation ===")
	fmt.Println("Enter password length (min 8): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		showError(err.Error())
		return err
	}
	input = strings.TrimSpace(input)
	length, err := strconv.Atoi(input)
	if err != nil {
		showError(err.Error())
		return err
	}
	password, err := pm.GeneratePassword(length)
	if err != nil {
		showError(err.Error())
		return err
	}
	showSuccess("Password generated successfully")
	fmt.Println("Generated password: ", password)
	waitForEnter()
	return nil
}
func HandlePasswordAdd(pm *PasswordManager) error {
	clearScreen()
	fmt.Println("=== Add New Password ===")
	fmt.Println("Enter service name: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		showError(err.Error())
		return err
	}
	name := strings.TrimSpace(input)
	fmt.Println("Enter password (or press Enter to generate): ")
	input, err = reader.ReadString('\n')
	if err != nil {
		showError(err.Error())
		return err
	}
	password := strings.TrimSpace(input)
	if password == "" {
		password, err = pm.GeneratePassword(8)
		if err != nil {
			showError(err.Error())
			return err
		}
		showInfo("Password generated successfully: " + password)
	}
	fmt.Println("Enter category: ")
	input, err = reader.ReadString('\n')
	if err != nil {
		showError(err.Error())
		return err
	}
	category := strings.TrimSpace(input)
	err = pm.SavePassword(name, category, password)
	if err != nil {
		showError(err.Error())
		return err
	}
	showSuccess("Password saved successfully")
	waitForEnter()
	return nil
}
func HandlePasswordSearch(pm *PasswordManager) error {
	clearScreen()
	fmt.Println("=== Search Password ===")
	fmt.Println("Enter service name: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		showError(err.Error())
		return err
	}
	name := strings.TrimSpace(input)
	password, ok := pm.passwords[name]
	if !ok {
		showError("password not found")
		return errors.New("password not found")
	}
	fmt.Println("Password Details:")
	fmt.Println("Service", password.Name)
	fmt.Println("Category", password.Category)
	fmt.Println("Password", password.Value)
	fmt.Println("Created", password.CreatedAt.Format(time.DateTime))
	fmt.Println("Last Modified", password.LastModified.Format(time.DateTime))
	waitForEnter()
	return nil
}
func HandlePasswordUpdate(pm *PasswordManager) error {
	clearScreen()
	fmt.Println("=== Password Update ===")
	fmt.Println("Enter service name: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		showError(err.Error())
		return err
	}
	name := strings.TrimSpace(input)
	_, ok := pm.passwords[name]
	if !ok {
		showError("password not found")
		return errors.New("password not found")
	}
	fmt.Println("Enter password (or press Enter to generate): ")
	input, err = reader.ReadString('\n')
	if err != nil {
		showError(err.Error())
		return err
	}
	passwordValue := strings.TrimSpace(input)
	if passwordValue == "" {
		passwordValue, err = pm.GeneratePassword(8)
		if err != nil {
			showError(err.Error())
			return err
		}
		showInfo("Password generated successfully: " + passwordValue)
	}
	err = pm.UpdatePassword(name, passwordValue)
	if err != nil {
		showError(err.Error())
		return err
	}
	showSuccess("Password updated successfully")
	waitForEnter()
	return nil
}
