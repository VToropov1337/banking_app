package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go_banking/helpers"
	"go_banking/interfaces"
	"go_banking/users"
	"io/ioutil"
	"log"
	"net/http"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username string
	Email    string
	Password string
}

func readBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	return body
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	if call["message"] == "all is fine" {
		resp := call
		// возвращаем json
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := interfaces.ErrResponse{Message: "Wrong username or password"}
		// возвращаем json
		json.NewEncoder(w).Encode(resp)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	// читаем боди
	body := readBody(r)

	// обрабатываем логин
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	fmt.Println("api.formattedBody--->", formattedBody)
	login := users.Login(formattedBody.Username, formattedBody.Password)
	fmt.Println("login===>", login)

	// подготавливаем ответ
	apiResponse(login, w)

}

func register(w http.ResponseWriter, r *http.Request) {
	// читаем боди
	body := readBody(r)

	// обрабатываем логин
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	fmt.Println("api.register.formattedBody--->", formattedBody)
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
	fmt.Println("register===>", register)

	// подготавливаем ответ
	apiResponse(register, w)

}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	auth := r.Header.Get("Authorization")

	user := users.GetUser(userId, auth)
	apiResponse(user, w)
}

func StartApi() {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	fmt.Println("App is working on port :8888")
	log.Fatal(http.ListenAndServe(":8888", router))

}
