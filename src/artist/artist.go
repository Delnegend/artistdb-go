package artist

import (
	"artistdb-go/src/socials"
	"artistdb-go/src/utils"
	"fmt"
	"log/slog"
	"strings"
)

type Artist struct {
	Username    string
	DisplayName string
	Avatar      string
	Alias       []string
	Socials     []Social

	Original string
}

var (
	WRONG_AVATAR_FORMAT  = "avatar must have a format of username@socialcode, leave empty or use underscore to auto infer"
	UNAVATAR_NOT_SUPPORT = "this social code is not supported for unavatar"
)

func (artist *Artist) Unmarshal(supportedSocials *socials.Supported, artistStrign string) *utils.SlogError {
	// split, rm empty lines, check length
	lines := make([]string, 0)
	for _, line := range strings.Split(artistStrign, "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) < 2 {
		return &utils.SlogError{
			Message: "artist has no social info",
			Errors:  []any{"artist", artistStrign},
		}
	}

	// parse info
	infoData := strings.Split(lines[0], ",")
	artist.Original = lines[0]

	// username
	if len(infoData) > 0 {
		artist.Username = strings.ToLower(infoData[0])
		artist.DisplayName = artist.Username
	}

	// display name
	if len(infoData) > 1 {
		artist.DisplayName = infoData[1]
	}
	if artist.DisplayName == "" || artist.DisplayName == "_" {
		artist.DisplayName = artist.Username
	}

	// avatar
	if len(infoData) > 2 {
		autoInfer := infoData[2] == "_"
		usingAtSocial := strings.Contains(infoData[2], "@")
		usingAbsPath := strings.HasPrefix(infoData[2], "/")

		if !autoInfer && !usingAtSocial && !usingAbsPath {
			return &utils.SlogError{
				Message: WRONG_AVATAR_FORMAT,
				Errors:  []any{"artist", artist.Username, "avatar", infoData[2]},
			}
		}

		if usingAtSocial {
			components := strings.Split(infoData[2], "@")
			if len(components) != 2 {
				return &utils.SlogError{
					Message: WRONG_AVATAR_FORMAT,
					Errors:  []any{"artist", artist.Username, "avatar", infoData[2]},
				}
			}

			socialUsername := components[0]
			socialCode := components[1]
			if socialCode == "x" {
				socialCode = "twitter"
			}

			if !(*supportedSocials).IsUnavatarSupported(socialCode) {
				return &utils.SlogError{
					Message: UNAVATAR_NOT_SUPPORT,
					Errors:  []any{"artist", artist.Username, "social", socialCode},
				}
			}

			artist.Avatar = fmt.Sprintf("%s/%s", socialCode, socialUsername)
		} else if usingAbsPath {
			artist.Avatar = infoData[2]
		} else if autoInfer {
			artist.Avatar = ""
		}

		// using auto-infer -> left blank, handle after parsing socials done
	}

	// alias
	artist.Alias = make([]string, 0)
	if len(infoData) > 3 {
		for i := 3; i < len(infoData); i++ {
			if infoData[i] != "" {
				artist.Alias = append(artist.Alias, strings.ToLower(infoData[i]))
			}
		}
	}

	// socials
	artist.Socials = make([]Social, 0)
next_social:
	for _, social := range lines[1:] {
		artistSocial := Social{}
		if err := artistSocial.Unmarshal(supportedSocials, artist.Username, social); err != nil {
			slog.Error((*err).Message, (*err).Errors...)
			continue next_social
		}
		artist.Socials = append(artist.Socials, artistSocial)
	}

	// handle auto-infer avatar
	if artist.Avatar == "" {
		for _, artistSocial := range artist.Socials {
			if !supportedSocials.IsUnavatarSupported(artistSocial.SocialCode) {
				continue
			}
			socialCode := artistSocial.SocialCode
			if socialCode == "x" {
				socialCode = "twitter"
			}
			artist.Avatar = fmt.Sprintf("%s/%s", socialCode, artistSocial.Username)
			break
		}
	}
	if artist.Avatar == "" {
		return &utils.SlogError{
			Message: "could not infer avatar from socials",
			Errors:  []any{"artist", artist.Username, "socials", artist.Socials},
		}
	}

	return nil
}

// Returns a map of file names to content to write
func (artist *Artist) Marshal(supportedSocials *socials.Supported) (map[string]string, error) {
	if artist.Username == "" || artist.DisplayName == "" || artist.Avatar == "" {
		return nil, fmt.Errorf("can't marshal artist with empty username, display name, or avatar")
	}

	infoLine := fmt.Sprintf("%s,%s", artist.DisplayName, artist.Avatar)

	socialLines := make([]string, 0)
	for _, social := range artist.Socials {
		socialLine, err := social.Marshal(supportedSocials)
		if err != nil {
			slog.Warn("can't marshal social", "artist", artist.Username, "error", err.Error(), "social", social)
		}
		socialLines = append(socialLines, socialLine)
	}

	result := map[string]string{
		artist.Username: fmt.Sprintf("%s\n%s", infoLine, strings.Join(socialLines, "\n"))}

	for _, alias := range artist.Alias {
		result[alias] = fmt.Sprintf("@%s", artist.Username)
	}

	return result, nil
}
