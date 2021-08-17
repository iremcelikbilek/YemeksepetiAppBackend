package search

import (
	"net/http"
	"strings"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	listing "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/RestaurantListing"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
)

func HandleSearchListing(w http.ResponseWriter, r *http.Request) {
	util.HeaderManager(&w)
	var response util.GeneralResponseModel

	cityId, cityIdOk := r.URL.Query()["city"]
	categoryId, categoryIdOk := r.URL.Query()["category"]
	searchKey, searchKeyOk := r.URL.Query()["search"]

	if !searchKeyOk || len(searchKey) < 1 {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "Lütfen arama yapın", nil,
		}
		w.Write(response.ToJson())
		return
	}

	fetchedData := fb.ReadData("/restaurants")

	if fetchedData == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "Üzgünüz herhangi bir restoran bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	var restaurants []listing.RestaurantModel
	mapstructure.Decode(fetchedData, &restaurants)

	var filteredRestaurants []listing.RestaurantModel

	if cityIdOk && len(cityId) == 1 {
		for _, value := range restaurants {
			if value.CityId == cityId[0] {
				filteredRestaurants = append(filteredRestaurants, value)
			}
		}
	}
	if categoryIdOk && len(categoryId) == 1 {
		for _, value := range restaurants {
			if value.CategoryId == categoryId[0] {
				filteredRestaurants = append(filteredRestaurants, value)
			}
		}
	}
	if !cityIdOk && !categoryIdOk {
		filteredRestaurants = restaurants
	}

	var searchedRestaurants []listing.RestaurantModel
	for _, value := range filteredRestaurants {
		if strings.Contains(value.Name, searchKey[0]) {
			filteredRestaurants = append(filteredRestaurants, value)
		}
	}

	response = util.GeneralResponseModel{
		false, "Arama Başarılı", searchedRestaurants,
	}
	w.Write(response.ToJson())
}
