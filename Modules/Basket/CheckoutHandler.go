package basket

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
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

	fetchedData := fb.ReadData("/basket/" + util.MailToPath(userMail))
	if fetchedData == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "Sepetiniz boş", nil,
		}
		w.Write(response.ToJson())
		return
	}

	oldHistory := fb.ReadData("/history/" + util.MailToPath(userMail))
	if oldHistory == nil {
		error := fb.WriteData("/history/"+util.MailToPath(userMail), fetchedData)
		if error != nil {
			response = util.GeneralResponseModel{
				true, "Sepetten kaldırma başarısız", nil,
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response.ToJson())
			return
		}
	} else {
		var historyItems []BasketModel
		var basketItems []BasketModel
		mapstructure.Decode(oldHistory, &historyItems)
		mapstructure.Decode(fetchedData, &basketItems)

		for _, value := range basketItems {
			historyItems = append(historyItems, value)
		}

		error := fb.WriteData("/history/"+util.MailToPath(userMail), historyItems)
		if error != nil {
			response = util.GeneralResponseModel{
				true, "Sepetten kaldırma başarısız", nil,
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response.ToJson())
			return
		}
	}

	_ = fb.DeletePath("/basket/" + util.MailToPath(userMail))

	response = util.GeneralResponseModel{
		false, "Başarılı", nil,
	}
	w.Write(response.ToJson())
}
