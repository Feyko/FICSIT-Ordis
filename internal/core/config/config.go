package config

type OrdisConfig struct {
	Commands CommandsConfig
}

type CommandsConfig struct {
	Persistent bool
}
