package main

import (
	"net/http"

	signup "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignUp"
	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
)

func main() {
	go fb.ConnectFirebase()
	go http.HandleFunc("/signup", signup.HandleSignUp)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
