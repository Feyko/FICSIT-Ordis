package config

type OrdisConfig struct {
	Arango   ArangoConfig
	Commands CommandsConfig
}

type CommandsConfig struct {
}

type ArangoConfig struct {
	Username,
	Password,
	SuperUsername,
	SuperPassword,
	DBName string
	Endpoints []string
}
