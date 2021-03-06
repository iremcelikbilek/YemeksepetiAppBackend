package basket

import (
	"net/http"
	"strconv"
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
	basketCount, basketCountOk := r.URL.Query()["count"]
	var totalBAsketCount int = 1

	if basketCountOk && len(basketCount) == 1 {
		value, err := strconv.Atoi(basketCount[0])
		if err == nil {
			totalBAsketCount = value
		}
	}

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
	var restaurantMenu listing.MunuModel

	for _, value := range restaurants {
		if value.Id == restaurantId[0] {
			restaurantName = value.Name
			for _, menu := range value.Menu {
				if menu.Id == menuId[0] {
					restaurantMenu = menu
					break
				}
			}
			break
		}
	}

	var newItem = BasketModel{restaurantMenu, restaurantName, restaurantId[0]}

	basketData := fb.ReadData("/basket/" + util.MailToPath(userMail))
	if basketData == nil {
		var data []BasketModel
		for i := 1; i <= totalBAsketCount; i++ {
			data = append(data, newItem)
		}
		error := fb.WriteData("/basket/"+util.MailToPath(userMail), data)
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

		for i := 1; i <= totalBAsketCount; i++ {
			basketItems = append(basketItems, newItem)
		}
		
		error := fb.WriteData("/basket/"+util.MailToPath(userMail), basketItems)
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
