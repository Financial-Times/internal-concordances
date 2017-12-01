# Enriched Concepts

UPP Concepts API which accepts a list of concept ids, concords them, and returns them as an array of standard UPP concepts.

## Installation

Download the source code, dependencies and test dependencies:

```
go get -u github.com/kardianos/govendor
mkdir $GOPATH/src/github.com/Financial-Times/enriched-concepts
cd $GOPATH/src/github.com/Financial-Times
git clone https://github.com/Financial-Times/enriched-concepts.git
cd enriched-concepts && govendor sync
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
$GOPATH/bin/enriched-concepts [--help]

Usage: enriched-concepts [OPTIONS]

UPP Enriched Concepts

Options:                  
      --app-system-code   System Code of the application (env $APP_SYSTEM_CODE) (default "enriched-concepts")
      --app-name          Application name (env $APP_NAME) (default "enriched-concepts")
      --port              Port to listen on (env $APP_PORT) (default "8080")
      --api-yml           Location of the API Swagger YML file. (env $API_YML) (default "./api.yml")
```

3. Test:

```
curl http://localhost:8080/__health | jq
```

## Build and deployment 

* Built by Docker Hub on merge to master: [coco/enriched-concepts](https://hub.docker.com/r/coco/enriched-concepts/)
* CI provided by CircleCI: [enriched-concepts](https://circleci.com/gh/Financial-Times/enriched-concepts)

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
