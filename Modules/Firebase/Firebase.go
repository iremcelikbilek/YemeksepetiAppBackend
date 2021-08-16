package Firebase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	database "firebase.google.com/go/db"
)

var ctx = context.Background()
var client *database.Client

func ConnectFirebase() {

	opt := option.WithCredentialsJSON([]byte(credentials))

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}

	firebaseClient, err := app.DatabaseWithURL(ctx, "https://irem-yemeksepeti-app-default-rtdb.europe-west1.firebasedatabase.app/")
	if err != nil {
		log.Fatal(err)
	}

	client = firebaseClient
}

func WriteData(path string, data interface{}) error {
	err := client.NewRef(path).Set(ctx, data)
	return err

}

func PushData(path string, data interface{}) error {
	_, err := client.NewRef(path).Push(ctx, data)
	return err
}

func PushFilteredData(path string, child string, equal string, newChild string, pushData interface{}) error {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}
	PushData(path+"/"+dataParentName+"/"+newChild, pushData)
	if err != nil {
		return err
	}
	return nil
}

func GetFilteredData(path string, child string, equal string) interface{} {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		fmt.Println(err)
	}
	itemsMap := data.(map[string]interface{})

	var responseData interface{}
	for _, v := range itemsMap {
		responseData = v
		break
	}
	return responseData
}

func UpdateFilteredData(path string, child string, equal string, updatedData interface{}) error {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	newData := map[string]interface{}{
		dataParentName: updatedData,
	}

	err = client.NewRef(path).Update(ctx, newData)
	if err != nil {
		return err
	}
	return nil
}

func CommentAdd(advertisementID string, comment interface{}) error {
	var data interface{}
	err := client.NewRef("/advertisement").OrderByChild("advertisementID").EqualTo(advertisementID).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	if dataParentName == "" {
		return errors.New("İlan bulunamadı")
	}

	return PushData("/advertisement/"+dataParentName+"/comments", comment)
}

func UpdateUserSpesificData(path string, child string, equal string, updatedData interface{}, user string) error {
	var data interface{}
	err := client.NewRef("/persons").OrderByChild("personEmail").EqualTo(user).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	return UpdateFilteredData("/persons/"+dataParentName+path, child, equal, updatedData)
}

func DeleteComment(advID string, commentID string, userMail string) error {
	var data interface{}
	err := client.NewRef("/advertisement").OrderByChild("advertisementID").EqualTo(advID).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})
	if itemsMap == nil {
		return errors.New("İlan bulunamadı")
	}

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	if dataParentName == "" {
		return errors.New("İlan bulunamadı")
	}

	advData := itemsMap[dataParentName]
	if advData == nil {
		return errors.New("İlan bulunamadı")
	}
	advDataMap := advData.(map[string]interface{})
	if advDataMap == nil {
		return errors.New("İlan bulunamadı")
	}

	commentsData := advDataMap["comments"]
	if commentsData == nil {
		return errors.New("Yorum bulunamadı")
	}
	commentsMap := commentsData.(map[string]interface{})

	if commentsMap == nil {
		return errors.New("Yorum bulunamadı")
	}
	for i, v := range commentsMap {
		comment := v.(map[string]interface{})
		if comment == nil {
			return errors.New("Yorum bulunamadı")
		}
		if comment["commentID"] == commentID {
			if comment["personEmail"].(string) == userMail {
				return client.NewRef("advertisement/" + dataParentName + "/comments/" + i).Delete(ctx)
			} else {
				return errors.New("Yorumu silmeye yetkiniz yok")
			}
			break
		}
	}
	return errors.New("Yorumu bulunamadı")
}

func Delete(path string, child string, equal string) error {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	err = client.NewRef(path + "/" + dataParentName).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func ReadData(path string) interface{} {
	var data interface{}
	if err := client.NewRef(path).Get(ctx, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func DeleteAllFilteredDatas(path string, child string, equal string) {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return
	}
	itemsMap := data.(map[string]interface{})

	for i, _ := range itemsMap {
		client.NewRef(path + "/" + i).Delete(ctx)
	}
}

func DeleteFavoriteAdvertisement(id string, user string) error {
	var data interface{}
	err := client.NewRef("/persons").OrderByChild("personEmail").EqualTo(user).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	userData := itemsMap[dataParentName].(map[string]interface{})
	favorites := userData["favorites"].(map[string]interface{})

	for index, value := range favorites {
		if value == id {
			client.NewRef("/persons/" + dataParentName + "/favorites/" + index).Delete(ctx)
		}
	}

	return nil
}

var credentials string = `
{
	"type": "service_account",
	"project_id": "irem-yemeksepeti-app",
	"private_key_id": "5154511c969c66905ea4d6bc794847c6491febd5",
	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC/gTYCovcGY/0A\nNlqqRu0IZqPfL3rzqA9JrLuo0CG1G5YNoYr9Q42RoQrZmj06U6fVKJZCyJ0hd1/F\nGaCfKOYFkQ55z7qJA2MR1Y1wQCR2W97l7WmKE00v6u9sjU0j3riZnIpf8eBNGKXw\nVtnsItnSQKzT2dV7wO7wp5/3+hDNU7+i0r/kW3ktjSr8Cz+2gR9NrIEg5ZVQRh4P\nit/XBjwLp2zrkrLRLYNPZCCDDeKe4mNUn3g+pqNm97Bd5jfOHQz9ECQBjeyoqDml\nHK/AaWSm5Hy89xRkYn/rXx0lbDgpozrE3na5DJcfkWkJvK1/03j6VxozaBYALz/q\ntmI0AAT3AgMBAAECggEACPBr8tPkP7DzfTWMUON5AAnj1LC03A0jEuIpwQJfeG0+\nZAoqNeftLU2STJWwYvnHB/evipQuZYKwKT/+KT/G+56JnvJX4jpoFEvey3fAW/gD\nimxlH0lcXS6/GKxAhrqdl6oLeYQNtNDq5LJL1NkhnXi1uI6JXQ/591N+b5xazsK+\nhOvWL41Xb2ZHUGBE4cVNTbA2qQ9Lekgu5eWA7D/1j3BCq7Sq1YGdjturKzw2ZLRa\nnaw0sBwBwVsbKzr4UyXYfWrUDCsyZRnNqUWXKjl+sfUxQe3Kyy4PqEw705jDBJJm\nMBd+qAppNkAk/FCIbCH8N0m9gTLKk3kcdKf7owQISQKBgQDeh1cAqkfwCncYZCB5\nLiIsjw2bHrYoDkcGE28nO/VvgEh9wIbSB/yM2CbFjlCIj7ziFkxtxiCGG3FXrwj5\npbUN6q1D+ug9UkzmoND3NF2BeRufho4104kVjc2lydcIizUBkrn74k72pOR2+wYi\nVXeCDwZoQoYPYLpMClWCxNJE7wKBgQDcT0RsciQ+1ZeGtFkv/2wfQDPIdaS/Qrjh\n8mU1t+d9TiTV0KeNiOYsOGcYzHvxXWWw1DOS50MNzBnrTylDnYGXgZL+57vbBdZQ\n9dNSBec9sygefS/yRoXp3kSJi4ok699B3ytnnd9fLda5s1DZWnFbsUHC1gmrKchy\nPkMHfZmQeQKBgQDOeLT5IQXua0dlkkGvLmb3ASSWsUBCmjy8Hnwb4z4vXs/kHib5\n6f8ij8wpsYp3qyaOgDIaCKNUy1G3Eek5+c6sQvrRAJVLkHlZ5Az/0c6Qu1YuBiMd\nPlELdq9BDK5AdymPdBys4aZyozx4SSG/6Z0hR9+iDVdmHVG+DDibRRP0cQKBgAZA\nqax6QNUXssk77RwTn7nzVITn8dkLx7uB6aVwpr1Drn/zAA5gSEgRAbwOcaYUBILU\nQvJ0Zc7KcCHhiUZF/huSrd1WLlq0+7QoheraCAoUP5s96lJx9fMBP+i3cSBDIX75\nGn5CWMiWwHVcxXqlunnjuf4RnQyijvHPGo/n3KfhAoGBAIlESS3x1ts9CoCt+o+y\nCEdLpDlHzm2WdkwY3vV+uXopHAQdPCjASK8hSi0IBYXuvHrKuOO9zdI6YcXvbPZe\nlxJfByTiskGAv2Jy5+kNY+WuFoBH37bqi0ufjqGnRqTtEgzZuG2vRw2Wtc6G9Kjp\nFoeuly+URuYwzehIpTFVysIs\n-----END PRIVATE KEY-----\n",
	"client_email": "firebase-adminsdk-ypjna@irem-yemeksepeti-app.iam.gserviceaccount.com",
	"client_id": "118389491062640486144",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-ypjna%40irem-yemeksepeti-app.iam.gserviceaccount.com"
  }

`
