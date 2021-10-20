package concepts

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	uuid "github.com/google/uuid"
	"github.com/husobee/vestigo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const expectedUserAgent = "UPP internal-concordances"

func TestGetConcordancesEmptyResponse(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{uuid.New().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesEmptyResponse", requestedUUIDs)
	serverMock.On("getResponse").Return(`{}`, http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL)
	identifiers, err := concordances.GetConcordances("tid_TestGetConcordancesEmptyResponse", NoAuthority, requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, identifiers, 0)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesAtLeastOneNonEmptyID(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.New().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesAtLeastOneNonEmptyID", requestedUUIDs)
	serverMock.On("getResponse").Return(`{}`, http.StatusOK) // respond with an empty body, so no data will be returned, but if the test passes then the test case is successful.

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL)
	identifiers, err := concordances.GetConcordances("tid_TestGetConcordancesAtLeastOneNonEmptyID", NoAuthority, requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, identifiers, 0)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordances(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{uuid.New().String()}
	serverMock.On("getRequest").Return("tid_TestGetEmptyConcordances", requestedUUIDs)

	concordancesResp, err := ioutil.ReadFile("./_fixtures/concordances_response.json")
	require.NoError(t, err)
	serverMock.On("getResponse").Return(string(concordancesResp), http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL)
	identifiers, err := concordances.GetConcordances("tid_TestGetEmptyConcordances", NoAuthority, requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, identifiers, 1)

	ids, ok := identifiers["2753c50c-b256-4814-9f0d-65c8e755aa14"]
	assert.True(t, ok)
	assert.Len(t, ids, 4)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesByAuthority(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{uuid.New().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesByAuthority", requestedUUIDs)

	concordancesResp, err := ioutil.ReadFile("./_fixtures/concordances_by_authority_response.json")
	require.NoError(t, err)
	serverMock.On("getResponse").Return(string(concordancesResp), http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL)
	uppAuthority := "http://api.ft.com/system/UPP"
	identifiers, err := concordances.GetConcordances("tid_TestGetConcordancesByAuthority", uppAuthority, requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, identifiers, 1)

	ids, ok := identifiers["2753c50c-b256-4814-9f0d-65c8e755aa14"]
	assert.True(t, ok)
	assert.Len(t, ids, 1)
	assert.Equal(t, uppAuthority, ids[0].Authority)
	assert.Equal(t, "6b43a14b-a5e0-3b63-a428-aa55def05fcb", ids[0].IdentifierValue)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesFailsWhenNoIDsSupplied(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, "")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailsWhenNoIDsSupplied", NoAuthority)

	assert.EqualError(t, err, ErrNoConceptsToSearch.Error())
}

func TestGetConcordancesFailsWhenEmptyIDsSupplied(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, "")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailsWhenEmptyIDsSupplied", NoAuthority, "", "")

	assert.EqualError(t, err, ErrConceptIDsAreEmpty.Error())
}

func TestGetConcordancesFailsInvalidURL(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, ":#") // this triggers a invalid url during the http.NewRequest() line
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailsInvalidURL", NoAuthority, uuid.New().String())

	assert.Error(t, err)
}

func TestGetConcordancesRequestFails(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, "#:") // this triggers a protocol error in the client.Do()
	_, err := concordances.GetConcordances("tid_TestGetConcordancesRequestFails", NoAuthority, uuid.New().String())

	assert.Error(t, err)
}

func TestGetConcordancesResponseJSONInvalid(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.New().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesResponseJSONInvalid", requestedUUIDs)
	serverMock.On("getResponse").Return(`{`, http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL)
	_, err := concordances.GetConcordances("tid_TestGetConcordancesResponseJSONInvalid", NoAuthority, requestedUUIDs...)

	assert.Error(t, err)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesFailedResponse(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.New().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesFailedResponse", requestedUUIDs)
	serverMock.On("getResponse").Return(`{"message":"uh oh"}`, http.StatusServiceUnavailable)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL)
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailedResponse", NoAuthority, requestedUUIDs...)

	assert.EqualError(t, err, "503 Service Unavailable: uh oh")
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesFailedResponseMessageDecodingAlsoFailed(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.New().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesFailedResponse", requestedUUIDs)
	serverMock.On("getResponse").Return(`{`, http.StatusBadRequest)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL)
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailedResponse", NoAuthority, requestedUUIDs...)

	assert.EqualError(t, err, "400 Bad Request: Failed to decode message from response")
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

type mockPublicConcordancesServer struct {
	mock.Mock
}

func (m *mockPublicConcordancesServer) getRequest() (string, []string) {
	args := m.Called()
	return args.String(0), args.Get(1).([]string)
}

func (m *mockPublicConcordancesServer) getResponse() (string, int) {
	args := m.Called()
	return args.String(0), args.Int(1)
}

func (m *mockPublicConcordancesServer) startServer(t *testing.T) *httptest.Server {
	r := vestigo.NewRouter()
	r.Get("/concordances", func(w http.ResponseWriter, r *http.Request) {
		tid, expectedIDs := m.getRequest()

		assert.Equal(t, tid, r.Header.Get("X-Request-Id"))
		assert.Equal(t, expectedUserAgent, r.Header.Get("User-Agent"))

		query := r.URL.Query()
		authorityParam, foundAuthority := query[authorityQueryParam]
		queryParam := concordancesQueryParam
		if foundAuthority {
			queryParam = identifierValueQueryParam
			assert.NotEmpty(t, authorityParam)
			assert.Empty(t, query[concordancesQueryParam])
		}
		actualIDs, foundConceptId := query[queryParam]
		assert.True(t, foundConceptId)
		assert.Equal(t, expectedIDs, actualIDs)

		json, status := m.getResponse()
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(json))
	})

	return httptest.NewServer(r)
}

func TestConcordanceHappyCheck(t *testing.T) {
	gtgServerMock := newPublicConcordanceAPIGTGMock(t, http.StatusOK)
	defer gtgServerMock.Close()

	search := NewConcordances(&http.Client{}, gtgServerMock.URL)
	check := search.Check()
	assertConcordanceCheckConsistency(t, check)
	msg, err := check.Checker()
	assert.NoError(t, err)
	assert.Equal(t, "Public Concordance API is good to go", msg)
}

func TestConcordanceUnhappyCheckDueInvalidURL(t *testing.T) {
	search := NewConcordances(&http.Client{}, ":#")
	check := search.Check()
	assertConcordanceCheckConsistency(t, check)
	_, err := check.Checker()
	var urlErr *url.Error
	assert.True(t, errors.As(err, &urlErr))
	assert.Equal(t, urlErr.Op, "parse")
}

func TestConcordanceUnhappyCheckDueHTTPCallError(t *testing.T) {
	search := NewConcordances(&http.Client{}, "")
	check := search.Check()
	assertConcordanceCheckConsistency(t, check)
	_, err := check.Checker()
	var urlErr *url.Error
	assert.True(t, errors.As(err, &urlErr))
	assert.Equal(t, urlErr.Op, "Get")
}

func TestConcordanceUnhappyCheckDueNon200HTTPStatus(t *testing.T) {
	gtgServerMock := newConceptSearchAPIGTGMock(t, http.StatusServiceUnavailable)
	defer gtgServerMock.Close()

	search := NewConcordances(&http.Client{}, gtgServerMock.URL)
	check := search.Check()
	assertConcordanceCheckConsistency(t, check)
	_, err := check.Checker()
	assert.EqualError(t, err, "GTG returned a non-200 HTTP status: 503")
}

func assertConcordanceCheckConsistency(t *testing.T, check fthealth.Check) {
	assert.Equal(t, "public-concordance-api", check.ID)
	assert.Equal(t, "Concorded concepts can not be returned to clients", check.BusinessImpact)
	assert.Equal(t, "Public Concordance API Healthcheck", check.Name)
	assert.Equal(t, "https://runbooks.in.ft.com/internal-concordances", check.PanicGuide)
	assert.Equal(t, uint8(2), check.Severity)
	assert.Equal(t, "Public Concordance API is not available", check.TechnicalSummary)
}

func newPublicConcordanceAPIGTGMock(t *testing.T, status int) *httptest.Server {
	r := vestigo.NewRouter()
	r.Get("/__gtg", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedUserAgent, r.Header.Get("User-Agent"))
		w.WriteHeader(status)
	})
	return httptest.NewServer(r)
}
