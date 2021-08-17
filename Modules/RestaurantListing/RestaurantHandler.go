package restaurantlisting

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
)

func HandleRestaurantListing(w http.ResponseWriter, r *http.Request) {
	util.HeaderManager(&w)
	var response util.GeneralResponseModel

	cityId, cityIdOk := r.URL.Query()["city"]
	categoryId, categoryIdOk := r.URL.Query()["category"]

	fetchedData := fb.ReadData("/restaurants")

	if fetchedData == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "Üzgünüz herhangi bir restoran bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	var restaurants []RestaurantModel
	mapstructure.Decode(fetchedData, &restaurants)

	var filteredRestaurants []RestaurantModel

	if cityIdOk || len(cityId) == 1 {
		for _, value := range restaurants {
			if value.City_Id == cityId[0] {
				filteredRestaurants = append(filteredRestaurants, value)
			}
		}
	} else if categoryIdOk || len(categoryId) == 1 {
		for _, value := range restaurants {
			if value.CategoryId == categoryId[0] {
				filteredRestaurants = append(filteredRestaurants, value)
			}
		}
	}

	response = util.GeneralResponseModel{
		false, "Listeleme Başarılı", filteredRestaurants,
	}
	w.Write(response.ToJson())
}
