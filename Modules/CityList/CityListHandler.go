package citylist

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
)

func HandleCityListing(w http.ResponseWriter, r *http.Request) {
	util.HeaderManager(&w)
	var response util.GeneralResponseModel

	fetchedData := fb.ReadData("/cities")
	if fetchedData == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "Üzgünüz herhangi bir il bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	var cities []CityModel
	mapstructure.Decode(fetchedData, &cities)

	response = util.GeneralResponseModel{
		false, "Başarılı", fetchedData,
	}
	w.Write(response.ToJson())
}
