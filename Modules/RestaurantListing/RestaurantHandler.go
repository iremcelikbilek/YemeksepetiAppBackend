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

	response = util.GeneralResponseModel{
		false, "Giriş Başarılı", fetchedData,
	}
	w.Write(response.ToJson())
}
