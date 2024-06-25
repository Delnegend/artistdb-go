package utils

import (
	"database/sql"
	"log/slog"
	"os"
	"strconv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type AppState struct {
	port      string
	inFile    string
	avatarDir string

	SocialLinkTmpl     *HTMLTemplate
	ArtistPageTmpl     *HTMLTemplate
	ArtistNotFoundTmpl *HTMLTemplate

	// To check duplicate usernames and aliases
	UsernameSet map[string]struct{}
	AliasSet    map[string]struct{}

	SupportedSocials SupportedSocials

	DB *bun.DB
}

func NewAppState() (*AppState, error) {
	as := AppState{
		port: "8080",
	}

	sqldbPath := os.Getenv("SQLITE")
	if sqldbPath == "" {
		sqldbPath = "./sqlite.db?mode=rwc"
	}
	sqldb, err := sql.Open(sqliteshim.ShimName, sqldbPath)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	as.DB = bun.NewDB(sqldb, sqlitedialect.New())

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
	fileStat, err := os.Stat(as.inFile)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if fileStat.IsDir() {
		slog.Error(as.inFile + " is a directory")
		os.Exit(1)
	}

	as.avatarDir = os.Getenv("AVATAR_DIR")
	if as.avatarDir == "" {
		as.avatarDir = "avatars"
	}

	as.SocialLinkTmpl = &HTMLTemplate{}
	if err := as.SocialLinkTmpl.Read("./frontend/link.html"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	as.ArtistPageTmpl = &HTMLTemplate{}
	if err := as.ArtistPageTmpl.Read("./frontend/index.html"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	as.ArtistNotFoundTmpl = &HTMLTemplate{}
	if err := as.ArtistNotFoundTmpl.Read("./frontend/not_found.html"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	as.SupportedSocials = NewSocialDBInstance()

	return &as, nil
}

func (as *AppState) GetPort() string {
	return as.port
}
func (as *AppState) GetInFile() string {
	return as.inFile
}
func (as *AppState) GetAvatarDir() string {
	return as.avatarDir
}
