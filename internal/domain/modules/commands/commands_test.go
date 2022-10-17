package commands

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

func TestExampleModuleTestSuite(t *testing.T) {
	suite.Run(t, new(CommandsModuleTestSuite))
}

type CommandsModuleTestSuite struct {
	suite.Suite
	rep repo.Repository[domain.Command]
	mod *Module
}

func (s *CommandsModuleTestSuite) SetupSuite() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)
	rep, err = repos.Retype[domain.Command, id.IDer](rep)
	s.Require().NoError(err)
	s.rep = rep
}

func (s *CommandsModuleTestSuite) SetupTest() {
	s.setupTest(true)
}

func (s *CommandsModuleTestSuite) setupTest(noAuth bool) {
	authModule, err := auth.New(auth.Config{Secret: "test-secret"})
	s.Require().NoError(err)

	mod, err := New(
		Config{
			Auth:   authModule,
			NoAuth: noAuth,
		}, s.rep)
	s.Require().NoError(err)

	s.mod = mod

	err = s.mod.Create(nil, defaultCommand)
	s.Require().NoError(err)
}

var defaultCommand = domain.Command{
	Name:    "default",
	Aliases: nil,
	Response: domain.Response{
		Text:       "Text",
		MediaLinks: []string{"https://SomeLink"},
	},
}

func (s *CommandsModuleTestSuite) TearDownTest() {
	err := s.rep.DeleteCollection("Commands")
	s.Require().NoError(err)
}

func (s *CommandsModuleTestSuite) TestCreateUniqueAliasName() {
	err := s.mod.Create(nil, domain.Command{
		Name:    "default",
		Aliases: []string{"uniquealias"},
	})
	s.Require().Error(err)
}

func (s *CommandsModuleTestSuite) TestCreateUniqueAliasAlias() {
	err := s.mod.Create(nil, domain.Command{
		Name:    "uniquename",
		Aliases: []string{"default"},
	})
	s.Require().Error(err)
}
