package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testProject/models"
	u "testProject/utils"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
	if err != nil {
		u.Abort(w, 400, fmt.Sprintf("Invalid request. Field to parse request body. %s", err))
		return
	}

	acc, err := account.Create() //Создать аккаунт
	if err != nil {
		u.Abort(w, 400, err)
		return
	}
	u.Respond(w, acc)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
	if err != nil {
		u.Abort(w, 400, fmt.Sprintf("Invalid request. Field to parse request body. %s", err))
		return
	}

	acc, err := models.Login(account.Email, account.Password)
	if err != nil {
		u.Abort(w, 400, err)
		return
	}
	u.Respond(w, acc)
}
