package auth

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/domain"
	"context"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
	"testing"
)

const adminToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJSb2xlSURzIjpbMV19.--kQQn4mF12neGr69Z47ZVol7BTeqSWVs3vujbXPAjDwAiMe-mbR1LqUyt4dXgrVAAcu_gzhaEXtaY1-ksv0uA"

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
	allPerms := []domain.Permission{domain.PermissionContentEditing, domain.PermissionTicketManagement}
	domain.Roles[len(domain.Roles)] = domain.Role{ID: len(domain.Roles), Permissions: []domain.Permission{domain.PermissionContentEditing}}
	domain.Roles[len(domain.Roles)] = domain.Role{ID: len(domain.Roles), Permissions: []domain.Permission{domain.PermissionTicketManagement}}
	defer func() {
		delete(domain.Roles, len(domain.Roles)-1)
		delete(domain.Roles, len(domain.Roles)-1)
	}()

	token, err := s.mod.NewToken(domain.Roles[len(domain.Roles)-1], domain.Roles[len(domain.Roles)-2])
	s.Require().NoError(err)
	slices.Sort(token.Permissions)
	slices.Sort(allPerms)
	s.Equal(allPerms, token.Permissions)
}

func (s *AuthModuleTestSuite) TestTokenHasPermission() {
	domain.Roles[len(domain.Roles)] = domain.Role{ID: len(domain.Roles), Permissions: []domain.Permission{domain.PermissionTicketManagement}}
	defer func() {
		delete(domain.Roles, len(domain.Roles)-1)
	}()
	token, err := s.mod.NewToken(domain.Roles[len(domain.Roles)-1])
	s.Require().NoError(err)
	s.Require().True(token.HasPermissions(domain.PermissionTicketManagement))
	s.Require().False(token.HasPermissions(domain.PermissionContentEditing))
}

func (s *AuthModuleTestSuite) TestAuthorize() {
	ctx := context.WithValue(context.Background(), "Authorization", adminToken)
	err := s.mod.Authorize(&ctx)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestAuthorizePermissions() {
	ctx := context.WithValue(context.Background(), "Authorization", adminToken)
	err := s.mod.Authorize(&ctx, domain.PermissionTokenCreation)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestAuthorizePermissionsBearer() {
	ctx := context.WithValue(context.Background(), "Authorization", "Bearer "+adminToken)
	err := s.mod.Authorize(&ctx, domain.PermissionTokenCreation)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestAuthorizeWrongPermissions() {
	ctx := context.WithValue(context.Background(), "Authorization", adminToken)
	err := s.mod.Authorize(&ctx, "NonexistentPermission")
	s.Require().Error(err)
}

func (s *AuthModuleTestSuite) TestReauthorize() {
	ctx := context.WithValue(context.Background(), "Authorization", adminToken)
	err := s.mod.Authorize(&ctx, domain.PermissionTokenCreation)
	s.Require().NoError(err)
	err = s.mod.Authorize(&ctx, domain.PermissionTokenCreation)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestReauthorizeWrongPermissions() {
	ctx := context.WithValue(context.Background(), "Authorization", adminToken)
	err := s.mod.Authorize(&ctx, "NonexistentPermission")
	s.Require().Error(err)
	err = s.mod.Authorize(&ctx, "NonexistentPermission")
	s.Require().Error(err)
}
