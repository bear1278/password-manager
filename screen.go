package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
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
