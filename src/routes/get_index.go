package routes

import "net/http"

func GetIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/Delnegend/artistdb-go", http.StatusTemporaryRedirect)
}
