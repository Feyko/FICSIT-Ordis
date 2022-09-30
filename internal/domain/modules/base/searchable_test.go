package base

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type SearchableTestSuite struct {
	suite.Suite
	mod *Searchable[ExampleElement]
}

func (s *SearchableTestSuite) SafeCreateElement(element ExampleElement) {
	err := s.mod.Create(element)
	s.Require().NoErrorf(err, "Could not create the element: %v", err)
}

func (s *SearchableTestSuite) SafeCreateElements(elements []ExampleElement) {
	for _, cmd := range elements {
		err := s.mod.Create(cmd)
		s.Require().NoErrorf(err, "Could not create an element: %v", err)
	}
}

func (s *SearchableTestSuite) SetupTest() {
	mod, err := newDefaultSearchable[ExampleElement]()
	if err != nil {
		s.Fail("Error when setting up: %+v", err)
	}
	s.mod = mod
}

func (s *SearchableTestSuite) TestSearchValid() {
	input := []ExampleElement{
		{"Name", "SearchMe", "Media"},
		{"SearchMe", "Response", "Media"},
		{"UniqueName", "Response", "DONTSearchMe"},
	}
	expected := []ExampleElement{
		{"Name", "SearchMe", "Media"},
		{"SearchMe", "Response", "Media"},
	}

	s.SafeCreateElements(input)
	actual, err := s.mod.Search("SearchMe")
	s.Require().NoErrorf(err, "Could not search for elements: %v", err)
	s.Equal(expected, actual)
}

func TestSearchableSuite(t *testing.T) {
	suite.Run(t, new(SearchableTestSuite))
}
