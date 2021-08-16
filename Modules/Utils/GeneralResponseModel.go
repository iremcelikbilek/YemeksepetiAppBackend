package Utils

import "encoding/json"

type GeneralResponseModel struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (res GeneralResponseModel) ToJson() []byte {
	responseJSON, err := json.Marshal(res)
	if err != nil {
		return []byte("Undefined Error")
	}
	return responseJSON
}
