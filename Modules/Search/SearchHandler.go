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

	for _, value := range restaurants {
		if cityIdOk && len(cityId) == 1 {
			if value.CityId != cityId[0] {
				continue
			}
		}

		if categoryIdOk && len(categoryId) == 1 {
			if value.CategoryId != categoryId[0] {
				continue
			}
		}

		filteredRestaurants = append(filteredRestaurants, value)
	}

	var searchedRestaurants []listing.RestaurantModel
	for _, value := range filteredRestaurants {
		if strings.Contains(strings.ToLower(value.Name), strings.ToLower(searchData.Search)) {
			searchedRestaurants = append(searchedRestaurants, value)
		}
	}

	if len(searchedRestaurants) == 0 {
		response = util.GeneralResponseModel{
			true, "Aramanıza uygun sonuç bulunamadı, filtrenizi ya da aramanızı değiştirin.", nil,
		}
	} else {
		response = util.GeneralResponseModel{
			false, "Arama Başarılı", searchedRestaurants,
		}
	}

	w.Write(response.ToJson())
}

type SearchModel struct {
	Search string `json:"search"`
}
