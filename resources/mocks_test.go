package resources

import (
	"github.com/Financial-Times/internal-concordances/concepts"
	"github.com/stretchr/testify/mock"
)

type mockConcordances struct {
	mock.Mock
}

func (m *mockConcordances) GetConcordances(tid string, uuids ...string) (map[string][]concepts.Identifier, error) {
	args := m.Called(tid, uuids)
	return args.Get(0).(map[string][]concepts.Identifier), args.Error(1)
}

type mockSearch struct {
	mock.Mock
}

func (m *mockSearch) ByIDs(tid string, uuids ...string) (map[string]concepts.Concept, error) {
	args := m.Called(tid, uuids)
	return args.Get(0).(map[string]concepts.Concept), args.Error(1)
}
