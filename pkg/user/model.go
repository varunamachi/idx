package user

type User struct {
	IdNum       int                 `json:"idNum" db:"idNum"`
	UserId      string              `json:"userId" db:"userId"`
	EMail       string              `json:"email" db:"email"`
	FirstName   string              `json:"firstName" db:"firstName"`
	LastName    string              `json:"lastName" db:"lastName"`
	Role        string              `json:"role" db:"role"`
	GroupsIDs   []string            `json:"groups"`
	Permissions map[string]struct{} `json:"permissions"`
}

type Group struct {
	IdNum       int               `json:"idNum" db:"idNum"`
	GroupID     string            `json:"groupID" db:"groupID"`
	Name        string            `json:"name" db:"name"`
	Permissions []*PermissionNode `json:"permissions"`
}

type PermissionNode struct {
	IdNum      int               `json:"idNum" db:"idNum"`
	PermId     string            `json:"permId" db:"permId"`
	Name       string            `json:"name" db:"name"`
	Predefined bool              `json:"predefined" db:"predefined"`
	Base       string            `json:"base" db:"base"`
	Children   []*PermissionNode `json:"children"`
}
