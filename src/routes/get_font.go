package routes

import (
	"log/slog"
	"net/http"
	"os"
	"path"
	"regexp"
)

func GetFont(w http.ResponseWriter, r *http.Request) {
	fontName := r.PathValue("fontName")
	if match, _ := regexp.MatchString(`[^a-zA-Z0-9\-_\.]`, fontName); match {
		http.Error(w, "invalid font name", http.StatusBadRequest)
		return
	}

	fontPath := path.Clean("./src/frontend/fonts/" + fontName)
	if _, err := os.Stat(fontPath); err != nil {
		http.Error(w, fontName+" not found", http.StatusNotFound)
		slog.Error(err.Error())
		return
	}

	http.ServeFile(w, r, fontPath)
}
