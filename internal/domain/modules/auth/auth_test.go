package auth

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/domain"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
	"testing"
)

func TestAuthModuleTestSuite(t *testing.T) {
	suite.Run(t, new(AuthModuleTestSuite))
}

type AuthModuleTestSuite struct {
	suite.Suite
	mod *Module
}

func (s *AuthModuleTestSuite) SetupTest() {
	mod, err := New(config.AuthConfig{Secret: "test-secret"})
	s.Require().NoError(err)
	s.mod = mod
}

func (s *AuthModuleTestSuite) TestNewToken() {
	_, err := s.mod.NewToken()
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestNewTokenIsValid() {
	token, err := s.mod.NewToken()
	s.Require().NoError(err)
	err = s.mod.ValidateToken(&token)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestInvalidToken() {
	err := s.mod.ValidateToken(&Token{})
	s.Require().Error(err)
}

func (s *AuthModuleTestSuite) TestNewTokenWithRolesIsValid() {
	token, err := s.mod.NewToken(domain.RoleAdmin, domain.RoleModerator)
	s.Require().NoError(err)
	err = s.mod.ValidateToken(&token)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestNewTokenWithRolesHasRoles() {
	roles := []domain.Role{domain.RoleAdmin, domain.RoleModerator}
	token, err := s.mod.NewToken(roles...)
	s.Require().NoError(err)
	s.Equal(roles, token.Roles)
}

func (s *AuthModuleTestSuite) TestNewTokenWithRolesHasRolesValidated() {
	roles := []domain.Role{domain.RoleAdmin, domain.RoleModerator}
	token, err := s.mod.NewToken(roles...)
	s.Require().NoError(err)
	err = s.mod.ValidateToken(&token)
	s.Require().NoError(err)
	s.Equal(roles, token.Roles)
}

func (s *AuthModuleTestSuite) TestTokenCorrectPermissionsOneRole() {
	perms := domain.RoleAdmin.Permissions
	token, err := s.mod.NewToken(domain.RoleAdmin)
	s.Require().NoError(err)
	s.Equal(perms, token.Permissions)
}

func (s *AuthModuleTestSuite) TestTokenCorrectPermissionsMultipleRoles() {
	rolePermTickets := domain.Role{Permissions: []domain.Permission{domain.PermissionTicketManagement}}
	rolePermContent := domain.Role{Permissions: []domain.Permission{domain.PermissionContentEditing}}
	rolePermTokens := domain.Role{Permissions: []domain.Permission{domain.PermissionTokenCreation}}
	allPerms := []domain.Permission{domain.PermissionContentEditing, domain.PermissionTokenCreation, domain.PermissionTicketManagement}
	token, err := s.mod.NewToken(rolePermTickets, rolePermContent, rolePermTokens)
	s.Require().NoError(err)
	slices.Sort(token.Permissions)
	slices.Sort(allPerms)
	s.Equal(allPerms, token.Permissions)
}
