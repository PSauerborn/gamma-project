package roles

type Role int

func (r Role) String() string {
	return [...]string{"Standard User", "Clerk", "Planner", "Admin"}[r]
}

func (r Role) IsValid() bool {
	return r > 0 && r < 4
}

const (
	StandardUser Role = iota
	Clerk
	Planner
	Admin
)

type Persistence interface {
	GetUserRole(uid string) (Role, error)
	SetUserRole(uid string, role Role) error
}
