package users

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go_banking/helpers"
	"go_banking/interfaces"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Login(username string, pass string) map[string]interface{} {
	// создаем коннект к базе
	db := helpers.ConnectDB()
	user := &interfaces.User{}
	if db.Where("username=?", username).First(&user).RecordNotFound() {
		return map[string]interface{}{"Message": "User not found"}
	}
	// проверяем пароль
	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if passError == bcrypt.ErrMismatchedHashAndPassword && passError != nil {
		return map[string]interface{}{"Message": "Wrong password"}
	}

	// ищем аккаунт
	accounts := []interfaces.ResponseAccount{}
	fmt.Println("accounts --->", accounts)
	db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

	// готовим ответ
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	defer db.Close()

	// создаем токен

	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute ^ 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	// предварительный ответ
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseUser
	fmt.Println("response --->", response["data"])
	return response
}
