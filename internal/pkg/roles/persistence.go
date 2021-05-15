package roles

import "errors"

type Role int

func (r Role) String() string {
	return [...]string{"Standard", "Clerk", "Planner", "Admin"}[r-1]
}

func (r Role) IsValid() bool {
	return r >= 1 && r <= 4
}

var ErrInvalidRole = errors.New("cannot convert to role: invalid role")

func StringToRole(role string) (Role, error) {
	var r Role
	switch role {
	case "Standard":
		r = Standard
	case "Clerk":
		r = Clerk
	case "Planner":
		r = Planner
	case "Admin":
		r = Admin
	default:
		return r, ErrInvalidRole
	}
	return r, nil
}

const (
	Standard Role = iota + 1
	Clerk
	Planner
	Admin
)

type Persistence interface {
	GetUserRole(uid string) (Role, error)
	SetUserRole(uid string, role Role) error
}
