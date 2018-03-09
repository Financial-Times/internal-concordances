package concepts

import (
	"encoding/json"
	"fmt"
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
)

const (
	concordancesQueryParam    = "conceptId"
	authorityQueryParam       = "authority"
	identifierValueQueryParam = "identifierValue"
	NoAuthority               = ""
)

type Concordances interface {
	GetConcordances(tid, authority string, ids ...string) (map[string][]Identifier, error)
	Check() fthealth.Check
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

func (c *publicConcordancesAPI) GetConcordances(tid, authority string, ids ...string) (map[string][]Identifier, error) {
	if err := validateIDs(ids); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.uri+"/concordances", nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	reqParamName := concordancesQueryParam

	if authority != NoAuthority {
		queryParams.Add(authorityQueryParam, authority)
		reqParamName = identifierValueQueryParam
	}

	for _, id := range ids {
		queryParams.Add(reqParamName, id)
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

func (c *publicConcordancesAPI) Check() fthealth.Check {
	return fthealth.Check{
		ID:               "public-concordance-api",
		BusinessImpact:   "Concorded concepts can not be returned to clients",
		Name:             "Public Concordance API Healthcheck",
		PanicGuide:       "https://dewey.ft.com/internal-concordances.html",
		Severity:         2,
		TechnicalSummary: "Public Concordance API is not available",
		Checker:          c.gtg,
	}
}

func (c *publicConcordancesAPI) gtg() (string, error) {
	req, err := http.NewRequest("GET", c.uri+"/__gtg", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "UPP internal-concordances")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GTG returned a non-200 HTTP status: %v", resp.StatusCode)
	}
	return "Public Concordance API is good to go", nil
}
