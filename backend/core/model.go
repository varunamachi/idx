package core

import (
	"time"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
)

// UserState - state of the user account
type UserState string

// Verfied - user account is verified by the user
var Verfied UserState = "verified"

// Active - user is active
var Active UserState = "active"

// Disabled - user account is disabled by an admin
var Disabled UserState = "disabled"

// Flagged - user account is flagged by a user
var Flagged UserState = "flagged"

type User struct {
	ID         uint64    `json:"id" db:"id"`
	UserId     string    `json:"userId" db:"user_id"`
	EmailId    string    `json:"email" db:"email"`
	Auth       auth.Role `json:"auth" db:"auth"`
	FirstName  string    `json:"firstName" db:"first_name"`
	LastName   string    `json:"lastName" db:"last_name"`
	Title      string    `json:"title" db:"title"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	ModifiedAt time.Time `json:"modifiedAt" db:"modified_at"`
	Props      data.M    `json:"props,omitempty" db:"props"`
}

func (u *User) SeqId() int {
	return int(u.ID)
}

func (u *User) Id() string {
	return u.UserId
}

func (u *User) Email() string {
	return u.EmailId
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) Role() auth.Role {
	return u.Auth
}

func (u *User) GroupIds() []string {
	return []string{}
}

func (u *User) Permissions() auth.PermissionSet {
	return auth.PermissionSet{}
}

func (u *User) AddProp(key string, value any) {
	if u.Props == nil {
		u.Props = make(data.M)
	}
	u.Props[key] = value
}

func (u *User) Prop(key string) any {
	return u.Props[key]
}

type UserStorage interface {
	Save(user *User) (err error)
	Update(user *User) (err error)
	Get(id string) (user *User, err error)
	SetState(id string, state UserState) error
	Remove(id string) error
}
