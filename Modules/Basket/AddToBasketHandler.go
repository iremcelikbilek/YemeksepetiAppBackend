package basket

import (
	"net/http"

	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	listing "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/RestaurantListing"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
)

func HandleAddToBasket(w http.ResponseWriter, r *http.Request) {
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

	restaurantsData := fb.ReadData("/restaurants")
	var restaurants []listing.RestaurantModel
	mapstructure.Decode(restaurantsData, &restaurants)
	var restaurantName string
	var menu listing.MunuModel

	for _, value := range restaurants {
		if value.Id == restaurantId[0] {
			restaurantName = value.Name
			for _, menu := range value.Menu {
				if menu.Id == menuId[0] {
					menu = menu
					break
				}
			}
			break
		}
	}

	var newItem = BasketModel{menu, restaurantName, restaurantId[0]}

	basketData := fb.ReadData("/basket/" + userMail)
	if basketData == nil {
		error := fb.PushData("/basket/"+userMail, [1]listing.MunuModel{newItem})
		if error != nil {
			response = util.GeneralResponseModel{
				true, "Sepete ekleme başarısız", nil,
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response.ToJson())
			return
		}
	} else {
		var basketItems []BasketModel
		mapstructure.Decode(basketData, &basketItems)
		basketItems = append(basketItems, newItem)
		error := fb.PushData("/basket/"+userMail, basketItems)
		if error != nil {
			response = util.GeneralResponseModel{
				true, "Sepete ekleme başarısız", nil,
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response.ToJson())
			return
		}
	}

	response = util.GeneralResponseModel{
		false, "Başarılı", nil,
	}
	w.Write(response.ToJson())
}
