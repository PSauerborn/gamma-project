package roles

type Role int

func (r Role) String() string {
	return [...]string{"Standard User", "Clerk", "Planner", "Admin"}[r-1]
}

func (r Role) IsValid() bool {
	return r >= 1 && r <= 4
}

const (
	StandardUser Role = iota + 1
	Clerk
	Planner
	Admin
)

type Persistence interface {
	GetUserRole(uid string) (Role, error)
	SetUserRole(uid string, role Role) error
}
