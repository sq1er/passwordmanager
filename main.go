package main

import (
	"PasswordManager/account"
	"PasswordManager/encrypter"
	"PasswordManager/files"
	"PasswordManager/output"
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var menu = map[string]func(*account.VaultWithDb){
	"1": outputAccounts,
	"2": createAccount,
	"3": findAccountByUrl,
	"4": findAccountByLogin,
	"5": deleteAccount,
}

var menuVariants = []string{
	"1. Вывести все аккаунты",
	"2. Добавить аккаунт",
	"3. Найти аккаунт по URL",
	"4. Найти аккаунт по логину",
	"5. Удалить аккаунт",
	"6. Выйти",
	"Выберите вариант",
}

func main() {
	err := godotenv.Load()
	if err != nil {
		output.PrintError("Не удалось найти env файл")
	}
	vault := account.NewVault(files.NewJsonDb("data.vault"), *encrypter.NewEncrypter())
	// vault := account.NewVault(cloud.NewCloudDb("https://a.ru"))
	color.Magenta("___Менеджер паролей___")
Menu:
	for {
		variant := promptData(menuVariants...)
		menuFunc := menu[variant]
		if menuFunc == nil {
			break Menu
		}
		menuFunc(vault)
	}
}

func outputAccounts(vault *account.VaultWithDb) {
	fmt.Println()
	if len(vault.Accounts) == 0 {
		output.PrintError("Аккаунтов нет")
		fmt.Println()
		return
	}
	for _, account := range vault.Accounts {
		account.OutputAccount()
	}
}

func findAccountByUrl(vault *account.VaultWithDb) {
	inputUrl := promptData("Введите URL для поиска")
	fmt.Println()
	_, err := url.ParseRequestURI(inputUrl)
	if err != nil {
		output.PrintError("Неверный формат URL")
		return
	}
	accounts := vault.FindAccount(inputUrl, func(acc account.Account, str string) bool {
		return strings.EqualFold(acc.Url, str)
	})
	outputResult(&accounts)
}

func findAccountByLogin(vault *account.VaultWithDb) {
	login := promptData("Введите логин для поиска")
	fmt.Println()
	accounts := vault.FindAccount(login, func(acc account.Account, str string) bool {
		return strings.EqualFold(acc.Login, str)
	})
	outputResult(&accounts)
}

func outputResult(accounts *[]account.Account) {
	if len(*accounts) == 0 {
		color.Red("Аккаунт не найден")
		return
	}
	for _, account := range *accounts {
		color.Green("Найден аккаунт:")
		account.OutputAccount()
	}
}

func deleteAccount(vault *account.VaultWithDb) {
	inputUrl := promptData("Введите URL аккаунта для удаления")
	_, err := url.ParseRequestURI(inputUrl)
	if err != nil {
		output.PrintError("Неверный формат URL")
		fmt.Println()
		return
	}
	isDeleted := vault.DeleteAccountByUrl(inputUrl)
	if isDeleted {
		color.Green("Удалено")
		fmt.Println()
	} else {
		output.PrintError("Не найдено")
		fmt.Println()
	}
}

func createAccount(vault *account.VaultWithDb) {
	login := promptData("Введите логин")
	n := promptData("Нажмите 1 чтобы ввести пароль,\nenter чтобы сгенерировать пароль")
	password := ""
	if n == "1" {
		password = promptData("Введите пароль от аккаунта")
	}
	inputUrl := promptData("Введите URL аккаунта")
	myAccount, err := account.NewAccount(login, password, inputUrl)
	if err != nil {
		output.PrintError("Неверный формат URL или логин")
		return
	}
	vault.AddAccount(*myAccount)
	color.Green("Аккаунт добавлен")
	fmt.Println()
}

func promptData(prompt ...string) string {
	for i, line := range prompt {
		if i == len(prompt)-1 {
			fmt.Printf("%v: ", line)
		} else {
			fmt.Println(line)
		}
	}
	var res string
	fmt.Scanln(&res)
	return res
}
