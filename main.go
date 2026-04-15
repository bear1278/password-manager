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
			clearScreen()
			PrintPasswordList(pm.ListPasswords())
			waitForEnter()
		case 5:
			err = HandlePasswordUpdate(pm)
		case 6:
			clearScreen()
			input = ReadUserInput("Enter service name: ")
			err = pm.DeletePassword(input)
			if err != nil {
				showError(err.Error())
			} else {
				showSuccess("Password deleted successfully")
			}
			waitForEnter()
		case 7:
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
		case 8:
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
		case 9:
			clearScreen()
			dublicates := pm.FindDuplicatePasswords()
			if len(dublicates) == 0 {
				fmt.Println("No duplicate passwords found")
			} else {
				for k, v := range dublicates {
					fmt.Print("password: ", k, "for ")
					for _, vv := range v {
						fmt.Print(vv, ", ")
					}
					fmt.Println()
				}
			}
			waitForEnter()
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
