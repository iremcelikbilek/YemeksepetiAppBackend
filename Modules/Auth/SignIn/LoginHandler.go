package Login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	signup "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Auth/SignUp"
	fb "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Firebase"
	util "github.com/iremcelikbilek/YemeksepetiAppBackend/Modules/Utils"
	"github.com/mitchellh/mapstructure"
)

var JWT_Token = []byte("YEMEK_SEPETI_JWT_TOKEN")

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	util.HeaderManager(&w)
	var response util.GeneralResponseModel
	var loginData LoginModel

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, "Gelen veriler hatalı :(", nil,
		}
		w.Write(response.ToJson())
		return
	}

	if !util.IsEmailValid(loginData.PersonEmail) {
		response = util.GeneralResponseModel{
			true, "eMail geçersiz", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fetchedData := fb.GetFilteredData("/persons", "personEmail", loginData.PersonEmail)
	var userDbData signup.SignUpDbModel
	mapstructure.Decode(fetchedData, &userDbData)

	if userDbData.PersonEmail == "" {
		response = util.GeneralResponseModel{
			true, "Kullanıcı veya şifre hatalı", nil,
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.ToJson())
		return
	}

	if !util.ComparePasswords(userDbData.Password, loginData.Password) {
		response = util.GeneralResponseModel{
			true, "Kullanıcı veya şifre hatalı", nil,
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.ToJson())
		return
	}

	expirationTime := time.Now().Add(750 * time.Hour)
	claims := &Claims{
		Username: loginData.PersonEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_Token)
	if err != nil {
		response = util.GeneralResponseModel{
			true, "Bir Sorun Oluştu", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenData := LoginResponseData{
		Token:        tokenString,
		Expires:      expirationTime.String(),
		UserName:     userDbData.PersonName,
		UserLastName: userDbData.PersonLastName,
	}

	response = util.GeneralResponseModel{
		false, "Giriş Başarılı", tokenData,
	}

	w.Write(response.ToJson())

	nowDate, _ := time.Now().MarshalText()
	itemMap := fetchedData.(map[string]interface{})
	itemMap["signInDate"] = string(nowDate)

	if err := fb.UpdateFilteredData("/persons", "personEmail", loginData.PersonEmail, itemMap); err != nil {
		fmt.Println("Login tarihi güncellenemedi")
	}
}
