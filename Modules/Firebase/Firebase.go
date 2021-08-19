package Firebase

import (
	"context"
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

	firebaseClient, err := app.DatabaseWithURL(ctx, "https://irem-yemeksepeti-backend-default-rtdb.firebaseio.com/")
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

func DeletePath(path string) error {
	err := client.NewRef(path).Delete(ctx)
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

var credentials string = `
{
	"type": "service_account",
	"project_id": "irem-yemeksepeti-backend",
	"private_key_id": "0810c49de71c17bffe3ecb4357b9a32b56a242a6",
	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDGHdPTVdffixCf\nUKDPcVNwvwdrxhMIqven4XzN6HVT+RMu1axuRKaI12Dh3k0DAfEtqmHIwR4ZzQ12\niJZI/PIEyhcfDDbiBg3lijDm5iBkplo/8/CooPXmzfEubViGo2QorM5WSp8xQuKo\niF4jpNIkY+obkzIEWD7fiFxVYjP0lNn0NtzIkC4g5mR9eP8C38J3k9fYUUw1Be78\npdj1YX/y71NwBLv6vRpCuaoLfjhCW0EwZrnp9skHNJS7ith6tJquVAu839jZ7xMG\nd94VmMpYl/YxyLX/tzWUCCPBYpwisMviSPHA3KxKvaVyjx1en5Ki5T6DocCnJTIt\nwKpfDyErAgMBAAECggEAC1ZwpWL/wCC+ukdMdKKpIkYkYBQNSc0y8A4U7Nm1QTF7\nWg5LWGIgX6tntXVZ1ea6DSF3iBwZI2PbNeHaK+Ih3YlNKm8yAtxS1kSCyOv5hZkJ\niChnKNdRSzyU5VHHo6jdFgDRrBmII7MOspNfQ83uYru/DYXuclY0fulYU2CT1ZbH\nNTZqHSEitCy8DXhP61GPs5Nu/i6K4rRUyKHbeHdllqO4FzIWc1Itoc7fUP7ADlO8\ndrsfcKslNQA1YvFa+BIhFzvjA8R36iM32ca6Od95VDnepwoE0pQvOt7dTdW/UScI\nfqqdckiEAFfhmL8pSbC6zM1YhNychaPhOYRwgJdW4QKBgQD7i2SpLO8wgFuRfCme\n52BcU72RUs5dGh8svdlUsFYufuLusrGphjBVjI0gQZ7MGKuN9GUvP8vwbjvsRiIl\nr3UJdJtz2obPMi4F7NDT2+MqUjiClKDZbWROAOnrqtsHKJglD2qK/HgyFkpExUUZ\nbuDqHkPjazoH1Qb0uUDLf4a9iwKBgQDJoCttE+pyfzyeoznqfZ6F545uACr6wqEr\nxiTee+6uLcGaBTi/sv8r2Yf186RO10eNqS2XGhHfT5mmtS3dzR5iEZ8MSXO2yIr1\nNd8lwrBSRBBq6tI1MoYcXj8BQ3lFiC9E3Efiw5N0CYZSvl66rBDxhHMfoPrKcrwh\nCGcGBEze4QKBgHvGQ2nbanb7MhOMfQ5r28aSjh0MGe9GA0EIygAaJM4MMa4yz6kT\nFoWB+497up/DI+dd8swlIDzWgTXp7LOOepCEiFmhleQuVOcleDxHXqhcfOIEMIHM\niia33GLSV6RWHUdfJpXtVVeQEEt2pmG1ZYbODanCAXQJJrsUzQVVYv+xAoGAT9H5\n/yfIQ+W9QOxLrFpo3IgMKd4lJbrRhXve8rlLh2cT4v64NaQOQvTOT39SB+hQKnPU\nWaJ3etmPcaD+dHWU1qw1M+8MQUtpP6RBIDjQBvFtMnaeG3NSBn8FIGHu66j7VZ6D\nUvGsOV7f73fwFqSx3Htb/CSFxInhko46AvbG2+ECgYEAuR8uGuqTPPe7EJTFlEx/\n31B4yJzoMnGMmKLTVtEGcALZ42KH+BPp8YaTysl96Tl0ZnbWINeUXnEanWbSJLkt\nPaYHKdGlqom9TLEPh5lMb0DcLRMQ5/pqTGsvKifulpZ7DHm6Q89iRbmlSX9aT8rZ\nFEXj+H1WahGR2E2fGHotEVE=\n-----END PRIVATE KEY-----\n",
	"client_email": "firebase-adminsdk-2vi7n@irem-yemeksepeti-backend.iam.gserviceaccount.com",
	"client_id": "103029863676139527792",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-2vi7n%40irem-yemeksepeti-backend.iam.gserviceaccount.com"
  }
`
