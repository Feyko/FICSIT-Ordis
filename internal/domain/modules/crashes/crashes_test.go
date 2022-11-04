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

var defaultText = "Text"
var defaultDescription = "Default description"
var defaultCrash = domain.Crash{
	Name:        "default",
	Description: &defaultDescription,
	Response: domain.Response{
		Text:       &defaultText,
		MediaLinks: []string{"https://link.com"},
	},
	Regexes: []string{"default"},
}

func (s *CrashesModuleTestSuite) TearDownTest() {
	err := s.rep.DeleteCollection("Crashes")
	s.Require().NoError(err)
}

func (s *CrashesModuleTestSuite) TestCreateInvalidRegex() {
	crash := defaultCrash
	crash.Regexes = []string{"[ invalid"}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateInvalidLink() {
	crash := defaultCrash
	crash.Response.MediaLinks = []string{"notalink"}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateNoRegex() {
	crash := defaultCrash
	crash.Regexes = []string{}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateEmptyResponse() {
	crash := defaultCrash
	crash.Response = domain.Response{}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateResponseEmptyText() {
	crash := defaultCrash
	empty := ""
	crash.Response = domain.Response{Text: &empty}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}
