package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"testProject/models"
	u "testProject/utils"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login"} //Список эндпоинтов, для которых не требуется авторизация
		requestPath := r.URL.Path //текущий путь запроса

		//проверяем, не требует ли запрос аутентификации, обслуживаем запрос, если он не нужен
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Получение токена

		if tokenHeader == "" { //Токен отсутствует, возвращаем  403 http-код Unauthorized
			u.Abort(w, http.StatusForbidden, "Missing auth token")
			return
		}

		splitted := strings.Split(tokenHeader, " ") //Токен обычно поставляется в формате `Bearer {token-body}`, мы проверяем, соответствует ли полученный токен этому требованию
		if len(splitted) != 2 {
			u.Abort(w, http.StatusForbidden, "Invalid/Malformed auth token")
			return
		}

		tokenPart := splitted[1] //Получаем вторую часть токена
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Неправильный токен, как правило, возвращает 403 http-код
			u.Abort(w, http.StatusForbidden, "Malformed authentication token")
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			u.Abort(w, http.StatusForbidden, "Token is not valid.")
			return
		}

		//Всё прошло хорошо, продолжаем выполнение запроса
		fmt.Printf("User %d\n", tk.UserId) //Полезно для мониторинга
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //передать управление следующему обработчику!
	})
}
