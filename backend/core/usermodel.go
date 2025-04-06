package core

import (
	"context"

	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
)

// UserState - state of the user account
type UserState string

// Created - user account is created but not verified
var Created UserState = "created"

// Verfied - user account is verified by the user
var Verfied UserState = "verified"

// Active - user is active
var Active UserState = "active"

// Disabled - user account is disabled by an admin
var Disabled UserState = "disabled"

// Flagged - user account is flagged by a user
var Flagged UserState = "flagged"

// Invalid state
var None UserState = ""

var ValidUserStates = []UserState{
	Created,
	Verfied,
	Active,
	Disabled,
	Flagged,
}

type User struct {
	DbItem
	UName     string    `json:"userName" db:"user_name"`
	EmailId   string    `json:"email" db:"email"`
	AuthzRole auth.Role `json:"auth" db:"auth"`
	State     UserState `json:"state" bson:"state"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Title     string    `json:"title" db:"title"`
	Props     data.M    `json:"props,omitempty" db:"props"`
	// Perms     auth.PermissionSet `json:"perms,omitempty" db:"perms"`
}

func (u *User) Id() int64 {
	return u.DbItem.Id
}

func (u *User) Username() string {
	return u.UName
}

func (u *User) Email() string {
	return u.EmailId
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) Role() auth.Role {
	return u.AuthzRole
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
	if u.Props == nil {
		return nil
	}
	return u.Props[key]
}

func (u *User) SetProp(key string, value any) {
	if u.Props == nil {
		u.Props = data.M{}
	}
	u.Props[key] = value
}

type UserWithPassword struct {
	User     *User
	Password string `json:"password"`
}

type UserController interface {
	Save(gtx context.Context, user *User) (int64, error)
	Update(gtx context.Context, user *User) error
	GetOne(gtx context.Context, id int64) (*User, error)
	SetState(gtx context.Context, id int64, state UserState) error
	Remove(gtx context.Context, id int64) error
	Get(gtx context.Context, params *data.CommonParams) ([]*User, error)
	Count(gtx context.Context, filter *data.Filter) (int64, error)

	ByUsername(gtx context.Context, username string) (*User, error)
	Exists(gtx context.Context, username string) (bool, error)
	GetId(gtx context.Context, username string) (int64, error)

	CredentialStorage() SecretStorage

	Register(gtx context.Context, user *User, password string) (int64, error)
	Verify(gtx context.Context, userName, verToken string) error
	Approve(gtx context.Context,
		userId int64,
		role auth.Role,
		groups ...int64) error
	InitResetPassword(gtx context.Context, userName string) error
	ResetPassword(
		gtx context.Context, userName, token, newPassword string) error
	UpdatePassword(
		gtx context.Context,
		userName,
		oldPassword,
		newPassword string) error
}
