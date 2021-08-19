package main

import (
	"net/http"

	signin "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignIn"
	signup "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignUp"
	basket "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Basket"
	categoryList "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/CategoryList"
	cityList "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/CityList"
	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	listing "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/RestaurantListing"
	search "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Search"
	user "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/User"
)

func main() {
	go fb.ConnectFirebase()
	go http.HandleFunc("/signup", signup.HandleSignUp)
	go http.HandleFunc("/signin", signin.HandleLogin)
	go http.HandleFunc("/restaurantListing", listing.HandleRestaurantListing)
	go http.HandleFunc("/cityList", cityList.HandleCityListing)
	go http.HandleFunc("/user", user.HandleUserData)
	go http.HandleFunc("/categoryList", categoryList.HandleCategoryListing)
	go http.HandleFunc("/search", search.HandleSearchListing)
	go http.HandleFunc("/basket", basket.HandleBasket)
	go http.HandleFunc("/addToBasket", basket.HandleAddToBasket)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
