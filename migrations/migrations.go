package migrations

import (
	"fmt"
	"go_banking/helpers"
	"go_banking/interfaces"
)

func createAccounts() {
	db := helpers.ConnectDB()

	users := &[2]interfaces.User{
		{Username: "Martin", Email: "martin@gmail.com"},
		{Username: "Michael", Email: "michael@gmail.com"},
	}
	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := &interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		db.Create(&user)

		account := &interfaces.Account{Type: "Daily Account", Name: string(users[i].Username + "'s" + " account"),
			Balance: uint(1000 * int(i+1)), UserID: user.ID}
		db.Create(&account)
	}
	defer db.Close()
}

func Migrate() {
	User := &interfaces.User{}
	fmt.Println(User)
	Account := &interfaces.Account{}
	db := helpers.ConnectDB()
	// создаем таблицы
	db.AutoMigrate(&User, &Account)
	defer db.Close()
	// создаем записи
	createAccounts()
}

func MigrateTransactions() {
	Transaction := &interfaces.Transaction{}
	db := helpers.ConnectDB()
	// создаем таблицы
	db.AutoMigrate(&Transaction)
	defer db.Close()
}
