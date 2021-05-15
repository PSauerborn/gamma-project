package roles

import "errors"

var ErrInvalidRole = errors.New("cannot convert to role: invalid role")

func StringToRole(role string) (Role, error) {
	var r Role
	switch role {
	case "Standard User":
		r = StandardUser
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
