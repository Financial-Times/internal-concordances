package health

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/stretchr/testify/assert"
)

func TestHealthServiceHandler(t *testing.T) {
	health := NewHealthService("appSystemCode", "appName", "appDescription", happyCheckMock, unhappyCheckMock)

	handler := health.HealthCheckHandleFunc()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/__health", nil)

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)

	dec := json.NewDecoder(w.Body)
	var r fthealth.HealthResult
	err := dec.Decode(&r)
	assert.NoError(t, err)
	assert.Len(t, r.Checks, 2)

	assert.False(t, r.Ok)
}

func TestGTGAllGood(t *testing.T) {
	health := NewHealthService("appSystemCode", "appName", "appDescription", happyCheckMock, happyCheckMock)

	gtg := health.GTG()
	assert.True(t, gtg.GoodToGo)
	assert.Equal(t, "OK", gtg.Message)
}

func TestNotGTG(t *testing.T) {
	health := NewHealthService("appSystemCode", "appName", "appDescription", happyCheckMock, unhappyCheckMock)

	gtg := health.GTG()
	assert.False(t, gtg.GoodToGo)
	assert.Equal(t, "computer says no", gtg.Message)
}

var happyCheckMock = fthealth.Check{
	ID:               "happy-check",
	BusinessImpact:   "A big impact",
	Name:             "Mock API Healthcheck",
	PanicGuide:       "https://runbooks.in.ft.com/mock-api.html",
	Severity:         1,
	TechnicalSummary: "Mock API is not available",
	Checker:          happyChecker,
}

func happyChecker() (string, error) {
	return "I'm happy!", nil
}

var unhappyCheckMock = fthealth.Check{
	ID:               "unhappy-check",
	BusinessImpact:   "A massive impact",
	Name:             "Mock API Healthcheck",
	PanicGuide:       "https://runbooks.in.ft.com/mock-api.html",
	Severity:         1,
	TechnicalSummary: "Mock API is not available",
	Checker:          unhappyChecker,
}

func unhappyChecker() (string, error) {
	return "", errors.New("computer says no")
}
