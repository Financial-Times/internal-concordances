package concepts

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/husobee/vestigo"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSearchByIDsNoResults(t *testing.T) {
	serverMock := new(mockConceptSearchAPI)
	requestedUUIDs := []string{uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestSearchByIDsNoResults", requestedUUIDs)
	serverMock.On("getResponse").Return(`{}`, http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	search := NewSearch(&http.Client{}, server.URL)
	concepts, err := search.ByIDs("tid_TestSearchByIDsNoResults", requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, concepts, 0)
	serverMock.AssertExpectations(t) // failure here means the search API has not been called
}

func TestSearchByIDs(t *testing.T) {
	serverMock := new(mockConceptSearchAPI)
	requestedUUIDs := []string{uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestSearchByIDs", requestedUUIDs)

	searchResp, err := ioutil.ReadFile("./_fixtures/search_response.json")
	require.NoError(t, err)
	serverMock.On("getResponse").Return(string(searchResp), http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	search := NewSearch(&http.Client{}, server.URL)
	concepts, err := search.ByIDs("tid_TestSearchByIDs", requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, concepts, 1)
	serverMock.AssertExpectations(t) // failure here means the search API has not been called
}

func TestSearchNoIDsProvided(t *testing.T) {
	search := NewSearch(&http.Client{}, "")
	_, err := search.ByIDs("tid_TestSearchNoIDsProvided")

	assert.EqualError(t, err, ErrNoConceptsToSearch.Error())
}

func TestSearchAllIDsProvidedEmpty(t *testing.T) {
	search := NewSearch(&http.Client{}, "")
	_, err := search.ByIDs("tid_TestSearchNoIDsProvided", "", "", "", "")

	assert.EqualError(t, err, ErrConceptIDsAreEmpty.Error())
}

func TestSearchRequestURLInvalid(t *testing.T) {
	search := NewSearch(&http.Client{}, ":#")
	_, err := search.ByIDs("tid_TestSearchRequestURLInvalid", uuid.NewV4().String())

	assert.Error(t, err)
}

func TestSearchRequestFails(t *testing.T) {
	search := NewSearch(&http.Client{}, "#:")
	_, err := search.ByIDs("tid_TestSearchRequestFails", uuid.NewV4().String())

	assert.Error(t, err)
}

func TestSearchResponseFailed(t *testing.T) {
	serverMock := new(mockConceptSearchAPI)
	requestedUUIDs := []string{uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestSearchResponseFailed", requestedUUIDs)
	serverMock.On("getResponse").Return(`{"message":"forbidden!!!!!"}`, http.StatusForbidden)

	server := serverMock.startServer(t)
	defer server.Close()

	search := NewSearch(&http.Client{}, server.URL)
	_, err := search.ByIDs("tid_TestSearchResponseFailed", requestedUUIDs...)

	assert.EqualError(t, err, "403 Forbidden: forbidden!!!!!")
	serverMock.AssertExpectations(t) // failure here means the search API has not been called
}

func TestSearchResponseInvalidJSON(t *testing.T) {
	serverMock := new(mockConceptSearchAPI)
	requestedUUIDs := []string{uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestSearchResponseInvalidJSON", requestedUUIDs)
	serverMock.On("getResponse").Return(`{`, http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	search := NewSearch(&http.Client{}, server.URL)
	_, err := search.ByIDs("tid_TestSearchResponseInvalidJSON", requestedUUIDs...)

	assert.Error(t, err)
	serverMock.AssertExpectations(t) // failure here means the search API has not been called
}

type mockConceptSearchAPI struct {
	mock.Mock
}

func (m *mockConceptSearchAPI) getRequest() (string, []string) {
	args := m.Called()
	return args.String(0), args.Get(1).([]string)
}

func (m *mockConceptSearchAPI) getResponse() (string, int) {
	args := m.Called()
	return args.String(0), args.Int(1)
}

func (m *mockConceptSearchAPI) startServer(t *testing.T) *httptest.Server {
	r := vestigo.NewRouter()
	r.Get("/concepts", func(w http.ResponseWriter, r *http.Request) {
		tid, expectedIDs := m.getRequest()

		assert.Equal(t, tid, r.Header.Get("X-Request-Id"))
		assert.Equal(t, expectedUserAgent, r.Header.Get("User-Agent"))

		query := r.URL.Query()
		actualIDs, found := query[conceptSearchQueryParam]
		assert.True(t, found)
		assert.Equal(t, expectedIDs, actualIDs)

		json, status := m.getResponse()
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(json))
	})

	return httptest.NewServer(r)
}

func TestSearchHappyCheck(t *testing.T) {
	gtgServerMock := newConceptSearchAPIGTGMock(t, http.StatusOK)
	defer gtgServerMock.Close()

	search := NewSearch(&http.Client{}, gtgServerMock.URL)
	check := search.Check()
	assertSearchCheckConsistency(t, check)
	msg, err := check.Checker()
	assert.NoError(t, err)
	assert.Equal(t, "Concept Search API is good to go", msg)
}

func TestSearchUnhappyCheckDueInvalidURL(t *testing.T) {
	search := NewSearch(&http.Client{}, ":#")
	check := search.Check()
	assertSearchCheckConsistency(t, check)
	_, err := check.Checker()
	assert.EqualError(t, err, "parse \":\": missing protocol scheme")
}

func TestSearchUnhappyCheckDueHTTPCallError(t *testing.T) {
	search := NewSearch(&http.Client{}, "")
	check := search.Check()
	assertSearchCheckConsistency(t, check)
	_, err := check.Checker()
	assert.EqualError(t, err, "Get \"/__gtg\": unsupported protocol scheme \"\"")
}

func TestSearchUnhappyCheckDueNon200HTTPStatus(t *testing.T) {
	gtgServerMock := newConceptSearchAPIGTGMock(t, http.StatusServiceUnavailable)
	defer gtgServerMock.Close()

	search := NewSearch(&http.Client{}, gtgServerMock.URL)
	check := search.Check()
	assertSearchCheckConsistency(t, check)
	_, err := check.Checker()
	assert.EqualError(t, err, "GTG returned a non-200 HTTP status: 503")
}

func assertSearchCheckConsistency(t *testing.T, check fthealth.Check) {
	assert.Equal(t, "concept-search-api", check.ID)
	assert.Equal(t, "Concept information can not be returned to clients", check.BusinessImpact)
	assert.Equal(t, "Concept Search API Healthcheck", check.Name)
	assert.Equal(t, "https://runbooks.in.ft.com/internal-concordances", check.PanicGuide)
	assert.Equal(t, uint8(2), check.Severity)
	assert.Equal(t, "Concept Search API is not available", check.TechnicalSummary)
}

func newConceptSearchAPIGTGMock(t *testing.T, status int) *httptest.Server {
	r := vestigo.NewRouter()
	r.Get("/__gtg", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedUserAgent, r.Header.Get("User-Agent"))
		w.WriteHeader(status)
	})
	return httptest.NewServer(r)
}
