package config

type OrdisConfig struct {
	Arango   ArangoConfig
	Commands CommandsConfig
}

type CommandsConfig struct{}

type AuthConfig struct {
	Secret string
}

type ArangoConfig struct {
	Username,
	Password,
	SuperUsername,
	SuperPassword,
	DBName string
	Endpoints []string
}
