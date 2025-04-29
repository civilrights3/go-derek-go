package config

type DisplayMode string

const (
	DisplayPlain      DisplayMode = "plain"
	DisplayMonospaced DisplayMode = "mono"
	DisplayColor      DisplayMode = "color"
)

type Chat struct {
	Key         string      `yaml:"-"`
	GuildID     string      `yaml:"guild_id"`
	ChannelID   string      `yaml:"channel_id"`
	DisplayMode DisplayMode `yaml:"display_mode"`
}

func newDefaultChat() Chat {
	return Chat{
		DisplayMode: DisplayPlain,
	}
}
