package socials

import (
	"fmt"
	"strings"
)

type Social struct {
	DisplayName string
	Profile     string
}

type Supported struct {
	unavatar map[string]Social
	extended map[string]Social
	special  map[string]bool
}

func NewInstance() Supported {
	return Supported{
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
			"twitter":       {"𝕏", "twitter.com/<@>"},
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

func (supported *Supported) IsUnavatarSupported(socialCode string) bool {
	if socialCode == "x" {
		return true
	}
	if _, ok := supported.unavatar[socialCode]; ok {
		return true
	}
	return false
}
func (supported *Supported) IsSupported(socialCode string) bool {
	if socialCode == "x" {
		return true
	}
	if _, ok := supported.unavatar[socialCode]; ok {
		return true
	}
	if _, ok := supported.extended[socialCode]; ok {
		return true
	}
	return false
}

func (supported *Supported) ToProfileLink(username, socialCode string) (string, error) {
	if socialCode == "x" {
		socialCode = "twitter"
	}
	if _, ok := supported.unavatar[socialCode]; ok {
		return strings.Replace(supported.unavatar[socialCode].Profile, "<@>", username, 1), nil
	}
	if _, ok := supported.extended[socialCode]; ok {
		return strings.Replace(supported.extended[socialCode].Profile, "<@>", username, 1), nil
	}
	return "", fmt.Errorf("social code not found to create profile link")
}

func (supported *Supported) FormatDescription(socialCode, description string) (string, error) {
	var socialName string

	// find social display name
	if social, ok := supported.unavatar[socialCode]; ok {
		socialName = social.DisplayName
	}
	if socialName == "" {
		if social, ok := supported.extended[socialCode]; ok {
			socialName = social.DisplayName
		}
	}

	// no social name, return description
	if socialName == "" && description != "" {
		return description, nil
	}

	// have social name, no description
	if socialName != "" && description == "" {
		return socialName, nil
	}

	// have both
	if socialName != "" && description != "" {
		return fmt.Sprintf("%s | %s", socialName, description), nil
	}

	return "", fmt.Errorf("social code not found, description is empty, can't format new description")
}

func (supported *Supported) IsSpecial(socialCode string) bool {
	if _, ok := supported.special[socialCode]; ok {
		return true
	}
	return false
}
