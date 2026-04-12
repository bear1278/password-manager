package main

import (
	"bufio"
	"fmt"
	"os"
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
