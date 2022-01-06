package auth

type Role string

const (
	Normal Role = "Normal"
	Admin  Role = "Admin"
	Super  Role = "Super"
)

type User struct {
	IdNum       int           `json:"idNum" db:"idNum"`
	UserId      string        `json:"userId" db:"userId"`
	EMail       string        `json:"email" db:"email"`
	FirstName   string        `json:"firstName" db:"firstName"`
	LastName    string        `json:"lastName" db:"lastName"`
	Role        string        `json:"role" db:"role"`
	GroupsIDs   []string      `json:"groups"`
	Permissions PermissionSet `json:"permissions"`
}

type Group struct {
	IdNum   int    `json:"idNum" db:"idNum"`
	GroupID string `json:"groupID" db:"groupID"`
	Name    string `json:"name" db:"name"`
}
