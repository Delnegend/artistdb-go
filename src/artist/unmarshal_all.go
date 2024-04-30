package artist

import (
	"artistdb-go/src/socials"
	"log/slog"
	"regexp"
	"slices"
	"sort"
)

// Unmarshal parses the artist string and returns a map of artists. Key: pointer
// to username, Value: Artist struct.
func Unmarshal(supportedSocials *socials.Supported, artistString string) []Artist {
	rgx, _ := regexp.Compile(`\n{2,}`)
	artistsSplitted := rgx.Split(artistString, -1)

	artists := make([]Artist, 0)

	for _, artistString := range artistsSplitted {
		artist := Artist{}
		if err := artist.Unmarshal(
			supportedSocials,
			artistString,
		); err != nil {
			slog.Error((*err).Message, (*err).Errors...)
			continue
		}
		artists = append(artists, artist)
	}

	// handle duplicate usernames and aliases
	usernames := make([]string, 0)
	aliases := make([]string, 0)

	for _, artist := range artists {
		if slices.Contains(usernames, artist.Username) {
			slog.Warn("duplicate username", "username", artist.Username)
			continue
		}
		usernames = append(usernames, artist.Username)
	}

	for _, artist := range artists {
		for _, alias := range artist.Alias {
			aliasExisted := slices.Contains(aliases, alias)
			aliasInUsernames := slices.Contains(usernames, alias)

			if aliasExisted || aliasInUsernames {
				slog.Warn("duplicate alias", "alias", alias)
				continue
			}
			aliases = append(aliases, alias)
		}
	}

	sort.Slice(artists, func(i, j int) bool {
		return artists[i].Username < artists[j].Username
	})

	return artists
}
