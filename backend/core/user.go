package core

import (
	"context"

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
	DbItem
	UserId    string    `json:"userId" db:"user_id"`
	EmailId   string    `json:"email" db:"email"`
	Auth      auth.Role `json:"auth" db:"auth"`
	State     UserState `json:"state" bson:"state"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Title     string    `json:"title" db:"title"`
	Props     data.M    `json:"props,omitempty" db:"props"`
	// Perms     auth.PermissionSet `json:"perms,omitempty" db:"perms"`
}

func (u *User) SeqId() int {
	return int(u.DbItem.Id)
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
	Save(gtx context.Context, user *User) error
	Update(gtx context.Context, user *User) error
	GetOne(gtx context.Context, id int) (*User, error)
	GetByUserId(gtx context.Context, id string) (*User, error)
	SetState(gtx context.Context, id int, state UserState) error
	Remove(gtx context.Context, id int) error
	Get(gtx context.Context, params *data.CommonParams) ([]*User, error)
	AddToGroup(gtx context.Context, userId, groupId int) error
	RemoveFromGroup(gtx context.Context, userId, groupId int) error
	GetPermissionForService(
		gtx context.Context, userId, serviceId int) ([]string, error)

	SetPassword(userId, password string) error
	UpdatePassword(userId, oldPw, newPw string) error
	Verify(userId, password string) error
}

type Authenticator interface {
	Authenticate(gtx context.Context, creds map[string]any) error
}

type Authorizor interface {
	Authorize(gtx context.Context, serviceId, userId int) (*User, error)
}
