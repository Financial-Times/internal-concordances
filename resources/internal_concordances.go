package resources

import (
	"encoding/json"
	"net/http"

	"github.com/Financial-Times/internal-concordances/concepts"
	tidutils "github.com/Financial-Times/transactionid-utils-go"
)

type internalConcordancesResponse struct {
	Concepts map[string]concepts.Concept `json:"concepts"`
}

// InternalConcordances concords provided uuids, and enriches them with concept model
func InternalConcordances(concordances concepts.Concordances, search concepts.Search) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		tid := tidutils.GetTransactionIDFromRequest(req)

		ids, found := getMultiValuedParam(req, "ids")
		if !found {
			writeJSON("Please provide ids to concord, using the 'ids' query parameter", http.StatusBadRequest, w)
			return
		}

		identifiers, err := concordances.GetConcordances(tid, ids...)
		if err == concepts.ErrConceptUUIDsAreEmpty {
			writeJSON("Please provide non-empty ids to concord, using the 'ids' query parameter", http.StatusBadRequest, w)
			return
		}

		if err != nil {
			writeJSON("Public Concordances request failed, please try again", http.StatusServiceUnavailable, w)
			return
		}

		if len(identifiers) == 0 { // all requested concepts were either deleted or missing
			writeInternalConcordanceResponse(w, internalConcordancesResponse{Concepts: make(map[string]concepts.Concept)})
			return
		}

		concordedUUIDs := conceptIdentifiersToUUIDs(identifiers)
		concepts, err := search.ByIDs(tid, concordedUUIDs...)
		if err != nil {
			writeJSON("Concept Search request failed, please try again", http.StatusServiceUnavailable, w)
			return
		}

		merged := mergeConcordancesAndConcepts(ids, identifiers, concepts)
		resp := internalConcordancesResponse{Concepts: merged}

		writeInternalConcordanceResponse(w, resp)
	}
}

func writeInternalConcordanceResponse(w http.ResponseWriter, resp internalConcordancesResponse) {
	jsonBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func mergeConcordancesAndConcepts(requestedUUIDs []string, identifiers map[string][]concepts.Identifier, searchedConcepts map[string]concepts.Concept) map[string]concepts.Concept {
	merged := make(map[string]concepts.Concept)

	for uuid, concept := range searchedConcepts {
		concordances := identifiers[uuid]
		concept.Concordances = concordances

		for _, c := range concordances {
			for _, requestedUUID := range requestedUUIDs {
				if c.IdentifierValue == requestedUUID {
					merged[requestedUUID] = concept
				}
			}
		}
	}

	return merged
}

func conceptIdentifiersToUUIDs(identifiers map[string][]concepts.Identifier) []string {
	uuids := make([]string, 0)
	for uuid := range identifiers {
		uuids = append(uuids, uuid)
	}
	return uuids
}

func writeJSON(msg string, status int, w http.ResponseWriter) {
	resp := make(map[string]string)
	resp["message"] = msg

	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func getMultiValuedParam(req *http.Request, param string) ([]string, bool) {
	query := req.URL.Query()
	values, found := query[param]
	return values, found
}
