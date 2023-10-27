package entity

// type Role struct {
// 	ID          uint
// 	Title       string
// 	Description string
// }

type Role uint8

const (
	UserRole Role = iota + 1
	SuperAdminRole
)

func (r Role) String() string {
	switch r {
	case UserRole:
		return "user"
	case SuperAdminRole:
		return "admin"
	}
	return ""
}
