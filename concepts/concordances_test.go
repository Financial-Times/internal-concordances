package concepts

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/husobee/vestigo"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const expectedUserAgent = "UPP internal-concordances"

func TestGetConcordancesEmptyResponse(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesEmptyResponse", requestedUUIDs)
	serverMock.On("getResponse").Return(`{}`, http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL+"/concordances")
	identifiers, err := concordances.GetConcordances("tid_TestGetConcordancesEmptyResponse", requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, identifiers, 0)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesAtLeastOneNonEmptyID(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesAtLeastOneNonEmptyID", requestedUUIDs)
	serverMock.On("getResponse").Return(`{}`, http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL+"/concordances")
	identifiers, err := concordances.GetConcordances("tid_TestGetConcordancesAtLeastOneNonEmptyID", requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, identifiers, 0)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordances(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestGetEmptyConcordances", requestedUUIDs)

	concordancesResp, err := ioutil.ReadFile("./_fixtures/concordances_response.json")
	require.NoError(t, err)
	serverMock.On("getResponse").Return(string(concordancesResp), http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL+"/concordances")
	identifiers, err := concordances.GetConcordances("tid_TestGetEmptyConcordances", requestedUUIDs...)

	assert.NoError(t, err)
	assert.Len(t, identifiers, 1)

	ids, ok := identifiers["2753c50c-b256-4814-9f0d-65c8e755aa14"]
	assert.True(t, ok)
	assert.Len(t, ids, 4)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesFailsWhenNoIDsSupplied(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, "/concordances")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailsWhenNoIDsSupplied")

	assert.EqualError(t, err, ErrNoConceptsToSearch.Error())
}

func TestGetConcordancesFailsWhenEmptyIDsSupplied(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, "/concordances")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailsWhenEmptyIDsSupplied", "", "")

	assert.EqualError(t, err, ErrConceptUUIDsAreEmpty.Error())
}

func TestGetConcordancesFailsInvalidURL(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, ":#")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailsInvalidURL", uuid.NewV4().String())

	assert.Error(t, err)
}

func TestGetConcordancesRequestFails(t *testing.T) {
	concordances := NewConcordances(&http.Client{}, "#:")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesRequestFails", uuid.NewV4().String())

	assert.Error(t, err)
}

func TestGetConcordancesResponseJSONInvalid(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesResponseJSONInvalid", requestedUUIDs)
	serverMock.On("getResponse").Return(`{`, http.StatusOK)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL+"/concordances")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesResponseJSONInvalid", requestedUUIDs...)

	assert.Error(t, err)
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesFailedResponse(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesFailedResponse", requestedUUIDs)
	serverMock.On("getResponse").Return(`{"message":"uh oh"}`, http.StatusServiceUnavailable)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL+"/concordances")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailedResponse", requestedUUIDs...)

	assert.EqualError(t, err, "503 Service Unavailable: uh oh")
	serverMock.AssertExpectations(t) // failure here means the concordances API has not been called
}

func TestGetConcordancesFailedResponseMessageDecodingAlsoFailed(t *testing.T) {
	serverMock := new(mockPublicConcordancesServer)
	requestedUUIDs := []string{"", "", uuid.NewV4().String()}
	serverMock.On("getRequest").Return("tid_TestGetConcordancesFailedResponse", requestedUUIDs)
	serverMock.On("getResponse").Return(`{`, http.StatusBadRequest)

	server := serverMock.startServer(t)
	defer server.Close()

	concordances := NewConcordances(&http.Client{}, server.URL+"/concordances")
	_, err := concordances.GetConcordances("tid_TestGetConcordancesFailedResponse", requestedUUIDs...)

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
		actualIDs, found := query[concordancesQueryParam]
		assert.True(t, found)
		assert.Equal(t, expectedIDs, actualIDs)

		json, status := m.getResponse()
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(json))
	})

	return httptest.NewServer(r)
}
