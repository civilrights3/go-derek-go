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
	defaultVersion          = "0.5.0"
	defaultConnectionRetry  = 1800 // 30 minutes
	defaultMultiworldServer = "archipelago.gg"

	defaultCacheFilepath = "./cache"
)

type Multiworld struct {
	ClientID           string `yaml:"client_id,omitempty"`
	ClientVersion      string `yaml:"client_version,omitempty"`
	MaxConnectionRetry int    `yaml:"max_connection_retry"`
	World              World  `yaml:"world,omitempty"`
	Cache              Cache  `yaml:"cache,omitempty"`
}

type World struct {
	Server string `yaml:"server,omitempty"`
	Port   string `yaml:"port,omitempty"`
	Slot   string `yaml:"slot,omitempty"`
}

type Cache struct {
	Filepath string
}

func newDefaultMultiworld() Multiworld {
	return Multiworld{
		ClientID:           defaultClientID,
		ClientVersion:      defaultVersion,
		MaxConnectionRetry: defaultConnectionRetry,
		World: World{
			Server: defaultMultiworldServer,
		},
		Cache: Cache{
			Filepath: defaultCacheFilepath,
		},
	}
}
