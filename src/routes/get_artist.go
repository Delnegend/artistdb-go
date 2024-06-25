package routes

import (
	"artistdb-go/src/artist"
	"artistdb-go/src/utils"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

func GetArtist(appState *utils.AppState) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := strings.ToLower(r.PathValue("username"))

		aliasModel := new(artist.AliasDB)
		err := appState.DB.NewSelect().Model(aliasModel).Where("alias = ?", username).Scan(r.Context())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				appState.ArtistNotFoundTmpl.Tmpl.Execute(w, nil)
				return
			}
			slog.Error("failed to get artist", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		artistModel := new(artist.ArtistDB)
		err = appState.DB.NewSelect().Model(artistModel).Where("id = ?", aliasModel.ID).Scan(r.Context())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				slog.Error("alias found but artist not found", "alias", username)
				appState.ArtistNotFoundTmpl.Tmpl.Execute(w, nil)
				return
			}
			slog.Error("failed to get artist", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		socialLines := strings.Split(artistModel.Socials, "\n")
		socials := make([]template.HTML, 0, len(socialLines))
		for _, socialLine := range socialLines {
			isSpecial := false
			if strings.HasPrefix(socialLine, "*") {
				isSpecial = true
				socialLine = socialLine[1:]
			}
			socialInfo := strings.Split(socialLine, ",")
			if len(socialInfo) != 2 {
				http.Error(w, "DB contains invalid social line", http.StatusInternalServerError)
				slog.Error("invalid social line", "artist", username, "line", socialLine)
				return
			}
			socials = append(socials, appState.SocialLinkTmpl.RenderAsHTML(utils.LinkPageFields{
				IsSpecial:   isSpecial,
				Link:        socialInfo[0],
				Description: socialInfo[1],
			}))
		}

		avatar := artistModel.Avatar
		if strings.HasPrefix(avatar, "//") {
			avatar = "https:" + avatar
		}
		appState.ArtistPageTmpl.Execute(w, utils.IndexPageFields{
			Title:        fmt.Sprintf("%s | ArtistDB", artistModel.DisplayName),
			Favicon:      avatar,
			ArtistAvatar: avatar,
			DisplayName:  artistModel.DisplayName,
			Links:        socials,
		})
	}
}
