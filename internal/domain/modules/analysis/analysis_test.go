package analysis

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/domain/modules/crashes"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAnalysisModuleTestSuite(t *testing.T) {
	suite.Run(t, new(AnalysisModuleTestSuite))
}

type AnalysisModuleTestSuite struct {
	suite.Suite
	rep        repo.Repository[id.IDer]
	mod        *Module
	crashesMod *crashes.Module
}

func (s *AnalysisModuleTestSuite) SetupSuite() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)
	s.rep = rep

	authModule, err := auth.New(auth.Config{Secret: "test-secret"})
	s.Require().NoError(err)

	crashesMod, err := crashes.New(
		crashes.Config{
			auth.AuthedConfig{
				true,
				authModule,
			},
		}, s.rep)
	s.Require().NoError(err)

	s.crashesMod = crashesMod

	err = s.crashesMod.Create(nil, defaultCrash)
	s.Require().NoError(err)

	mod, err := New(Config{
		CrashesModule: s.crashesMod,
	})
	s.Require().NoError(err)
	s.mod = mod
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

func (s *AnalysisModuleTestSuite) TestEmptyText() {
	_, err := s.mod.AnalyseText(nil, nil)
	s.Require().NoError(err)
}

func (s *AnalysisModuleTestSuite) TestOneModInList() {
	result, err := s.mod.AnalyseText(nil, []byte("[2022.04.10-16.15.17:090][  0]LogPakFile: New pak file ../../../FactoryGame/Mods/SomeMod/Content/Paks/WindowsNoEditor/SomeModFactoryGame-WindowsNoEditor.pak added to pak precacher."))
	s.Require().NoError(err)
	s.Require().Len(result.ModList, 1)
	s.Equal("SomeMod", result.ModList[0])
}

func (s *AnalysisModuleTestSuite) TestSMLVersion() {
	result, err := s.mod.AnalyseText(nil, []byte("[2022.04.10-16.15.12:916][  0]LogSatisfactoryModLoader: Display: Satisfactory Mod Loader v.3.3.0+4b0a5be3 pre-initializing..."))
	s.Require().NoError(err)
	s.Require().NotNil(result.SMLVersion)
	s.Equal("3.3.0", *result.SMLVersion)
}

func (s *AnalysisModuleTestSuite) TestGameVersion() {
	result, err := s.mod.AnalyseText(nil, []byte("LogInit: Net CL: 202470"))
	s.Require().NoError(err)
	s.Require().NotNil(result.GameVersion)
	s.Equal("202470", *result.GameVersion)
}

func (s *AnalysisModuleTestSuite) TestPath() {
	result, err := s.mod.AnalyseText(nil, []byte("LogInit: Base Directory: Z:/home/feyko/Games/Heroic/SatisfactoryExperimental/Engine/Binaries/Win64/\nLogInit: Allocator: binned2\nLogInit: Installed Engine Build: 1"))
	s.Require().NoError(err)
	s.Require().NotNil(result.Path)
	s.Equal("Z:/home/feyko/Games/Heroic/SatisfactoryExperimental/Engine/Binaries/Win64/", *result.Path)
}

func (s *AnalysisModuleTestSuite) TestCommandLine() {
	result, err := s.mod.AnalyseText(nil, []byte("LogInit: Command Line:  -AUTH_LOGIN=unused -AUTH_PASSWORD=44ee2cbf3a0a4266aa052cc296a50979 -AUTH_TYPE=exchangecode -epicapp=CrabTest -epicenv=Prod -EpicPortal -epicusername=Feykoo -epicuserid=400dd576dacd47bb86815721d9dc3b28 -epiclocale=en -epicsandboxid=crab"))
	s.Require().NoError(err)
	s.Require().NotNil(result.CommandLine)
	s.Equal("-AUTH_LOGIN=unused -AUTH_PASSWORD=44ee2cbf3a0a4266aa052cc296a50979 -AUTH_TYPE=exchangecode -epicapp=CrabTest -epicenv=Prod -EpicPortal -epicusername=Feykoo -epicuserid=400dd576dacd47bb86815721d9dc3b28 -epiclocale=en -epicsandboxid=crab", *result.CommandLine)
}

func (s *AnalysisModuleTestSuite) TestLauncherID() {
	result, err := s.mod.AnalyseText(nil, []byte("LogInit: Launcher ID: epic"))
	s.Require().NoError(err)
	s.Require().NotNil(result.LauncherID)
	s.Equal("epic", *result.LauncherID)
}

func (s *AnalysisModuleTestSuite) TestLauncherArtifact() {
	result, err := s.mod.AnalyseText(nil, []byte("LogInit: Launcher Artifact: CrabTest"))
	s.Require().NoError(err)
	s.Require().NotNil(result.LauncherArtifact)
	s.Equal("CrabTest", *result.LauncherArtifact)
}

func (s *AnalysisModuleTestSuite) TestDesiredSMLVersion() {
	result, err := s.mod.AnalyseText(nil, []byte("LogInit: Net CL: 201345\nSatisfactory Mod Loader v.3.3.0"))
	s.Require().NoError(err)
	s.Require().NotNil(result.DesiredSMLVersion)
	s.Require().Equal("3.3.2", *result.DesiredSMLVersion)
}

func (s *AnalysisModuleTestSuite) TestNoDesiredSMLVersionWhenNoSMLVersion() {
	result, err := s.mod.AnalyseText(nil, []byte("LogInit: Net CL: 201345"))
	s.Require().NoError(err)
	s.Require().Nil(result.DesiredSMLVersion)
}
