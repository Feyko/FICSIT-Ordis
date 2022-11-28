package base

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SearchableTestSuite struct {
	suite.Suite
	mod *Searchable[ExampleElement]
	rep repo.Repository[id.IDer]
}

func (s *SearchableTestSuite) SafeCreateElement(element ExampleElement) {
	err := s.mod.Create(nil, element)
	s.Require().NoErrorf(err, "Could not create the element: %v", err)
}

func (s *SearchableTestSuite) SafeCreateElements(elements []ExampleElement) {
	for _, cmd := range elements {
		err := s.mod.Create(nil, cmd)
		s.Require().NoErrorf(err, "Could not create an element: %v", err)
	}
}

func (s *SearchableTestSuite) SetupSuite() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)
	s.rep = rep
}

func (s *SearchableTestSuite) SetupTest() {
	collection, err := repos.CreateCollection[ExampleElement](s.rep, "Searchable")
	s.Require().NoError(err)
	authModule, err := auth.New(auth.Config{Secret: "test-secret"}, s.rep)
	s.Require().NoError(err)
	s.mod = NewSearchable[ExampleElement](NewDefaultConfigNoPerm(authModule), collection)
}

func (s *SearchableTestSuite) TearDownTest() {
	err := s.rep.DeleteCollection("Searchable")
	s.Require().NoError(err)
}

var defaultResponseToSearchText = "SearchMe"
var defaultResponseToSearch = domain.Response{Text: &defaultResponseToSearchText}

func (s *SearchableTestSuite) TestSearchValid() {
	input := []ExampleElement{
		{"Name", defaultResponseToSearch, "Media"},
		{"SearchMe", defaultResponse, "Media"},
		{"UniqueName", defaultResponse, "DONTSearchMe"},
	}
	expected := []ExampleElement{
		{"Name", defaultResponseToSearch, "Media"},
		{"SearchMe", defaultResponse, "Media"},
	}

	s.SafeCreateElements(input)
	actual, err := s.mod.Search(nil, "SearchMe")
	s.Require().NoErrorf(err, "Could not search for elements: %v", err)
	s.Equal(expected, actual)
}

func TestSearchableSuite(t *testing.T) {
	suite.Run(t, new(SearchableTestSuite))
}
