package Login

import "github.com/dgrijalva/jwt-go"

type LoginModel struct {
	PersonEmail string `json:"personEmail"`
	Password    string `json:"password"`
}

type Claims struct {
	Username string `json:"personEmail"`
	jwt.StandardClaims
}

type LoginResponseData struct {
	Token        string `json:"token"`
	Expires      string `json:"expires"`
	UserName     string `json:"name"`
	UserLastName string `json:"lastName"`
}
