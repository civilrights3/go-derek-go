package config

type Config struct {
	Chat       Chat       `yaml:"chat"`
	Multiworld Multiworld `yaml:"multiworld"`
}

func NewDefaultConfig() Config {
	return Config{
		Chat:       newDefaultChat(),
		Multiworld: newDefaultMultiworld(),
	}
}

type Chat struct {
	Key       string `yaml:"-"`
	GuildID   string `yaml:"guild_id"`
	ChannelID string `yaml:"channel_id"`
}

func newDefaultChat() Chat {
	return Chat{} // no default values
}

const (
	defaultClientID         = "163519839402105"
	defaultVersion          = "0.3.8"
	defaultConnectionRetry  = 1800 // 30 minutes
	defaultMultiworldServer = "archipelago.gg"
)

type Multiworld struct {
	ClientID           string `yaml:"client_id,omitempty"`
	ClientVersion      string `yaml:"client_version,omitempty"`
	MaxConnectionRetry int    `yaml:"max_connection_retry"`
	World              World  `yaml:"world,omitempty"`
}

type World struct {
	Server string `yaml:"server,omitempty"`
	Port   string `yaml:"port,omitempty"`
	Slot   string `yaml:"slot,omitempty"`
}

func newDefaultMultiworld() Multiworld {
	return Multiworld{
		ClientID:           defaultClientID,
		ClientVersion:      defaultVersion,
		MaxConnectionRetry: defaultConnectionRetry,
		World: World{
			Server: defaultMultiworldServer,
		},
	}
}
