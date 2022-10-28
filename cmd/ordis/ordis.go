package main

import (
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/ordis"
	"FICSIT-Ordis/internal/ports/gql"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
)

var defaultConf = ordis.Config{
	Arango: arango.Config{
		Username:      "ordis",
		Password:      "pass",
		SuperUsername: "root",
		SuperPassword: "pass",
		DBName:        "ordis",
		Endpoints:     []string{"http://localhost:8529"},
	},
	Auth: auth.Config{
		Secret: "notsecret",
	},
}

func main() {
	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	ord, err := ordis.New(conf)
	if err != nil {
		log.Fatalf("Could not create an Ordis instance: %v", err)
	}

	err = gql.Server(&ord)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfig() (ordis.Config, error) {
	var conf ordis.Config
	useDefault := shouldUseDefaultConfig()
	if useDefault {
		conf = defaultConf
	}
	missing := fillConfigFromEnv(&conf)
	if !useDefault && len(missing) > 0 {
		return ordis.Config{}, errors.Errorf("Missing %v config env variables: %v", len(missing), missing)
	}
	return conf, nil
}

// fillConfigFromEnv fills an ordis.Config from environment variables and returns all missing env var names
func fillConfigFromEnv(conf *ordis.Config) []string {
	missing := make([]string, 0)
	setFromEnv(&conf.Arango.Username, "ORDIS_ARANGO_USERNAME", &missing)
	setFromEnv(&conf.Arango.Password, "ORDIS_ARANGO_PASSWORD", &missing)
	setFromEnv(&conf.Arango.SuperUsername, "ORDIS_ARANGO_SUPER_USERNAME", nil)
	setFromEnv(&conf.Arango.SuperPassword, "ORDIS_ARANGO_SUPER_PASSWORD", nil)
	setFromEnv(&conf.Arango.DBName, "ORDIS_ARANGO_DB_NAME", &missing)
	if conf.Arango.Endpoints == nil {
		conf.Arango.Endpoints = make([]string, 1)
	}
	setFromEnv(&conf.Arango.Endpoints[0], "ORDIS_ARANGO_DB_ENDPOINT", &missing)
	setFromEnv(&conf.Auth.Secret, "ORDIS_ARANGO_AUTH_SECRET", &missing)
	return missing
}

func setFromEnv(p *string, name string, missing *[]string) {
	value, ok := os.LookupEnv(name)
	if !ok {
		if missing != nil {
			*missing = append(*missing, name)
		}
		return
	}
	*p = value
}

func shouldUseDefaultConfig() bool {
	v, ok := os.LookupEnv("ORDIS_USE_DEFAULT_CONFIG")
	if !ok {
		return false
	}
	should, _ := strconv.ParseBool(v)
	return should
}
