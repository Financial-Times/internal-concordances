package resources

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Financial-Times/internal-concordances/concepts"
	"github.com/stretchr/testify/assert"
)

var errComputerSaysNo = errors.New("computer says no")

func TestInternalConcordancesNoIDsParamSupplied(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	InternalConcordances(nil, nil)(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"Please provide ids to concord, using the 'ids' query parameter"}`, strings.TrimSpace(w.Body.String()))
}

func TestGetConcordancesFailsDueToEmptyUUIDS(t *testing.T) {
	concordances := new(mockConcordances)
	search := new(mockSearch)

	req := httptest.NewRequest("GET", "/?ids=", nil)
	req.Header.Add("X-Request-Id", "tid_TestGetConcordancesFailsDueToEmptyUUIDS")
	w := httptest.NewRecorder()

	concordances.On("GetConcordances", "tid_TestGetConcordancesFailsDueToEmptyUUIDS", []string{""}).
		Return(make(map[string][]concepts.Identifier), concepts.ErrConceptUUIDsAreEmpty)

	InternalConcordances(concordances, search)(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"Please provide non-empty ids to concord, using the 'ids' query parameter"}`, strings.TrimSpace(w.Body.String()))

	concordances.AssertExpectations(t)
	search.AssertExpectations(t)
}

func TestGetConcordancesFails(t *testing.T) {
	concordances := new(mockConcordances)
	search := new(mockSearch)

	req := httptest.NewRequest("GET", "/?ids=a-uuid", nil)
	req.Header.Add("X-Request-Id", "tid_TestGetConcordancesFails")
	w := httptest.NewRecorder()

	concordances.On("GetConcordances", "tid_TestGetConcordancesFails", []string{"a-uuid"}).
		Return(make(map[string][]concepts.Identifier), errComputerSaysNo)

	InternalConcordances(concordances, search)(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, `{"message":"Public Concordances request failed, please try again"}`, strings.TrimSpace(w.Body.String()))

	concordances.AssertExpectations(t)
	search.AssertExpectations(t)
}

func TestGetConcordancesReturnsNoData(t *testing.T) {
	concordances := new(mockConcordances)
	search := new(mockSearch)

	req := httptest.NewRequest("GET", "/?ids=a-uuid", nil)
	req.Header.Add("X-Request-Id", "tid_TestGetConcordancesReturnsNoData")
	w := httptest.NewRecorder()

	concordances.On("GetConcordances", "tid_TestGetConcordancesReturnsNoData", []string{"a-uuid"}).
		Return(make(map[string][]concepts.Identifier), nil)

	InternalConcordances(concordances, search)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"concepts":{}}`, strings.TrimSpace(w.Body.String()))

	concordances.AssertExpectations(t)
	search.AssertExpectations(t)
}

func TestSearchByIDsFails(t *testing.T) {
	concordances := new(mockConcordances)
	search := new(mockSearch)

	req := httptest.NewRequest("GET", "/?ids=a-uuid", nil)
	req.Header.Add("X-Request-Id", "tid_TestSearchByIDsFails")
	w := httptest.NewRecorder()

	identifiers := map[string][]concepts.Identifier{
		"a-uuid": []concepts.Identifier{
			{Authority: "authority", IdentifierValue: "a-uuid"},
		},
	}

	concordances.On("GetConcordances", "tid_TestSearchByIDsFails", []string{"a-uuid"}).
		Return(identifiers, nil)

	search.On("ByIDs", "tid_TestSearchByIDsFails", []string{"a-uuid"}).Return(make(map[string]concepts.Concept), errComputerSaysNo)

	InternalConcordances(concordances, search)(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, `{"message":"Concept Search request failed, please try again"}`, strings.TrimSpace(w.Body.String()))

	concordances.AssertExpectations(t)
	search.AssertExpectations(t)
}

func TestSearchByIDs(t *testing.T) {
	concordances := new(mockConcordances)
	search := new(mockSearch)

	req := httptest.NewRequest("GET", "/?ids=a-concorded-uuid", nil)
	req.Header.Add("X-Request-Id", "tid_TestSearchByIDs")
	w := httptest.NewRecorder()

	identifiers := map[string][]concepts.Identifier{
		"a-uuid": []concepts.Identifier{
			{Authority: "authority", IdentifierValue: "a-concorded-uuid"},
			{Authority: "authority", IdentifierValue: "a-uuid"},
		},
	}

	concordances.On("GetConcordances", "tid_TestSearchByIDs", []string{"a-concorded-uuid"}).
		Return(identifiers, nil)

	expectedConcepts := map[string]concepts.Concept{
		"a-uuid": {ID: "http://www.ft.com/thing/a-uuid", PrefLabel: "Donald Trump"},
	}

	expectedResponse := internalConcordancesResponse{Concepts: map[string]concepts.Concept{
		"a-concorded-uuid": {
			ID:        "http://www.ft.com/thing/a-uuid",
			PrefLabel: "Donald Trump",
			Concordances: []concepts.Identifier{
				{Authority: "authority", IdentifierValue: "a-concorded-uuid"},
				{Authority: "authority", IdentifierValue: "a-uuid"},
			},
		},
	}}

	search.On("ByIDs", "tid_TestSearchByIDs", []string{"a-uuid"}).
		Return(expectedConcepts, nil)

	InternalConcordances(concordances, search)(w, req)

	b, _ := json.Marshal(expectedResponse)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(b), w.Body.String())

	concordances.AssertExpectations(t)
	search.AssertExpectations(t)
}

func TestSearchByIDsOneConceptNotFound(t *testing.T) {
	concordances := new(mockConcordances)
	search := new(mockSearch)

	req := httptest.NewRequest("GET", "/?ids=found-this-one&ids=but-not-this-one", nil)
	req.Header.Add("X-Request-Id", "tid_TestSearchByIDsOneConceptNotFound")
	w := httptest.NewRecorder()

	identifiers := map[string][]concepts.Identifier{
		"found-this-one": []concepts.Identifier{
			{Authority: "authority", IdentifierValue: "found-this-one"},
		},
	}

	concordances.On("GetConcordances", "tid_TestSearchByIDsOneConceptNotFound", []string{"found-this-one", "but-not-this-one"}).
		Return(identifiers, nil)

	expectedConcepts := map[string]concepts.Concept{
		"found-this-one": {ID: "http://www.ft.com/thing/found-this-one", PrefLabel: "Donald Trump"},
	}

	expectedResponse := internalConcordancesResponse{Concepts: map[string]concepts.Concept{
		"found-this-one": {
			ID:        "http://www.ft.com/thing/found-this-one",
			PrefLabel: "Donald Trump",
			Concordances: []concepts.Identifier{
				{Authority: "authority", IdentifierValue: "found-this-one"},
			},
		},
	}}

	search.On("ByIDs", "tid_TestSearchByIDsOneConceptNotFound", []string{"found-this-one"}).
		Return(expectedConcepts, nil)

	InternalConcordances(concordances, search)(w, req)

	b, _ := json.Marshal(expectedResponse)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(b), w.Body.String())

	concordances.AssertExpectations(t)
	search.AssertExpectations(t)
}
