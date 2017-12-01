package health

import (
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
)

// HealthService runs application health checks, and provides the /__health http endpoint
type HealthService struct {
	fthealth.HealthCheck
	gtgChecks []fthealth.Check
}

// NewHealthService returns a new HealthService
func NewHealthService(appSystemCode string, appName string, appDescription string) *HealthService {
	service := &HealthService{}
	service.SystemCode = appSystemCode
	service.Name = appName
	service.Description = appDescription
	service.Checks = []fthealth.Check{
		service.skeletonCheck(),
	}

	return service
}

// HealthCheckHandleFunc provides the http endpoint function
func (service *HealthService) HealthCheckHandleFunc() func(w http.ResponseWriter, r *http.Request) {
	return fthealth.Handler(service)
}

func (service *HealthService) skeletonCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "skeleton-healthcheck",
		BusinessImpact:   "None",
		Name:             "Skeleton Healthcheck",
		PanicGuide:       "https://dewey.ft.com/enriched-concepts.html",
		Severity:         1,
		TechnicalSummary: "The app is not running",
		Checker:          service.skeletonPingCheck,
	}
}

func (service *HealthService) skeletonPingCheck() (string, error) {
	return "UPP enriched-concepts is healthy", nil
}

// GTG returns the current gtg status
func (service *HealthService) GTG() gtg.Status {
	return gtg.Status{GoodToGo: true, Message: "OK"}
}
