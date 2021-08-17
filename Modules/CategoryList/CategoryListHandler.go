package categorylist

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
)

func HandlecategorylistListing(w http.ResponseWriter, r *http.Request) {
	util.HeaderManager(&w)
	var response util.GeneralResponseModel

	fetchedData := fb.ReadData("/categories")
	if fetchedData == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "Üzgünüz herhangi bir kategori bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		false, "Başarılı", fetchedData,
	}
	w.Write(response.ToJson())
}
