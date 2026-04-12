package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
	"text/tabwriter"
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tCategory\tCreated\tLast Modified")
	fmt.Println(strings.Repeat("-", 70))
	for _, password := range passwords {
		fmt.Fprintln(w, password.Name+"\t"+password.Category+"\t"+password.CreatedAt.Format(time.DateTime)+"\t"+password.LastModified.Format(time.DateTime))
	}
	w.Flush()
}

func ShowPasswordDetails(password Password) {
	fmt.Println("=== Password details ===")
	fmt.Println("Service", password.Name)
	fmt.Println("Category", password.Category)
	fmt.Println("Password", password.Value)
	fmt.Println("Created", password.CreatedAt.Format(time.DateTime))
	fmt.Println("Last Modified", password.LastModified.Format(time.DateTime))
}
