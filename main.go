package main

import (
	"net/http"

	signin "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignIn"
	signup "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignUp"
	cityList "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/CityList"
	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	listing "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/RestaurantListing"
	user "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/User"
)

func main() {
	go fb.ConnectFirebase()
	go http.HandleFunc("/signup", signup.HandleSignUp)
	go http.HandleFunc("/signin", signin.HandleLogin)
	go http.HandleFunc("/restaurantListing", listing.HandleRestaurantListing)
	go http.HandleFunc("/cityList", cityList.HandleCityListing)
	go http.HandleFunc("/user", user.HandleUserData)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
