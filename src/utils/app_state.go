package utils

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type AppState struct {
	port           string
	inFile         string
	outDir         string
	avatarDir      string
	fallbackAvatar string
	formatAndExit  bool

	SocialLinkTmpl     *HTMLTemplate
	ArtistPageTmpl     *HTMLTemplate
	ArtistNotFoundTmpl *HTMLTemplate
}

func NewAppState() (*AppState, error) {
	as := AppState{
		port: "8080",
	}

	port := os.Getenv("PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		slog.Error("invalid port number")
		os.Exit(1)
	}
	if portInt < 1024 || portInt > 65535 {
		slog.Error("invalid port number")
		os.Exit(1)
	}
	as.port = port
	as.inFile = os.Getenv("IN_FILE")
	if as.inFile == "" {
		as.inFile = "artists.txt"
	}
	as.outDir = os.Getenv("OUT_DIR")
	if as.outDir == "" {
		as.outDir = "artists"
	}
	as.avatarDir = os.Getenv("AVATAR_DIR")
	if as.avatarDir == "" {
		as.avatarDir = "avatars"
	}
	as.formatAndExit = false
	formatAndExit := os.Getenv("FORMAT_AND_EXIT")
	if strings.ToLower(formatAndExit) == "true" {
		as.formatAndExit = true
	}
	as.fallbackAvatar = os.Getenv("FALLBACK_AVATAR")

	as.SocialLinkTmpl = &HTMLTemplate{}
	if err := as.SocialLinkTmpl.Read("./src/frontend/link.html"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	as.ArtistPageTmpl = &HTMLTemplate{}
	if err := as.ArtistPageTmpl.Read("./src/frontend/index.html"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	as.ArtistNotFoundTmpl = &HTMLTemplate{}
	if err := as.ArtistNotFoundTmpl.Read("./src/frontend/not_found.html"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return &as, nil
}

func (as *AppState) GetPort() string {
	return as.port
}

func (as *AppState) GetInFile() string {
	return as.inFile
}

func (as *AppState) GetOutDir() string {
	return as.outDir
}

func (as *AppState) GetAvatarDir() string {
	return as.avatarDir
}

func (as *AppState) GetFallbackAvatar() string {
	return as.fallbackAvatar
}

func (as *AppState) GetFormatAndExit() bool {
	return as.formatAndExit
}
