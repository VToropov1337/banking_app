package users

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go_banking/helpers"
	"go_banking/interfaces"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func prepareToken(user *interfaces.User) string {
	// создаем токен
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute ^ 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount) map[string]interface{} {
	// готовим ответ
	var responseUser = &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}
	fmt.Println("users.responseUser ---------->", responseUser)
	// предварительный ответ
	var token = prepareToken(user)
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	//responseUser.ID = 42
	response["data"] = responseUser
	fmt.Println("users.response --->", response["data"])
	return response
}

func Login(username string, pass string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// создаем коннект к базе
		db := helpers.ConnectDB()
		//user := &interfaces.User{}
		user := interfaces.User{}
		fmt.Println("users.user==>", user)
		if db.Where("username=?", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"Message": "User not found"}
		}
		fmt.Println("users.user======>", user)
		// проверяем пароль
		passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if passError == bcrypt.ErrMismatchedHashAndPassword && passError != nil {
			return map[string]interface{}{"Message": "Wrong password"}
		}

		// ищем аккаунт
		var accounts []interfaces.ResponseAccount
		fmt.Println("users.accounts --->", accounts)
		//db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)
		db.Table("accounts").Select("id,name,balance").Where("user_id in (?)", []int{1, 3}).Scan(&accounts)
		fmt.Println("users.accounts ------->", accounts)

		defer db.Close()

		//var response = prepareResponse(user, accounts)
		var response = prepareResponse(&user, accounts)
		return response

	} else {
		return map[string]interface{}{"message":"not valid values"}
	}
}

func Register(username string, email string, pass string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	fmt.Println("users.register.valid--->", valid)
	if valid {
		db := helpers.ConnectDB()
		generatedPassword := helpers.HashAndSalt([]byte(pass))
		user := &interfaces.User{Username: username, Email: email, Password: generatedPassword}
		db.Create(user)

		account := &interfaces.Account{Type: "Daily Account", Name: username + "'s" + " account",
			Balance: 0, UserID: user.ID}
		db.Create(account)

		defer db.Close()

		var accounts []interfaces.ResponseAccount
		respAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: int(account.Balance)}
		accounts = append(accounts, respAccount)
		var response = prepareResponse(user, accounts)

		return response


	} else {
		return map[string]interface{}{"message":"not valid values"}
	}

}
