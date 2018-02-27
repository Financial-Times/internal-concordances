# Internal Concordances

Internal UPP Concordances API which accepts a list of concept ids, concords them, and returns them as an array of UPP concepts.

## Installation

Download the source code, dependencies and test dependencies:

```
go get -u github.com/kardianos/govendor
mkdir $GOPATH/src/github.com/Financial-Times/internal-concordances
cd $GOPATH/src/github.com/Financial-Times
git clone https://github.com/Financial-Times/internal-concordances.git
cd internal-concordances && govendor sync
go build .
```

## Running locally

1. Run the tests and install the binary:

```
govendor sync
govendor test -v -race +local
go install
```

2. Run the binary (using the `help` flag to see the available optional arguments):

```
$GOPATH/bin/internal-concordances [--help]

Usage: internal-concordances [OPTIONS]

UPP Internal Concordances

Options:
      --app-system-code                System Code of the application (env $APP_SYSTEM_CODE) (default "internal-concordances")
      --app-name                       Application name (env $APP_NAME) (default "internal-concordances")
      --concept-search-api-endpoint    Endpoint to query for concepts (env $CONCEPT_SEARCH_ENDPOINT) (default "http://localhost:8111")
      --public-concordances-endpoint   Endpoint to concord ids with (env $PUBLIC_CONCORDANCES_ENDPOINT) (default "http://localhost:8222")
      --port                           Port to listen on (env $APP_PORT) (default "8080")
      --api-yml                        Location of the API Swagger YML file. (env $API_YML) (default "./api.yml")
```

3. Test:

```
curl http://localhost:8080/__health | jq
```

## Build and deployment

* Built by Docker Hub on merge to master: [coco/internal-concordances](https://hub.docker.com/r/coco/internal-concordances/)
* CI provided by CircleCI: [internal-concordances](https://circleci.com/gh/Financial-Times/internal-concordances)

## Service endpoints

For a full description of API endpoints for the service, please see the [Open API specification](./api/api.yml).

## Healthchecks

Admin endpoints are:

`/__gtg`
`/__health`
`/__build-info`

At the moment, both the `/__gtg` and `/__health` endpoints perform no checks (effectively a ping of the service).

### Logging

* The application uses [logrus](https://github.com/sirupsen/logrus); the log file is initialised in [main.go](main.go).
* Logging requires an `env` app parameter, for all environments other than `local` logs are written to file.
* When running locally, logs are written to console. If you want to log locally to file, you need to pass in an env parameter that is != `local`.
* NOTE: `/__build-info` and `/__gtg` endpoints are not logged as they are called every second from varnish/vulcand and this information is not needed in logs/splunk.
