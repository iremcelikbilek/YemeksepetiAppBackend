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

	var fetchedData interface{}

	if cityIdOk || len(cityId) == 1 {
		fetchedData = fb.GetFilteredData("/restaurants", "city_id", cityId[0])
	} else if categoryIdOk || len(categoryId) == 1 {
		fetchedData = fb.GetFilteredData("/restaurants", "category_id", categoryId[0])
	} else {
		fetchedData = fb.ReadData("/restaurants")
	}

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

	response = util.GeneralResponseModel{
		false, "Listeleme Başarılı", fetchedData,
	}
	w.Write(response.ToJson())
}
