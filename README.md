# FICSIT-Ordis
Helper and ticket manager backend for the Satisfactory Modding community meant to integrate with Discord through [FICSIT-Fred](https://github.com/Feyko/FICSIT-Fred) as well as our mod manager

## Development

### Prerequisites
- Install Go 1.18+
- Run `go generate -x ./...` 
  - The command will error. You will have to go to the mentioned file and replace the `errors` import by `github.com/pkg/errors`. This is an issue with gqlgen but I will implement a workaround
- Have an Arango database (docker-compose file provided)

### Running
Configure the following env vars  
You can use the env var `ORDIS_USE_DEFAULT_CONFIG` with a truthy value to pull the default values, which also supresses the missing env vars error

|             Name            | Optional |        Default        |
|:---------------------------:|:--------:|:---------------------:|
| ORDIS_ARANGO_USERNAME       |    no    |         ordis         |
| ORDIS_ARANGO_PASSWORD       |    no    |          pass         |
| ORDIS_ARANGO_SUPER_USERNAME |    yes   |          root         |
| ORDIS_ARANGO_SUPER_PASSWORD |    yes   |          pass         |
| ORDIS_ARANGO_DB_NAME        |    no    |         ordis         |
| ORDIS_ARANGO_DB_ENDPOINT    |    no    | http://localhost:8259 |
| ORDIS_ARANGO_AUTH_SECRET    |    no    |       notsecret       |

Then simply run `go run cmd/ordis`

### Testing
**All tests must pass, both with the in-memory and arango repos**  
Run `go test ./...` to run all tests

Configure the following env vars to make the tests use the Arango db:  
ORDIS_TEST_ARANGO_USER  
ORDIS_TEST_ARANGO_PASSWORD  
ORDIS_TEST_ARANGO_ENDPOINT