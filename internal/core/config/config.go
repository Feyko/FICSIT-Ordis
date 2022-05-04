package config

type OrdisConfig struct {
	Commands CommandsConfig
}

type CommandsConfig struct {
	Persistent bool
	Arango     ArangoConfig
}

type ArangoConfig struct {
	Username,
	Password,
	RootPassword,
	DBName string
	Endpoints []string
}
