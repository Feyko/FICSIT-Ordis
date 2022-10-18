# FICSIT-Ordis
Helper and ticket manager backend for the Satisfactory Modding community meant to integrate with Discord through [FICSIT-Fred](https://github.com/Feyko/FICSIT-Fred) as well as our mod manager

## Development

### Prerequisites
- Install Go 1.18+
- Run `go generate -x ./...` 
  - The command will error. You will have to go to the mentioned file and replace the `errors` import by `github.com/pkg/errors`. This is an issue with gqlgen but I will implement a workaround
- Have an Arango database on localhost:8529 and with `pass` as the root password

### Running
Simply `go run cmd/ordis` as it is not configurable yet

### Testing
**All tests must pass, both with the in-memory and arango repos**  
Run `go test ./...` to run all tests

Configure the following env vars to make the tests use the Arango db:  
ORDIS_TEST_ARANGO_USER  
ORDIS_TEST_ARANGO_PASSWORD  
ORDIS_TEST_ARANGO_ENDPOINT