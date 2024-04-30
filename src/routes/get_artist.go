package routes

import (
	"artistdb-go/src/utils"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func GetArtist(appState *utils.AppState) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		if match, _ := regexp.MatchString(`[^a-zA-Z0-9\._]`, username); match {
			http.Error(w, "invalid username", http.StatusBadRequest)
			return
		}

		// check if file exists
		filePath := path.Join(appState.GetOutDir(), username)
		if _, err := os.Stat(filePath); err != nil {
			http.Error(w, "artist not found", http.StatusNotFound)
			slog.Error(err.Error())
			return
		}

		// read the file
		file, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "artist not found", http.StatusNotFound)
			slog.Error(err.Error())
			return
		}

		// if file is an alias
		data := string(file)
		if strings.HasPrefix(data, "@") {
			file, err := os.ReadFile(path.Join(appState.GetOutDir(), data[1:]))
			if err != nil {
				http.Error(w, "avatar not found", http.StatusNotFound)
				slog.Error(err.Error())
				return
			}
			data = string(file)
		}

		lines := strings.Split(data, "\n")

		// prepare display name and avatar
		infoLine := strings.Split(lines[0], ",")
		if len(infoLine) != 2 {
			http.Error(w, "invalid artist file", http.StatusBadRequest)
			return
		}
		displayName := infoLine[0]

		avatar := infoLine[1]
		if strings.HasPrefix(avatar, "/") {
			avatar = "/avatar" + avatar
		} else {
			avatar = "https://unavatar.io/" + avatar + "?size=400"
			if appState.GetFallbackAvatar() != "" {
				avatar += "&fallback=" + appState.GetFallbackAvatar()
			}
		}

		// prepare social links
		socialLines := lines[1:]
		socials := make([]template.HTML, 0, len(socialLines))
		for _, socialLine := range socialLines {
			isSpecial := false
			if strings.HasPrefix(socialLine, "*") {
				isSpecial = true
				socialLine = socialLine[1:]
			}
			socialInfo := strings.Split(socialLine, ",")
			if len(socialInfo) != 2 {
				http.Error(w, "invalid artist file", http.StatusBadRequest)
				slog.Error("invalid social line", "artist", username, "line", socialLine)
				return
			}
			socials = append(socials, appState.SocialLinkTmpl.RenderAsHTML(utils.LinkPageFields{
				IsSpecial:   isSpecial,
				Link:        socialInfo[0],
				Description: socialInfo[1],
			}))
		}

		appState.ArtistPageTmpl.Execute(w, utils.IndexPageFields{
			Title:        fmt.Sprintf("%s | ArtistDB", displayName),
			ArtistAvatar: avatar,
			DisplayName:  displayName,
			Links:        socials,
		})
	}
}
