package roles

type Role int

func (r Role) String() string {
	return [...]string{"StandardUser", "Clerk", "Planner", "Admin"}[r]
}

const (
	StandardUser Role = iota
	Clerk
	Planner
	Admin
)

type RolesPersistence interface {
	GetUserRoles(uid string) (Role, error)
	SetUserRoles(uid string, role Role) error
}
