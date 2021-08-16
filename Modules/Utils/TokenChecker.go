package Utils

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Mail string `json:"personEmail"`
	jwt.StandardClaims
}

var JWT_Token = []byte("PETNER_JWT_TOKEN")

func CheckToken(r *http.Request) (bool, string) {
	tokenString := r.Header.Get("token")
	if tokenString == "" {
		return false, "Token boş ya da bulunamadı"
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_Token, nil
	})

	if err != nil {
		return false, err.Error()
	}
	if !token.Valid {
		return false, "Token geçerli değil"
	}
	return true, claims.Mail
}
