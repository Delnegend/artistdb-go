package main

import (
	"artistdb-go/src/artist"
	"artistdb-go/src/routes"
	"artistdb-go/src/utils"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/rjeczalik/notify"
)

func init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC1123Z,
		}),
	))
}

func main() {
	appState := utils.NewAppState()

	// read file & parse for 1st time
	func() {
		rawBytes, err := os.ReadFile(appState.GetInFile())
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		artistCount, err2 := artist.ParseToNewDB(appState, string(rawBytes))
		if err2 != nil {
			slog.Error(err2.Message, err2.Props...)
			os.Exit(1)
		}
		slog.Info("parsed artists successfully", "count", artistCount)
	}()

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
				rawBytes, err := os.ReadFile(appState.GetInFile())
				if err != nil {
					slog.Error(err.Error())
				}
				artistCount, err2 := artist.ParseToNewDB(appState, string(rawBytes))
				if err2 != nil {
					slog.Error(err2.Message, err2.Props...)

				}
				slog.Info("parsed artists successfully", "count", artistCount)
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
