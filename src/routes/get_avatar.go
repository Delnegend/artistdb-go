package routes

import (
	"artistdb-go/src/utils"
	"log/slog"
	"net/http"
	"os"
	"path"
	"regexp"
)

func GetAvatar(appState *utils.AppState) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("fileName")

		if fileName == "default" {
			if _, err := os.Stat("./src/frontend/avatar.svg"); err != nil {
				http.Error(w, "avatar.svg not found", http.StatusNotFound)
				slog.Error(err.Error())
				return
			}
			http.ServeFile(w, r, "./src/frontend/avatar.svg")
			return
		}

		if match, _ := regexp.MatchString(`[^a-zA-Z0-9\._]`, fileName); match {
			http.Error(w, "invalid file name", http.StatusBadRequest)
			return
		}

		// check file exists in AVATAR_DIR
		filePath := path.Join(appState.GetAvatarDir(), fileName)
		if _, err := os.Stat(filePath); err != nil {
			http.Error(w, "avatar not found", http.StatusNotFound)
			slog.Error(err.Error())
			return
		}

		http.ServeFile(w, r, filePath)
	}
}
