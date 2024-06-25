package routes

import (
	"artistdb-go/src/utils"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
)

func GetAvatar(appState *utils.AppState) func(w http.ResponseWriter, r *http.Request) {
	filesInAvatarDir := make(map[string]struct{})
	err := filepath.WalkDir(appState.GetAvatarDir(), func(path string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err
		case d.IsDir():
			return nil
		default:
			filesInAvatarDir[d.Name()] = struct{}{}
			return nil
		}
	})
	if err != nil {
		slog.Error("can't scan avatar dir", "error", err.Error())
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("fileName")
		if fileName == "" || fileName == "default" {
			http.ServeFile(w, r, "./frontend/avatar.svg")

		}
		if _, ok := filesInAvatarDir[fileName]; ok {
			http.ServeFile(w, r, fileName)
		}

		http.Error(w, "avatar not found", http.StatusNotFound)
	}
}
