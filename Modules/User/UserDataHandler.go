package user

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
)

var JWT_Token = []byte("YEMEK_SEPETI_JWT_TOKEN")

func HandleUserData(w http.ResponseWriter, r *http.Request) {
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

	fetchedData := fb.GetFilteredData("/persons", "personEmail", userMail)

	response = util.GeneralResponseModel{
		false, "Başarılı", fetchedData,
	}

	w.Write(response.ToJson())
}
