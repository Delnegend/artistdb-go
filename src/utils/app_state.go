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

func NewAppState() *AppState {
	return &AppState{
		port: func() string {
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
			return port
		}(),
		inFile: func() string {
			inFile := os.Getenv("IN_FILE")
			if inFile == "" {
				return "artists.txt"
			}
			fileStat, err := os.Stat(inFile)
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			if fileStat.IsDir() {
				slog.Error(inFile + " is a directory")
				os.Exit(1)
			}
			return inFile
		}(),
		avatarDir: func() string {
			avatarDir := os.Getenv("AVATAR_DIR")
			if avatarDir == "" {
				return "avatars"
			}
			fileStat, err := os.Stat(avatarDir)
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			if !fileStat.IsDir() {
				slog.Error(avatarDir + " is not a directory")
				os.Exit(1)
			}
			return avatarDir
		}(),

		SocialLinkTmpl: func() *HTMLTemplate {
			st := &HTMLTemplate{}
			if err := st.Read("./frontend/link.html"); err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			return st
		}(),
		ArtistPageTmpl: func() *HTMLTemplate {
			st := &HTMLTemplate{}
			if err := st.Read("./frontend/index.html"); err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			return st
		}(),
		ArtistNotFoundTmpl: func() *HTMLTemplate {
			st := &HTMLTemplate{}
			if err := st.Read("./frontend/not_found.html"); err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			return st
		}(),

		UsernameSet:      make(map[string]struct{}),
		AliasSet:         make(map[string]struct{}),
		SupportedSocials: NewSocialDBInstance(),

		DB: func() *bun.DB {
			sqldbPath := os.Getenv("SQLITE")
			if sqldbPath == "" {
				sqldbPath = "./sqlite.db?mode=rwc"
			}
			sqldb, err := sql.Open(sqliteshim.ShimName, sqldbPath)
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			return bun.NewDB(sqldb, sqlitedialect.New())
		}(),
	}
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
