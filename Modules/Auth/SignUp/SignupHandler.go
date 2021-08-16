package SignUp

import (
	"fmt"
	"net/http"
	//fb "../../Firebase"
	//util "../../Utils"
)

var JWT_Token = []byte("YEMEK_SEPETI_JWT_TOKEN")

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}
