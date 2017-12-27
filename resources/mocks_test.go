package resources

import (
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
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

func (m *mockConcordances) Check() fthealth.Check {
	args := m.Called()
	return args.Get(0).(fthealth.Check)
}

type mockSearch struct {
	mock.Mock
}

func (m *mockSearch) ByIDs(tid string, uuids ...string) (map[string]concepts.Concept, error) {
	args := m.Called(tid, uuids)
	return args.Get(0).(map[string]concepts.Concept), args.Error(1)
}

func (m *mockSearch) Check() fthealth.Check {
	args := m.Called()
	return args.Get(0).(fthealth.Check)
}
