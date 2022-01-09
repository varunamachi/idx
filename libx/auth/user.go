package auth

type Role string

const (
	None   Role = "None"
	Normal Role = "Normal"
	Admin  Role = "Admin"
	Super  Role = "Super"
)

func (r Role) EqualOrAbove(another Role) bool {
	// Following logic only checks above condition
	switch r {
	case None:
		return true
	case Normal:
		return another == None
	case Admin:
		return another == None || another == Normal
	case Super:
		return false
	}

	// Following checks equal
	return r == another
}

type User struct {
	SeqId       int           `json:"seqId" db:"seqId"`
	Id          string        `json:"id" db:"id"`
	EMail       string        `json:"email" db:"email"`
	FirstName   string        `json:"firstName" db:"firstName"`
	LastName    string        `json:"lastName" db:"lastName"`
	Role        Role          `json:"role" db:"role"`
	GroupsIDs   []string      `json:"groups"`
	Permissions PermissionSet `json:"permissions"`
}

func (u *User) HasRole(role Role) bool {
	return role.EqualOrAbove(u.Role)
}

func (u *User) HasPerms(permIds ...string) bool {
	for _, perm := range permIds {
		if !u.Permissions.HasPerm(perm) {
			return false
		}
	}
	return true
}

type Group struct {
	SeqId int    `json:"seqId" db:"seqId"`
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
}

func ToRole(roleStr string) Role {
	switch roleStr {
	case "None":
		return None
	case "Normal":
		return Normal
	case "Admin":
		return Admin
	case "Super":
		return Super
	}
	return None
}

type UserRetrieverFunc func(userId string) (*User, error)
