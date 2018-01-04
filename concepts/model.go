package concepts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Concept struct {
	ID           string       `json:"id"`
	APIURL       string       `json:"apiUrl,omitempty"`
	Type         string       `json:"type,omitempty"`
	PrefLabel    string       `json:"prefLabel,omitempty"`
	Concordances []Identifier `json:"concordances,omitempty"`
	IsFTAuthor   *bool        `json:"isFTAuthor,omitempty"`
}

type Identifier struct {
	IdentifierValue string `json:"identifierValue"`
	Authority       string `json:"authority"`
}

type Concordance struct {
	Concept    Concept    `json:"concept"`
	Identifier Identifier `json:"identifier"`
}

type ResponseError struct {
	Status  string
	Message string `json:"message"`
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("%v: %s", r.Status, r.Message)
}

func decodeResponseError(resp *http.Response) error {
	err := ResponseError{Status: resp.Status}
	dec := json.NewDecoder(resp.Body)
	decodeErr := dec.Decode(&err)
	if decodeErr != nil {
		err.Message = "Failed to decode message from response"
	}
	return err
}

const apiIDPrefix = "http://api.ft.com/things/"
const ftIDPrefix = "http://www.ft.com/thing/"

func stripThingPrefix(conceptID string) (string, bool) {
	if strings.HasPrefix(conceptID, apiIDPrefix) {
		return strings.TrimPrefix(conceptID, apiIDPrefix), true
	}

	if strings.HasPrefix(conceptID, ftIDPrefix) {
		return strings.TrimPrefix(conceptID, ftIDPrefix), true
	}

	return "", false
}
