package main

import (
	"artistdb-go/src/artist"
	"artistdb-go/src/routes"
	"artistdb-go/src/socials"
	"artistdb-go/src/utils"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/rjeczalik/notify"
)

var SUPPORTED_SOCIALS = socials.NewInstance()

func init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC1123Z,
		}),
	))
}

func main() {
	appState, err := utils.NewAppState()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// check if inFile is a file
	fileStat, err := os.Stat(appState.GetInFile())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if fileStat.IsDir() {
		slog.Error(appState.GetInFile() + " is a directory")
		os.Exit(1)
	}

	// read file & parse
	artistsByte, err := os.ReadFile(appState.GetInFile())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	artistsString := string(artistsByte)
	artists := artist.Unmarshal(&SUPPORTED_SOCIALS, artistsString)
	slog.Info("parsed artists successfully", "count", len(artists))

	// sort artist by name & overwrite artists.txt
	if appState.GetFormatAndExit() {
		result := make([]string, 0, len(artists))
		for _, artist := range artists {
			result = append(result, artist.Original)
			for _, social := range artist.Socials {
				result = append(result, social.Original)
			}
			result = append(result, "")
		}

		dataToWrite := strings.Join(result, "\n")
		if err := os.WriteFile(appState.GetInFile(), []byte(dataToWrite), 0644); err != nil {
			slog.Error(err.Error())
		}
		slog.Info("formatted "+appState.GetInFile()+" successfully", "count", len(artists))
		os.Exit(0)
	}

	if err := artist.BuildArtistsDir(&SUPPORTED_SOCIALS, &artists, appState.GetOutDir()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	artists = nil

	// watch artists.txt for changes and re-parse
	eventInfoCh := make(chan notify.EventInfo, 1)
	if err := notify.Watch(appState.GetInFile(), eventInfoCh, notify.InCloseWrite); err != nil {
		slog.Error(err.Error())
	}
	defer notify.Stop(eventInfoCh)
	go func() {
		for {
			switch event := <-eventInfoCh; event.Event() {
			case notify.InCloseWrite:
				artistsByte, err := os.ReadFile(appState.GetInFile())
				if err != nil {
					slog.Error(err.Error())
				}
				artistsString = string(artistsByte)
				artists = artist.Unmarshal(&SUPPORTED_SOCIALS, artistsString)
				slog.Info("parsed artists successfully", "count", len(artists))

				if err := artist.BuildArtistsDir(&SUPPORTED_SOCIALS, &artists, appState.GetOutDir()); err != nil {
					slog.Error(err.Error())
					os.Exit(1)
				}
				artists = nil
			}
		}
	}()

	http.HandleFunc("GET /", routes.GetIndex)
	http.HandleFunc("GET /{username}", routes.GetArtist(appState))
	http.HandleFunc("GET /style.css", routes.StyleCSS)
	http.HandleFunc("GET /avatar/{fileName}", routes.GetAvatar(appState))
	http.HandleFunc("GET /font/{fontName}", routes.GetFont)

	slog.Info("listening on port " + appState.GetPort())
	http.ListenAndServe(":"+appState.GetPort(), nil)
}
