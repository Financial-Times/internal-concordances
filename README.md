# Internal Concordances

[![Circle CI](https://circleci.com/gh/Financial-Times/internal-concordances.svg?style=shield)](https://circleci.com/gh/Financial-Times/internal-concordances)
[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/internal-concordances)](https://goreportcard.com/report/github.com/Financial-Times/internal-concordances)
[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/internal-concordances/badge.svg)](https://coveralls.io/github/Financial-Times/internal-concordances)

Internal UPP Concordances API which accepts a list of concept IDs, concords them, and returns them as an array of UPP concepts.

## Installation

Download the source code, dependencies and test dependencies:

```
git clone https://github.com/Financial-Times/internal-concordances.git
cd internal-concordances 
go build -mod=readonly
```

## Running locally

1. Run the tests and install the binary:

```
go test -mod=readonly -race ./...
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
      --concept-search-api-endpoint    Endpoint to query for concepts (env $CONCEPT_SEARCH_ENDPOINT) (default "http://concept-search-api:8080")
      --public-concordances-endpoint   Endpoint to concord ids with (env $PUBLIC_CONCORDANCES_ENDPOINT) (default "http://public-concordances-api:8080")
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

For a full description of API endpoints for the service, please see the [Open API specification](./_ft/api.yml).

## Healthchecks

Admin endpoints are:

`/__gtg`
`/__health`
`/__build-info`

At the moment, both the `/__gtg` and `/__health` endpoints perform no checks (effectively a ping of the service).

### Logging

* The application uses [logrus](https://github.com/sirupsen/logrus) wrapped by [go-logger](https://github.com/Financial-Times/go-logger); the log file is initialised in [main.go](main.go).
* NOTE: `/__build-info` and `/__gtg` endpoints are not logged as they are called every second from varnish and this information is not needed in logs/splunk.
txt
