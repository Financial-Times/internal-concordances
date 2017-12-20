package concepts

import (
	"encoding/json"
	"errors"
	"net/http"
)

const conceptSearchQueryParam = "ids"

var (
	// ErrNoConceptsToSearch indicates the provided uuids array was empty
	ErrNoConceptsToSearch = errors.New("no concept ids to search for")
	// ErrConceptUUIDsAreEmpty indicates the provided uuids array only contained empty string
	ErrConceptUUIDsAreEmpty = errors.New("provided concept ids are empty")
)

type Search interface {
	ByIDs(tid string, uuids ...string) (map[string]Concept, error)
}

type conceptSearchAPI struct {
	client *http.Client
	uri    string
}

type conceptSearchResponse struct {
	Concepts []Concept `json:"concepts"`
}

func NewSearch(client *http.Client, uri string) Search {
	return &conceptSearchAPI{client: client, uri: uri}
}

func (c *conceptSearchAPI) ByIDs(tid string, uuids ...string) (map[string]Concept, error) {
	if err := validateUUIDs(uuids); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.uri, nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	for _, uuid := range uuids {
		queryParams.Add(conceptSearchQueryParam, uuid)
	}

	req.URL.RawQuery = queryParams.Encode()

	stampRequest(req, tid)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeResponseError(resp)
	}

	searchResp := conceptSearchResponse{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&searchResp)

	if err != nil {
		return nil, err
	}

	concepts := make(map[string]Concept)
	for _, c := range searchResp.Concepts {
		if uuid, ok := stripThingPrefix(c.ID); ok {
			concepts[uuid] = c
		}
	}

	return concepts, nil
}

func stampRequest(req *http.Request, tid string) {
	req.Header.Add("User-Agent", "UPP internal-concordances")
	req.Header.Add("X-Request-Id", tid)
}

func validateUUIDs(uuids []string) error {
	if len(uuids) == 0 {
		return ErrNoConceptsToSearch
	}

	for _, v := range uuids {
		if v != "" {
			return nil
		}
	}

	return ErrConceptUUIDsAreEmpty
}
