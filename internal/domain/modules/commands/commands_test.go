package commands

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCommandsModuleTestSuite(t *testing.T) {
	suite.Run(t, new(CommandsModuleTestSuite))
}

type CommandsModuleTestSuite struct {
	suite.Suite
	rep repo.Repository[id.IDer]
	mod *Module
}

func (s *CommandsModuleTestSuite) SetupSuite() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)
	s.rep = rep
}

func (s *CommandsModuleTestSuite) SetupTest() {
	s.setupTest(true)
}

func (s *CommandsModuleTestSuite) setupTest(noAuth bool) {
	authModule, err := auth.New(auth.Config{Secret: "test-secret"}, s.rep)
	s.Require().NoError(err)

	mod, err := New(
		Config{
			AuthedConfig: auth.AuthedConfig{
				AuthModule: authModule,
				NoAuth:     noAuth,
			},
		}, s.rep)
	s.Require().NoError(err)

	s.mod = mod

	err = s.mod.Create(nil, defaultCommand)
	s.Require().NoError(err)
}

var defaultText = "Text"
var defaultCommand = domain.Command{
	Name:    "default",
	Aliases: []string{"defaultalias"},
	Response: domain.Response{
		Text:       &defaultText,
		MediaLinks: []string{"https://SomeLink"},
	},
}

func (s *CommandsModuleTestSuite) TearDownTest() {
	err := s.rep.DeleteCollection("Commands")
	s.Require().NoError(err)
}

func (s *CommandsModuleTestSuite) TestGetByAlias() {
	_, err := s.mod.Get(nil, "defaultalias")
	s.Require().NoError(err)
}

func (s *CommandsModuleTestSuite) TestCreateUniqueAliasName() {
	err := s.mod.Create(nil, domain.Command{
		Name:    "uniquename",
		Aliases: []string{"defaultalias"},
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

func (s *CommandsModuleTestSuite) TestCreateCommandsSimilarName() {
	err := s.mod.Create(nil, domain.Command{
		Name: "ayo",
	})
	s.Require().NoError(err)
	err = s.mod.Create(nil, domain.Command{
		Name: "ayooo",
	})
	s.Require().NoError(err)
}

func (s *CommandsModuleTestSuite) TestDeleteByAlias() {
	err := s.mod.Delete(nil, "defaultalias")
	s.Require().NoError(err)
}

func (s *CommandsModuleTestSuite) TestUpdateByAlias() {
	expected := &defaultCommand
	expected.Name = "newname"
	updated, err := s.mod.Update(nil, "defaultalias", domain.CommandUpdate{Name: &expected.Name})
	s.Require().NoError(err)
	s.Require().Equal(expected, updated)
}
