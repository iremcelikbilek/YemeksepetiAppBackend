package restaurantlisting

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
)

var allRestaurants []RestaurantModel

func HandleRestaurantListing(w http.ResponseWriter, r *http.Request) {
	util.HeaderManager(&w)
	var response util.GeneralResponseModel

	cityId, cityIdOk := r.URL.Query()["city"]
	categoryId, categoryIdOk := r.URL.Query()["category"]

	if allRestaurants == nil {
		fetchedData := fb.ReadData("/restaurants")
		if fetchedData == nil {
			w.WriteHeader(http.StatusNotFound)
			response = util.GeneralResponseModel{
				true, "Üzgünüz herhangi bir restoran bulunamadı", nil,
			}
			w.Write(response.ToJson())
			return
		}

		mapstructure.Decode(fetchedData, &allRestaurants)
	}

	var filteredRestaurants []RestaurantModel

	for _, value := range allRestaurants {
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

	if len(filteredRestaurants) == 0 {
		response = util.GeneralResponseModel{
			true, "Filtrenize uygun sonuç bulunamadı, filtrenizi değiştirin.", nil,
		}
	} else {
		response = util.GeneralResponseModel{
			false, "Arama Başarılı", filteredRestaurants,
		}
	}
	w.Write(response.ToJson())
}
