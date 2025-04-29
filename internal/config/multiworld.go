package config

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
