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
