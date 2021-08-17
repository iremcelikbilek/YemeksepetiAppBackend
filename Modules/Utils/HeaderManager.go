package Utils

import "net/http"

func HeaderManager(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")

}
