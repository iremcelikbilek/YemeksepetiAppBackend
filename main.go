package main

import (
	"fmt"
	"net/http"

	signup "./Modules/Auth/SignUp"
)

func main() {
	fmt.Println("Hello")

	go http.HandleFunc("/signup", signup.HandleSignUp)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
