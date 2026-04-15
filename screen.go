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
	fmt.Println(
		`1. Generate new password
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
	defer waitForEnter()
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
	return nil
}
func HandlePasswordAdd(pm *PasswordManager) error {
	var err error
	clearScreen()
	defer waitForEnter()
	fmt.Println("=== Add New Password ===")
	name := ReadUserInput("Enter service name: ")
	password := ReadUserInput("Enter password (or press Enter to generate): ")
	if password == "" {
		for err = pm.CheckPasswordStrength(password); err != nil; err = pm.CheckPasswordStrength(password) {
			password, err = pm.GeneratePassword(12)
			if err != nil {
				showError(err.Error())
				return err
			}
		}
		showInfo("Password generated successfully: " + password)
	}
	category := ReadUserInput("Enter category: ")
	err = pm.SavePassword(name, password, category)
	if err != nil {
		showError(err.Error())
		return err
	}
	showSuccess("Password saved successfully")
	return nil
}
func HandlePasswordSearch(pm *PasswordManager) error {
	clearScreen()
	defer waitForEnter()
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
	showSuccess("Password search successfully")
	ShowPasswordDetails(password)
	return nil
}
func HandlePasswordUpdate(pm *PasswordManager) error {
	var err error
	clearScreen()
	defer waitForEnter()
	fmt.Println("=== Password Update ===")
	name := ReadUserInput("Enter service name: ")
	_, ok := pm.passwords[name]
	if !ok {
		showError("password not found")
		return errors.New("password not found")
	}
	passwordValue := ReadUserInput("Enter password (or press Enter to generate): ")
	if passwordValue == "" {
		for err = pm.CheckPasswordStrength(passwordValue); err != nil; err = pm.CheckPasswordStrength(passwordValue) {
			passwordValue, err = pm.GeneratePassword(12)
			if err != nil {
				showError(err.Error())
				return err
			}
		}
		showInfo("Password generated successfully: " + passwordValue)
	}
	err = pm.UpdatePassword(name, passwordValue)
	if err != nil {
		showError(err.Error())
		return err
	}
	showSuccess("Password updated successfully")
	return nil
}

func HandleExitAndSave(pm *PasswordManager) error {
	clearScreen()
	fmt.Println("Saving changes...")
	err := pm.SaveToFile()
	if err != nil {
		err = fmt.Errorf("error saving data: %w", err)
		showError(err.Error())
		return err
	}
	showSuccess("Password saved successfully")
	showSuccess("Goodbye!")
	return nil
}

func HandleFindingDuplicates(pm *PasswordManager) error {
	clearScreen()
	duplicates := pm.FindDuplicatePasswords()
	if len(duplicates) == 0 {
		showInfo("No duplicate passwords found")
	} else {
		for k, v := range duplicates {
			fmt.Print("password: ", k, "for ")
			for _, vv := range v {
				fmt.Print(vv, ", ")
			}
			fmt.Println()
		}
	}
	waitForEnter()
	return nil
}

func ShowStats(pm *PasswordManager) {
	clearScreen()
	stats := pm.GetPasswordStats()
	fmt.Println("total password", stats["total"].(int))
	categories := pm.ListCategories()
	if len(categories) != 0 {
		fmt.Println("Count of passwords by category:")
		for _, category := range categories {
			fmt.Println(category, stats[category].(int))
		}
	}
	fmt.Println("oldest password", stats["oldest"].(string))
	fmt.Println("newest password", stats["newest"].(string))
	waitForEnter()
}

func ShowPasswordList(pm *PasswordManager) {
	clearScreen()
	PrintPasswordList(pm.ListPasswords())
	waitForEnter()
}

func HandlePasswordDelete(pm *PasswordManager) error {
	clearScreen()
	input := ReadUserInput("Enter service name: ")
	err := pm.DeletePassword(input)
	if err != nil {
		showError(err.Error())
		return err
	} else {
		showSuccess("Password deleted successfully")
	}
	waitForEnter()
	return nil
}

func ShowCategories(pm *PasswordManager) {
	clearScreen()
	categories := pm.ListCategories()
	if len(categories) == 0 {
		showSuccess("Empty list of categories")
	} else {
		for _, category := range categories {
			fmt.Println(category)
		}
	}
	waitForEnter()
}
