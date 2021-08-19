package basket

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
)

func HandleRemoveToBasket(w http.ResponseWriter, r *http.Request) {
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

	restaurantId, restaurantIdOk := r.URL.Query()["restaurant"]
	menuId, menuIdOk := r.URL.Query()["menu"]

	if !menuIdOk || !restaurantIdOk || len(restaurantId) != 1 || len(menuId) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, "Gerekli bilgileri göndermediniz", nil,
		}
		w.Write(response.ToJson())
		return
	}

	var newBasketItems []BasketModel
	basketData := fb.ReadData("/basket/" + util.MailToPath(userMail))
	if basketData == nil {
		response = util.GeneralResponseModel{
			true, "Sepete zaten boş", nil,
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.ToJson())
		return
	} else {
		var basketItems []BasketModel
		mapstructure.Decode(basketData, &basketItems)

		for _, value := range basketItems {
			if value.Id != restaurantId[0] && value.Menu.Id != menuId[0] {
				newBasketItems = append(newBasketItems, value)
			}
		}

		error := fb.WriteData("/basket/"+util.MailToPath(userMail), newBasketItems)
		if error != nil {
			response = util.GeneralResponseModel{
				true, "Sepetten kaldırma başarısız", nil,
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response.ToJson())
			return
		}
	}

	response = util.GeneralResponseModel{
		false, "Başarılı", newBasketItems,
	}
	w.Write(response.ToJson())
}
