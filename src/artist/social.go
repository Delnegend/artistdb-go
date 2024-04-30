package artist

import (
	"artistdb-go/src/socials"
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

	Original string
}

var WRONG_SOCIAL_FORMAT = "social must have a format of username@socialcode[,description] or //link,description"

func (social *Social) Unmarshal(
	supportedSocials *socials.Supported,
	username, original string,
) *utils.SlogError {
	social.Original = original
	social.IsSpecial = false

	fields := strings.Split(original, ",")
	if len(fields) < 1 || len(fields) > 2 {
		return &utils.SlogError{
			Message: WRONG_SOCIAL_FORMAT,
			Errors:  []any{"artist", username, "social", original, "kind", "split with ,", "split result", fields},
		}
	}

	// isSpecial
	if strings.HasPrefix(fields[0], "*") {
		social.IsSpecial = true
	}
	fields[0] = strings.TrimPrefix(fields[0], "*")

	usingCustomLink := strings.HasPrefix(fields[0], "//")
	usingAtSocial := strings.Contains(fields[0], "@")

	if !usingCustomLink && !usingAtSocial {
		return &utils.SlogError{
			Message: WRONG_SOCIAL_FORMAT,
			Errors:  []any{"artist", username, "social", original, "kind", "not using @ or //"},
		}
	}

	if usingCustomLink {
		if len(fields) < 2 {
			return &utils.SlogError{
				Message: "custom social link needs a description",
				Errors:  []any{"artist", username, "social", original},
			}
		}
		social.Link = fields[0]
		social.Description = fields[1]
		return nil
	}

	infoField := strings.Split(fields[0], "@")
	if len(infoField) != 2 {
		return &utils.SlogError{
			Message: WRONG_SOCIAL_FORMAT,
			Errors:  []any{"artist", username, "social", original, "kind", "split with @", "split result", infoField},
		}
	}
	if !supportedSocials.IsSupported(infoField[1]) {
		return &utils.SlogError{
			Message: "this social code is not supported",
			Errors:  []any{"artist", username, "social", infoField[1]},
		}
	}
	social.Username = infoField[0]
	social.SocialCode = infoField[1]

	var err error
	social.Link, err = supportedSocials.ToProfileLink(social.Username, social.SocialCode)
	if err != nil {
		return &utils.SlogError{
			Message: err.Error(),
			Errors:  []any{"artist", username, "socialCode", social.SocialCode},
		}
	}

	descriptionField := func() string {
		if len(fields) == 2 {
			return fields[1]
		}
		return ""
	}()
	social.Description, err = supportedSocials.FormatDescription(social.SocialCode, descriptionField)
	if err != nil {
		return &utils.SlogError{
			Message: err.Error(),
			Errors:  []any{"artist", username, "socialCode", social.SocialCode},
		}
	}

	social.IsSpecial = supportedSocials.IsSpecial(social.SocialCode)

	return nil
}

func (social *Social) Marshal(supportedSocials *socials.Supported) (string, error) {
	if social.Link == "" {
		return "", fmt.Errorf("social link is empty")
	}
	if social.Description == "" {
		return "", fmt.Errorf("social description is empty")
	}
	isSpecialStr := ""
	if social.IsSpecial {
		isSpecialStr = "*"
	}
	return fmt.Sprintf("%s%s,%s", isSpecialStr, social.Link, social.Description), nil
}
