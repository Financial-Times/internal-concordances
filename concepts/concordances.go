package concepts

import (
	"encoding/json"
	"net/http"
)

const concordancesQueryParam = "conceptId"

type Concordances interface {
	GetConcordances(tid string, uuids ...string) (map[string][]Identifier, error)
}

type publicConcordancesAPI struct {
	client *http.Client
	uri    string
}

type publicConcordancesResponse struct {
	Concordances []Concordance `json:"concordances"`
}

func NewConcordances(client *http.Client, uri string) Concordances {
	return &publicConcordancesAPI{client: client, uri: uri}
}

func (c *publicConcordancesAPI) GetConcordances(tid string, uuids ...string) (map[string][]Identifier, error) {
	if err := validateUUIDs(uuids); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.uri, nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	for _, uuid := range uuids {
		queryParams.Add(concordancesQueryParam, uuid)
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

	concordances := publicConcordancesResponse{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&concordances)

	if err != nil {
		return nil, err
	}

	return concordancesToIdentifiers(concordances.Concordances), nil
}

func concordancesToIdentifiers(concordances []Concordance) map[string][]Identifier {
	identifiers := make(map[string][]Identifier)
	for _, concordance := range concordances {
		if uuid, ok := stripThingPrefix(concordance.Concept.ID); ok {
			identifiers[uuid] = append(identifiers[uuid], concordance.Identifier)
		}
	}
	return identifiers
}
