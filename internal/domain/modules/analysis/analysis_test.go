package analysis

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAnalysisModuleTestSuite(t *testing.T) {
	suite.Run(t, new(AnalysisModuleTestSuite))
}

type AnalysisModuleTestSuite struct {
	suite.Suite
	rep repo.Repository[domain.Command]
	mod *Module
}

func (s *AnalysisModuleTestSuite) SetupSuite() {
	mod, err := New(Config{})
	s.Require().NoError(err)
	s.mod = mod
}
