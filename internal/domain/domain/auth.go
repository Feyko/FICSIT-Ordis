package domain

type Role struct {
	Name        string
	Permissions []Permission
}

type Permission string

const (
	PermissionContentEditing   Permission = "ContentEditing"
	PermissionTokenCreation    Permission = "TokenCreation"
	PermissionTicketManagement Permission = "TicketManagement"
)

var (
	RoleAdmin = Role{
		Name:        "Admin",
		Permissions: []Permission{PermissionContentEditing, PermissionTokenCreation, PermissionTicketManagement},
	}
	RoleModerator = Role{
		Name:        "Moderator",
		Permissions: []Permission{PermissionContentEditing, PermissionTicketManagement},
	}
)
