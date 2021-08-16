package main

import (
	"net/http"

	signin "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignIn"
	signup "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignUp"
	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
)

func main() {
	go fb.ConnectFirebase()
	go http.HandleFunc("/signup", signup.HandleSignUp)
	go http.HandleFunc("/signin", signin.HandleLogin)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
