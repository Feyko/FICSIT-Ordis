package domain

type Role struct {
	ID          int
	Name        string
	Permissions []Permission
}

type Permission string

const (
	PermissionContentEditing   Permission = "ContentEditing"
	PermissionTokenCreation    Permission = "TokenCreation"
	PermissionTicketManagement Permission = "TicketManagement"
)

var Roles = map[int]Role{
	0: {
		ID:          0,
		Name:        "Default",
		Permissions: []Permission{},
	},
	1: {
		ID:          1,
		Name:        "Admin",
		Permissions: []Permission{PermissionContentEditing, PermissionTokenCreation, PermissionTicketManagement},
	},
	2: {
		ID:          2,
		Name:        "Moderator",
		Permissions: []Permission{PermissionContentEditing, PermissionTicketManagement},
	},
	3: {
		ID:          3,
		Name:        "Mediator",
		Permissions: []Permission{PermissionTicketManagement},
	},
}

var (
	RoleAdmin     = Roles[1]
	RoleModerator = Roles[2]
	RoleMediator  = Roles[3]
)
