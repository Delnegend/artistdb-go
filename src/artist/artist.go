package artist

import (
	"artistdb-go/src/utils"
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

var (
	WRONG_AVATAR_FORMAT = "avatar must have a format of username@socialcode, leave empty or use underscore to auto infer"
)

type SlogErr struct {
	Message string
	Props   []any
}

type ArtistDB struct {
	bun.BaseModel `bun:"table:artist"`

	ID          string `bun:"id,pk,unique,notnull"`
	DisplayName string `bun:"display_name"`
	Avatar      string `bun:"avatar"`
	Socials     string `bun:"socials"`

	Aliases []string `bun:"-"`
}

type AliasDB struct {
	bun.BaseModel `bun:"table:alias"`

	Alias string `bun:"alias,unique,pk"`
	ID    string `bun:"artist_id,notnull"`
}

// ParseToNewDB parses the artist string and returns a map of artists. Key: pointer
// to username, Value: Artist struct.
func ParseToNewDB(appState *utils.AppState, artistString string) (int, *SlogErr) {
	// reset database & temp vars
	startTimer := time.Now()
	if err := appState.DB.
		ResetModel(
			context.Background(),
			(*ArtistDB)(nil), (*AliasDB)(nil)); err != nil {
		return 0, &SlogErr{
			Message: err.Error(),
		}
	}
	appState.UsernameSet = make(map[string]struct{})
	appState.AliasSet = make(map[string]struct{})
	slog.Info("reset DB, clear temp vars", "time", time.Since(startTimer))

	// split, rm empty lines, sort
	rgx, _ := regexp.Compile(`\n{2,}`)
	artistRawStrings := rgx.Split(artistString, -1)
	sort.Slice(artistRawStrings, func(i, j int) bool {
		return artistRawStrings[i] < artistRawStrings[j]
	})

	// parse artists into DB models
	startTimer = time.Now()
	artistsToDB := make([]ArtistDB, 0)
	for _, artistString := range artistRawStrings {
		artist := Artist{}
		artistModel, err := artist.Unmarshal(appState, artistString)
		if err != nil {
			return 0, err
		}
		artistsToDB = append(artistsToDB, artistModel)
	}
	slog.Info("artists parsed to DB models", "time", time.Since(startTimer))

	// insert into DB
	startTimer = time.Now()
	if _, err := appState.DB.NewInsert().
		Model(&artistsToDB).
		Exec(context.Background()); err != nil {
		return 0, &SlogErr{
			Message: err.Error(),
		}
	}
	slog.Info("artists inserted into DB", "time", time.Since(startTimer))

	// prepare alias into DB models
	aliasesToDB := make([]AliasDB, 0)
	for _, artist := range artistsToDB {
		for _, alias := range artist.Aliases {
			aliasesToDB = append(aliasesToDB, AliasDB{
				ID:    artist.ID,
				Alias: alias,
			})
		}
		aliasesToDB = append(aliasesToDB, AliasDB{
			ID:    artist.ID,
			Alias: artist.ID,
		})
	}

	// insert
	startTimer = time.Now()
	if _, err := appState.DB.NewInsert().
		Model(&aliasesToDB).
		Exec(context.Background()); err != nil {
		return 0, &SlogErr{
			Message: err.Error(),
		}
	}
	slog.Info("aliases inserted into DB", "time", time.Since(startTimer))

	appState.UsernameSet = make(map[string]struct{})
	appState.AliasSet = make(map[string]struct{})
	return len(artistsToDB), nil
}

type Artist struct {
	Original string
}

// Unmarshal parses the artists.txt and inserts into database
func (artist *Artist) Unmarshal(appState *utils.AppState, rawString string) (ArtistDB, *SlogErr) {
	artist.Original = rawString

	// split, rm empty lines, check length
	lines := make([]string, 0)
	for _, line := range strings.Split(rawString, "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) < 2 {
		return ArtistDB{}, &SlogErr{
			Message: "artist has no social info",
			Props:   []any{"artist", rawString},
		}
	}

	// parse info
	infoData := strings.Split(lines[0], ",")
	artist.Original = lines[0]

	// username, display name & check duplicate username
	var username, displayName string
	username = strings.ToLower(infoData[0])
	if (len(infoData) > 1) && (infoData[1] != "") && (infoData[1] != "_") {
		displayName = infoData[1]
	} else {
		displayName = username
	}
	if _, ok := appState.UsernameSet[username]; ok {
		return ArtistDB{}, &SlogErr{
			Message: "duplicate username found in username pool",
			Props:   []any{"artist", username},
		}
	}
	if _, ok := appState.AliasSet[username]; ok {
		return ArtistDB{}, &SlogErr{
			Message: "username found in alias pool",
			Props:   []any{"artist", username},
		}
	}
	appState.UsernameSet[username] = struct{}{}

	// alias & check duplicate
	alias := func() []string {
		aliasMap := make(map[string]struct{}, 0)
		if len(infoData) > 3 {
			for i := 3; i < len(infoData); i++ {
				if infoData[i] != "" {
					aliasMap[strings.ToLower(infoData[i])] = struct{}{}
				}
			}
		}
		aliasSlice := make([]string, 0)
		for alias := range aliasMap {
			aliasSlice = append(aliasSlice, alias)
		}
		return aliasSlice
	}()
	for _, alias := range alias {
		if _, ok := appState.UsernameSet[alias]; ok {
			return ArtistDB{}, &SlogErr{
				Message: "alias found in username pool",
				Props:   []any{"artist", username, "alias", alias},
			}
		}

		if _, ok := appState.AliasSet[alias]; ok {
			return ArtistDB{}, &SlogErr{
				Message: "duplicate alias found in alias pool",
				Props:   []any{"artist", username, "alias", alias},
			}
		}
		appState.AliasSet[alias] = struct{}{}
	}

	// socials
	socials := make([]Social, 0)
next_social:
	for _, social := range lines[1:] {
		artistSocial := Social{}
		if err := artistSocial.Unmarshal(appState, username, social); err != nil {
			slog.Error((*err).Message, (*err).Props...)
			continue next_social
		}
		socials = append(socials, artistSocial)
	}

	// avatar
	var avatar string
	if len(infoData) > 2 {
		autoInfer := infoData[2] == "_" || infoData[2] == ""
		usingAtSocial := strings.Contains(infoData[2], "@")
		usingAbsPath := strings.HasPrefix(infoData[2], "/")

		switch {
		case usingAtSocial:
			components := strings.Split(infoData[2], "@")
			if len(components) != 2 {
				return ArtistDB{}, &SlogErr{
					Message: WRONG_AVATAR_FORMAT,
					Props:   []any{"artist", username, "avatar", infoData[2]},
				}
			}

			result, err := appState.SupportedSocials.
				ToUnavatarLink(components[0], components[1])
			if err != nil {
				return ArtistDB{}, &SlogErr{
					Message: err.Error(),
					Props:   []any{"artist", username, "social", components[1]},
				}
			}
			avatar = result
		case usingAbsPath:
			avatar = fmt.Sprintf("/avatar%s", infoData[2])
		case autoInfer:
			for _, social := range socials {
				result, err := appState.SupportedSocials.
					ToUnavatarLink(social.Username, social.SocialCode)
				if err != nil {
					continue
				}
				avatar = result
				break
			}
			if avatar == "" {
				return ArtistDB{}, &SlogErr{
					Message: "could not infer avatar from socials",
					Props:   []any{"artist", username, "socials", socials},
				}
			}
		default:
			return ArtistDB{}, &SlogErr{
				Message: WRONG_AVATAR_FORMAT,
				Props:   []any{"artist", username, "avatar", infoData[2]},
			}
		}
	}

	socialsMarshaled := make([]string, 0)
	for _, social := range socials {
		socialMarshaled, err := social.Marshal()
		if err != nil {
			slog.Error("can't marshal social", "artist", username, "error", err.Error(), "social", social)
			continue
		}
		socialsMarshaled = append(socialsMarshaled, socialMarshaled)
	}

	return ArtistDB{
		ID:          username,
		DisplayName: displayName,
		Avatar:      avatar,
		Socials:     strings.Join(socialsMarshaled, "\n"),
		Aliases:     alias,
	}, nil
}
