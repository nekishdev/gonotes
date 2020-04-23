package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

/*
Структура прав доступа JWT
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//структура для учётной записи пользователя
type Account struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

//Проверить входящие данные пользователя ...
func (account *Account) Validate() error {

	if !strings.Contains(account.Email, "@") {
		return errors.New("Email address is required")
	}

	if len(account.Password) < 6 {
		return errors.New("Password is required")
	}

	//Email должен быть уникальным
	temp := &Account{}

	//проверка на наличие ошибок и дубликатов электронных писем
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("Connection error. Please retry")
	}
	if temp.Email != "" {
		return errors.New("Email address already in use by another user.")
	}

	return nil
}

func (account *Account) Create() (*Account, error) {

	if err := account.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return nil, errors.New("Failed to create account, connection error.")
	}

	//Создать новый токен JWT для новой зарегистрированной учётной записи
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString
	account.Password = "" //удалить пароль

	return account, nil
}

func Login(email, password string) (*Account, error) {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New( "Email address not found")
		}
		return nil, errors.New("Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Пароль не совпадает!!
		return nil, errors.New("Invalid login credentials. Please try again")
	}
	//Работает! Войти в систему
	account.Password = ""

	//Создать токен JWT
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString // Сохраните токен в ответе

	return account, nil
}

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //Пользователь не найден!
		return nil
	}

	acc.Password = ""
	return acc
}