package user

import (
	"net/http"

	signup "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignUp"
	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
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
	var userDbData signup.SignUpDbModel
	mapstructure.Decode(fetchedData, &userDbData)

	response = util.GeneralResponseModel{
		false, "Başarılı", userDbData,
	}

	w.Write(response.ToJson())
}
