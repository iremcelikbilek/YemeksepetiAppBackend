package search

import (
	"encoding/json"
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
	var searchData SearchModel

	cityId, cityIdOk := r.URL.Query()["city"]
	categoryId, categoryIdOk := r.URL.Query()["category"]

	err := json.NewDecoder(r.Body).Decode(&searchData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		if strings.Contains(value.Name, searchData.Search) {
			searchedRestaurants = append(searchedRestaurants, value)
		}
	}

	response = util.GeneralResponseModel{
		false, "Arama Başarılı", searchedRestaurants,
	}
	w.Write(response.ToJson())
}

type SearchModel struct {
	Search string `json:"search"`
}
