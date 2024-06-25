package artist

import (
	"artistdb-go/src/utils"
	"fmt"
	"strings"
)

type Social struct {
	Username    string
	SocialCode  string
	Description string
	Link        string
	IsSpecial   bool
}

var WRONG_SOCIAL_FORMAT = "social must have a format of username@socialcode[,description] or //link,description"

func (social *Social) Unmarshal(
	appState *utils.AppState,
	username, rawString string,
) *SlogErr {
	social.IsSpecial = false

	// <link || username@socialcode>,description
	slice := strings.Split(rawString, ",")
	if len(slice) < 1 || len(slice) > 2 {
		return &SlogErr{
			Message: WRONG_SOCIAL_FORMAT,
			Props:   []any{"artist", username, "social", rawString, "kind", "split with ,", "split result", slice},
		}
	}

	// *<link || username@socialcode>,description
	if strings.HasPrefix(slice[0], "*") {
		social.IsSpecial = true
		slice[0] = strings.TrimPrefix(slice[0], "*")
	}

	usingCustomLink := strings.HasPrefix(slice[0], "//")
	usingAtSocial := strings.Contains(slice[0], "@")

	switch {
	case usingCustomLink:
		if len(slice) < 2 {
			return &SlogErr{
				Message: "custom social link needs a description",
				Props:   []any{"artist", username, "social", rawString},
			}
		}
		social.Link = slice[0]
		social.Description = slice[1]
	case usingAtSocial:
		// username@socialcode
		subSlice := strings.Split(slice[0], "@")
		if len(subSlice) != 2 {
			return &SlogErr{
				Message: WRONG_SOCIAL_FORMAT,
				Props:   []any{"artist", username, "social", rawString, "kind", "split with @", "split result", subSlice},
			}
		}
		social.Username = subSlice[0]
		social.SocialCode = subSlice[1]

		socialLink, err := appState.SupportedSocials.
			ToProfileLink(social.Username, social.SocialCode)
		if err != nil {
			return &SlogErr{
				Message: err.Error(),
				Props:   []any{"artist", username, "socialCode", social.SocialCode},
			}
		}
		social.Link = socialLink

		var description string
		if len(slice) == 2 {
			description = slice[1]
		}
		description, err = appState.SupportedSocials.
			FormatDescription(social.SocialCode, description)
		if err != nil {
			return &SlogErr{
				Message: err.Error(),
				Props:   []any{"artist", username, "socialCode", social.SocialCode},
			}
		}
		social.Description = description
	default:
		return &SlogErr{
			Message: WRONG_SOCIAL_FORMAT,
			Props:   []any{"artist", username, "social", rawString, "kind", "not using @ or //"},
		}
	}
	return nil
}

func (social *Social) Marshal() (string, error) {
	switch {
	case social.Link != "" && social.Description != "" && social.IsSpecial:
		return fmt.Sprintf("*%s,%s", social.Link, social.Description), nil
	case social.Link != "" && social.Description != "":
		return fmt.Sprintf("%s,%s", social.Link, social.Description), nil
	default:
		return "", fmt.Errorf("social link and description are empty")
	}
}
