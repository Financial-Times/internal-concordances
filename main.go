package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Financial-Times/api-endpoint"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	"github.com/Financial-Times/internal-concordances/concepts"
	"github.com/Financial-Times/internal-concordances/health"
	"github.com/Financial-Times/internal-concordances/resources"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/husobee/vestigo"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

const appDescription = "UPP Internal Concordances"

func main() {
	app := cli.App("internal-concordances", appDescription)

	appSystemCode := app.String(cli.StringOpt{
		Name:   "app-system-code",
		Value:  "internal-concordances",
		Desc:   "System Code of the application",
		EnvVar: "APP_SYSTEM_CODE",
	})

	appName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  "internal-concordances",
		Desc:   "Application name",
		EnvVar: "APP_NAME",
	})

	conceptSearchEndpoint := app.String(cli.StringOpt{
		Name:   "concept-search-api-endpoint",
		Value:  "http://concept-search-api:8080",
		Desc:   "Endpoint to query for concepts",
		EnvVar: "CONCEPT_SEARCH_ENDPOINT",
	})

	publicConcordancesEndpoint := app.String(cli.StringOpt{
		Name:   "public-concordances-endpoint",
		Value:  "http://public-concordances-api:8080",
		Desc:   "Endpoint to concord ids with",
		EnvVar: "PUBLIC_CONCORDANCES_ENDPOINT",
	})

	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})

	apiYml := app.String(cli.StringOpt{
		Name:   "api-yml",
		Value:  "./api.yml",
		Desc:   "Location of the API Swagger YML file.",
		EnvVar: "API_YML",
	})

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	app.Action = func() {
		log.Infof("[Startup] %v is starting", *appSystemCode)
		log.Infof("System code: %s, App Name: %s, Port: %s", *appSystemCode, *appName, *port)

		client := &http.Client{Timeout: 8 * time.Second}

		search := concepts.NewSearch(client, *conceptSearchEndpoint)
		concordances := concepts.NewConcordances(client, *publicConcordancesEndpoint)

		healthService := health.NewHealthService(*appSystemCode, *appName, appDescription, search.Check(), concordances.Check())

		serveEndpoints(*port, apiYml, healthService, search, concordances)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}

func serveEndpoints(port string, apiYml *string, healthService *health.HealthService, search concepts.Search, concordances concepts.Concordances) {
	r := vestigo.NewRouter()

	var monitoringRouter http.Handler = r
	monitoringRouter = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), monitoringRouter)
	monitoringRouter = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, monitoringRouter)

	r.Get("/__health", healthService.HealthCheckHandleFunc())
	r.Get(status.GTGPath, status.NewGoodToGoHandler(healthService.GTG))
	r.Get(status.BuildInfoPath, status.BuildInfoHandler)

	r.Get("/internalconcordances", resources.InternalConcordances(concordances, search))

	http.Handle("/", monitoringRouter)

	if apiYml != nil {
		apiEndpoint, err := api.NewAPIEndpointForFile(*apiYml)
		if err != nil {
			log.WithError(err).WithField("file", *apiYml).Warn("Failed to serve the API Endpoint for this service. Please validate the Swagger YML and the file location")
		} else {
			r.Get(api.DefaultPath, apiEndpoint.ServeHTTP)
		}
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Unable to start: %v", err)
	}
}
