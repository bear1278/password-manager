package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	pm := NewPasswordManager("passwords.txt")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter master password: ")
	password, err := readPassword()
	err = pm.SetMasterPassword(password)
	if err != nil {
		showError(err.Error())
		return
	}
	showSuccess("Password saved successfully")
	err = pm.LoadFromFile()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		showError(err.Error())
		return
	}
	for {
		ShowMainMenu()
		fmt.Println("Enter choice: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			showError(err.Error())
			return
		}
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			showError(err.Error())
			return
		}
		switch choice {
		case 1:
			err = HandlePasswordGeneration(pm)
		case 2:
			err = HandlePasswordAdd(pm)
		case 3:
			err = HandlePasswordSearch(pm)
		case 4:
			ShowPasswordList(pm)
		case 5:
			err = HandlePasswordUpdate(pm)
		case 6:
			err = HandlePasswordDelete(pm)
		case 7:
			ShowCategories(pm)
		case 8:
			ShowStats(pm)
		case 9:
			err = HandleFindingDuplicates(pm)
		case 0:
			err = HandleExitAndSave(pm)
			return
		default:
			clearScreen()
			showError("Invalid choice")
			waitForEnter()
		}
	}
}
