package routes

import (
	"log/slog"
	"net/http"
	"os"
)

func StyleCSS(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("./frontend/dist.css"); err != nil {
		http.Error(w, "style.css not found", http.StatusNotFound)
		slog.Error(err.Error())
		return
	}
	http.ServeFile(w, r, "./frontend/dist.css")
}
