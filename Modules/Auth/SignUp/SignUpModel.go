package SignUp

import "github.com/dgrijalva/jwt-go"

type SignUpModel struct {
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	PersonEmail    string `json:"personEmail"`
	PersonPhone    string `json:"personPhone"`
	Password       string `json:"password"`
}

type Claims struct {
	Username string `json:"personEmail"`
	jwt.StandardClaims
}

type SignUpDbModel struct {
	SignUpDate     string `json:"signUpDate"`
	SignInDate     string `json:"signInDate"`
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	PersonEmail    string `json:"personEmail"`
	PersonPhone    string `json:"personPhone"`
	Password       string `json:"password"`
}

type LoginResponseData struct {
	Token        string `json:"token"`
	Expires      string `json:"expires"`
	UserName     string `json:"name"`
	UserLastName string `json:"lastName"`
}
