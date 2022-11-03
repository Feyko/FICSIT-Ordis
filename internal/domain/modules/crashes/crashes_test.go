package crashes

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
	suite.Run(t, new(CrashesModuleTestSuite))
}

type CrashesModuleTestSuite struct {
	suite.Suite
	rep repo.Repository[domain.Crash]
	mod *Module
}

func (s *CrashesModuleTestSuite) SetupSuite() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)
	rep, err = repos.Retype[domain.Crash, id.IDer](rep)
	s.Require().NoError(err)
	s.rep = rep
}

func (s *CrashesModuleTestSuite) SetupTest() {
	s.setupTest(true)
}

func (s *CrashesModuleTestSuite) setupTest(noAuth bool) {
	authModule, err := auth.New(auth.Config{Secret: "test-secret"})
	s.Require().NoError(err)

	mod, err := New(
		Config{
			auth.AuthedConfig{
				noAuth,
				authModule,
			},
		}, s.rep)
	s.Require().NoError(err)

	s.mod = mod

	err = s.mod.Create(nil, defaultCrash)
	s.Require().NoError(err)
}

var defaultCrash = domain.Crash{
	Name:    "default",
	Description: "Default description",
	Response: &domain.Response{
		Text:       "Text",
		MediaLinks: []string{"https://SomeLink"},
	},
	Regexes: []string{"default"},
}

func (s *CrashesModuleTestSuite) TearDownTest() {
	err := s.rep.DeleteCollection("Crashes")
	s.Require().NoError(err)
}

func (s *CrashesModuleTestSuite) TestGetByAlias() {
	_, err := s.mod.Get(nil, "defaultalias")
	s.Require().NoError(err)
}