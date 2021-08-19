package basket

import (
	listing "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/RestaurantListing"
)

type BasketModel struct {
	Menu []listing.MunuModel `json:"menu"`
	Name string              `json:"name"`
	Id   string              `json:"id"`
}
