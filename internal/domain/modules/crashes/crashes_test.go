package crashes

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCrashesModuleTestSuite(t *testing.T) {
	suite.Run(t, new(CrashesModuleTestSuite))
}

// TODO: Add search tests
type CrashesModuleTestSuite struct {
	suite.Suite
	rep repo.Repository[id.IDer]
	mod *Module
}

func (s *CrashesModuleTestSuite) SetupSuite() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)
	s.rep = rep
}

func (s *CrashesModuleTestSuite) SetupTest() {
	s.setupTest(true)
}

func (s *CrashesModuleTestSuite) setupTest(noAuth bool) {
	authModule, err := auth.New(auth.Config{Secret: "test-secret"}, s.rep)
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

var otherCrash = domain.Crash{
	Name:        "other",
	Description: &defaultDescription,
	Response: domain.Response{
		Text:       &defaultText,
		MediaLinks: []string{"https://link.com"},
	},
	Regexes: []string{"other"},
}

func (s *CrashesModuleTestSuite) TearDownTest() {
	err := s.rep.DeleteCollection("Crashes")
	s.Require().NoError(err)
}

func (s *CrashesModuleTestSuite) TestCreate() {
	err := s.mod.Create(nil, otherCrash)
	s.Require().NoError(err)
}

func (s *CrashesModuleTestSuite) TestCreateInvalidRegex() {
	crash := otherCrash
	crash.Regexes = []string{"[ invalid"}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateInvalidLink() {
	crash := otherCrash
	crash.Response.MediaLinks = []string{"notalink"}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateNoRegex() {
	crash := otherCrash
	crash.Regexes = []string{}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateEmptyResponse() {
	crash := otherCrash
	crash.Response = domain.Response{}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestCreateResponseEmptyText() {
	crash := otherCrash
	empty := ""
	crash.Response = domain.Response{Text: &empty}

	err := s.mod.Create(nil, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestAnalyseNewRegex() {
	err := s.mod.Create(nil, otherCrash)
	s.Require().NoError(err)

	matches, err := s.mod.Analyse(nil, otherCrash.Regexes[0])
	s.Require().NoError(err)
	s.Require().Len(matches, 1)
}

func (s *CrashesModuleTestSuite) TestUpdateInvalidRegex() {
	var crash domain.CrashUpdate
	crash.Regexes = []string{"[ invalid"}

	_, err := s.mod.Update(nil, defaultCrash.Name, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestUpdateInvalidLink() {
	var crash domain.CrashUpdate
	crash.Response = &domain.Response{MediaLinks: []string{"notalink"}}

	_, err := s.mod.Update(nil, defaultCrash.Name, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestUpdateNoRegex() {
	var crash domain.CrashUpdate
	crash.Regexes = []string{}

	_, err := s.mod.Update(nil, defaultCrash.Name, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestUpdateEmptyResponse() {
	var crash domain.CrashUpdate
	crash.Response = &domain.Response{}

	_, err := s.mod.Update(nil, defaultCrash.Name, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestUpdateResponseEmptyText() {
	var crash domain.CrashUpdate
	empty := ""
	crash.Response = &domain.Response{Text: &empty}

	_, err := s.mod.Update(nil, defaultCrash.Name, crash)
	s.Require().Error(err)
}

func (s *CrashesModuleTestSuite) TestAnalyseUpdatedRegex() {
	_, err := s.mod.Update(nil, defaultCrash.Name, domain.CrashUpdate{Regexes: otherCrash.Regexes})
	s.Require().NoError(err)

	matches, err := s.mod.Analyse(nil, otherCrash.Regexes[0])
	s.Require().NoError(err)
	s.Require().Len(matches, 1)
}
