package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"os"
	"testProject/controllers"
	"testProject/jwt"
)

func main() {

	router := mux.NewRouter()
	router.Use(auth.JwtAuthentication) // добавляем middleware проверки JWT-токена
	router.HandleFunc("/api/user/new",
		controllers.CreateAccount).Methods("POST")

	router.HandleFunc("/api/user/login",
		controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/notes/{id:[0-9]+}",
		controllers.NoteDetailHandler).Methods("GET")

	router.HandleFunc("/api/notes",
		controllers.NoteCreateHandler).Methods("POST")

	router.HandleFunc("/api/notes",
		controllers.NoteListHandler).Methods("GET")

	router.HandleFunc("/api/notes/{id:[0-9]+}",
		controllers.NoteUpdateHandler).Methods("PUT")

	router.HandleFunc("/api/notes/{id:[0-9]+}",
		controllers.NoteDeleteHandler).Methods("DELETE")

	port := os.Getenv("listen_port") //Получить порт из файла .env; мы не указали порт, поэтому при локальном тестировании должна возвращаться пустая строка
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":" + port, router) //Запустите приложение, посетите localhost:8000/api

	if err != nil {
		fmt.Print(err)
	}
}
