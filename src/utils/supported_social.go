package utils

import (
	"fmt"
	"strings"
)

// Social
type Social struct {
	DisplayName string
	Profile     string
}

// SupportedSocials
type SupportedSocials struct {
	unavatar map[string]Social
	extended map[string]Social
	special  map[string]bool
}

func NewSocialDBInstance() SupportedSocials {
	return SupportedSocials{
		unavatar: map[string]Social{
			"deviantart":    {"DeviantArt", "deviantart.com/<@>"},
			"dribbble":      {"Dribbble", "dribbble.com/<@>"},
			"duckduckgo":    {"DuckDuckGo", ""},
			"facebook":      {"Facebook", "fb.com/<@>"},
			"fb":            {"Facebook", "fb.com/<@>"},
			"github":        {"GitHub", "github.com/<@>"},
			"google":        {"Google", ""},
			"gravatar":      {"Gravatar", ""},
			"instagram":     {"Instagram", "instagram.com/<@>"},
			"microlink":     {"Microlink", ""},
			"readcv":        {"ReadCV", "read.cv/<@>"},
			"reddit":        {"Reddit", "reddit.com/user/<@>"},
			"soundcloud":    {"SoundCloud", "soundcloud.com/<@>"},
			"subscribestar": {"SubscribeStar", "subscribestar.adult/<@>"},
			"substack":      {"Substack", "<@>.substack.com/"},
			"telegram":      {"Telegram", "t.me/<@>"},
			"x":             {"𝕏", "x.com/<@>"},
			"youtube":       {"YouTube", "youtube.com/@<@>"},
		},
		extended: map[string]Social{
			"artstation": {"ArtStation", "www.artstation.com/<@>"},
			"bluesky":    {"BlueSky", "bsky.app/profile/<@>"},
			"boosty":     {"Boosty", "boosty.to/<@>"},
			"booth":      {"Booth.pm", "<@>.booth.pm"},
			"bsky":       {"BlueSky", "bsky.app/profile/<@>"},
			"carrd.co":   {"Carrd.co", "<@>.carrd.co"},
			"fa":         {"FurAffinity 🐾", "www.furaffinity.net/user/<@>/"},
			"fanbox":     {"PixivFanbox", "<@>.fanbox.cc"},
			"gumroad":    {"Gumroad", "<@>.gumroad.com"},
			"itaku":      {"Itaku", "itaku.ee/profile/<@>"},
			"itch.io":    {"Itch.io", "itch.io/profile/<@>"},
			"kofi":       {"Ko-fi 🍵", "ko-fi.com/<@>"},
			"linktr.ee":  {"Linktr.ee 🌲", "linktr.ee/<@>"},
			"lit.link":   {"Lit.link", "lit.link/<@>"},
			"patreon":    {"Patreon", "www.patreon.com/<@>"},
			"picarto":    {"Picarto", "www.picarto.tv/<@>"},
			"pixiv":      {"Pixiv", "www.pixiv.net/en/users/<@>"},
			"plurk":      {"Plurk", "plurk.com/<@>"},
			"potofu.me":  {"Potofu.me", "potofu.me/<@>"},
			"skeb":       {"Skeb.jp", "skeb.jp/@<@>"},
			"threads":    {"Threads", "www.threads.net/@<@>"},
			"tumblr":     {"Tumblr", "<@>.tumblr.com"},
			"twitch":     {"Twitch", "www.twitch.tv/<@>"},
			"x":          {"𝕏", "twitter.com/<@>"},
		},
		special: map[string]bool{
			"potofu.me": true,
			"carrd.co":  true,
			"linktr.ee": true,
			"lit.link":  true,
		},
	}
}

func (ss *SupportedSocials) ToUnavatarLink(username, socialCode string) (string, error) {
	if _, ok := ss.unavatar[socialCode]; ok {
		return fmt.Sprintf("//unavatar.io/%s/%s", socialCode, username), nil
	}
	return "", fmt.Errorf("SupportedSocials.ToUnavatarLink: social code not found to create avatar link")
}

func (ss *SupportedSocials) ToProfileLink(username, socialCode string) (string, error) {
	if _, ok := ss.unavatar[socialCode]; ok {
		return strings.Replace(ss.unavatar[socialCode].Profile, "<@>", username, 1), nil
	}
	if _, ok := ss.extended[socialCode]; ok {
		return strings.Replace(ss.extended[socialCode].Profile, "<@>", username, 1), nil
	}
	return "", fmt.Errorf("SupportedSocials.ToProfileLink: social code not found to create profile link")
}

func (ss *SupportedSocials) FormatDescription(socialCode, description string) (string, error) {
	var socialName string

	// find social display name
	if social, ok := ss.unavatar[socialCode]; ok {
		socialName = social.DisplayName
	}
	if socialName == "" {
		if social, ok := ss.extended[socialCode]; ok {
			socialName = social.DisplayName
		}
	}

	switch {
	case socialName == "" && description != "":
		return description, nil
	case socialName != "" && description == "":
		return socialName, nil
	case socialName != "" && description != "":
		return fmt.Sprintf("%s | %s", socialName, description), nil
	default:
		return "", fmt.Errorf("SupportedSocials.FormatDescription: social code not found, description is empty, can't format new description")
	}
}

func (ss *SupportedSocials) IsSpecial(socialCode string) bool {
	if _, ok := ss.special[socialCode]; ok {
		return true
	}
	return false
}
