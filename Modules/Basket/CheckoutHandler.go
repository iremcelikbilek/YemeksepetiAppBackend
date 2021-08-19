package basket

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
)

func HandleCheckout(w http.ResponseWriter, r *http.Request) {
	util.HeaderManager(&w)
	var response util.GeneralResponseModel

	var userMail string
	if isSucessToken, message := util.CheckToken(r); !isSucessToken {
		response = util.GeneralResponseModel{
			true, message, nil,
		}
		w.Write(response.ToJson())
		return
	} else {
		userMail = message
	}

	fetchedData := fb.ReadData("/basket/" + userMail)
	if fetchedData == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "Sepetiniz boş", nil,
		}
		w.Write(response.ToJson())
		return
	}

	error := fb.PushData("/history/"+userMail, fetchedData)
	if error != nil {
		response = util.GeneralResponseModel{
			true, "Sepetten kaldırma başarısız", nil,
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response.ToJson())
		return
	}

	_ = fb.DeletePath("/basket/" + userMail)

	response = util.GeneralResponseModel{
		false, "Başarılı", nil,
	}
	w.Write(response.ToJson())
}
