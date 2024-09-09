package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func WriteJSON(w http.ResponseWriter,statusCode int ,response any){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func ExtractEmailFromJWT(token *jwt.Token)(string,error){
	claims:=token.Claims.(jwt.MapClaims)
	email,ok:=claims["user"].(string)
	if !ok{
		return "",fmt.Errorf("error")
	}
	return email,nil
}
func ExtractIDFromJWT(token *jwt.Token)(uint64,error){
	claims:=token.Claims.(jwt.MapClaims)
	id,ok:=claims["id"].(float64)
	if !ok{
		return 0,fmt.Errorf("error")
	}
	return  uint64(id),nil
}

func SetTokenInCookies(w http.ResponseWriter,jwtTokenoken string){
	cookie:=http.Cookie{
		Name:"authToken",
		Value: jwtTokenoken,
		Expires: time.Now().Add(24*time.Hour),
		HttpOnly: true,
		// Secure: true,
		Path: "/",
	}

	http.SetCookie(w,&cookie)

}