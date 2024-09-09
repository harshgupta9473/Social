package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type handlerFunctionToken func(http.ResponseWriter,*http.Request,*jwt.Token)

func WithJWTAuth(handlerFunc handlerFunctionToken)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		fmt.Println("calling JWT auth middleware")
		cookie,err:=r.Cookie("authToken")
		if err!=nil{
			// no cookie found
			// http.Redirect(w,r)
			http.Error(w,"cookie not present",http.StatusUnauthorized)
			return
		}
		tokenString:=cookie.Value
		token,err:=ValidateJWT(tokenString)
		if err!=nil || !token.Valid{
			http.Error(w,"not valid",http.StatusUnauthorized)
			return
		}
		if claims,ok:=token.Claims.(jwt.MapClaims); ok && token.Valid{
			exp:=int64(claims["exp"].(float64))
			if time.Now().Unix()>exp{
				http.Error(w,"not valid",http.StatusUnauthorized)
			return
			}
			handlerFunc(w,r,token)
		} else{
			http.Error(w,"not valid",http.StatusUnauthorized)
			return
		}
	}
}

func CreateJWT(email string, createdAt time.Time,id uint64) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	key := os.Getenv("secretKey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":     email,
		"id":       id,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
		"created_at": createdAt,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Println("Error Creating the token", err)
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	key := os.Getenv("secretKey")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
}
