package commands

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/ports/repos/memrepo"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CommandsTestSuite struct {
	suite.Suite
	mod *Module
}

func (s *CommandsTestSuite) SafeCreateCommand(cmd Command) {
	err := s.mod.Create(cmd)
	s.Require().NoErrorf(err, "Could not create the command: %v", err)
}

func (s *CommandsTestSuite) SafeCreateCommands(cmds []Command) {
	for _, cmd := range cmds {
		err := s.mod.Create(cmd)
		s.Require().NoErrorf(err, "Could not create a command: %v", err)
	}
}

func (s *CommandsTestSuite) SetupTest() {
	repo := memrepo.New()
	mod, err := New(config.CommandsConfig{}, repo)
	s.Require().NoErrorf(err, "Could not create the module: %v", err)
	s.mod = mod
}

func (s *CommandsTestSuite) TestSearchValid() {
	input := []Command{
		{"Name", "SearchMe", "Media"},
		{"SearchMe", "Response", "Media"},
		{"UniqueName", "Response", "DONTSearchMe"},
	}
	expected := []Command{
		{"Name", "SearchMe", "Media"},
		{"SearchMe", "Response", "Media"},
	}
	s.SafeCreateCommands(input)
	actual, err := s.mod.Search("SearchMe")
	s.Require().NoErrorf(err, "Could not search for commands: %v", err)
	s.Equal(expected, actual)
}

func TestCommandsSuite(t *testing.T) {
	suite.Run(t, new(CommandsTestSuite))
}
