package core

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

var (
	ErrPwPatternMismatch = errors.New(
		"password does not match required pattern")
)

var GitTag = "--"
var GitHash = "--"
var GitBranch = "--"
var BuildTime = "--"
var BuildHost = "--"
var BuildUser = "--"

var bi = libx.BuildInfo{
	GitTag:    GitTag,
	GitHash:   GitHash,
	GitBranch: GitBranch,
	BuildTime: BuildTime,
	BuildHost: BuildHost,
	BuildUser: BuildUser,
}

func GetBuildInfo() *libx.BuildInfo {
	return &bi
}

type DbItem struct {
	Id        int64     `db:"id" json:"id"`
	CreatedOn time.Time `db:"created_on" json:"createdOn"`
	CreatedBy int64     `db:"created_by" json:"createdBy"`
	UpdatedOn time.Time `db:"updated_on" json:"updatedOn"`
	UpdatedBy int64     `db:"updated_by" json:"updatedBy"`
}

type Token struct {
	Token      string `db:"token" json:"token"`
	UniqueName string `db:"unique_name" json:"uniqueName"`
	AssocType  string `db:"assoc_type" json:"assocType"`
	Operation  string `db:"operation" json:"operation"`
	CreatedOn  string `db:"createdOn" json:"created_on"`
}

type AuthEntity string

const (
	AuthUser    AuthEntity = "user"
	AuthService AuthEntity = "service"
)

type Creds struct {
	UniqueName string     `json:"uniqueName" db:"unique_name"`
	Password   string     `json:"password" db:"password"`
	Type       AuthEntity `json:"type" db:"type"`
}

type Secret struct {
	UniqueName    string           `json:"uniqueName" db:"unique_name"`
	PasswordHash  string           `json:"password_hash" db:"password_hash"`
	Type          AuthEntity       `json:"type" db:"type"`
	CreatedOn     time.Time        `json:"createdOn" db:"created_on"`
	NumFailedAuth int              `json:"numFailedAuth" db:"num_failed_auth"`
	LastFailedOn  time.Time        `json:"lastFailedOn" db:"last_failed_on"`
	PrevPasswords data.Vec[string] `json:"prevPasswords" db:"prev_passwords"`
}

type CredentialPolicy struct {
	ItemType       AuthEntity    `db:"itemType" json:"item_type"`
	Pattern        string        `db:"pattern" json:"pattern"`
	Expiry         time.Duration `db:"expiry" json:"expiry"`
	MaxRetries     int           `db:"max_retries" json:"maxRetries"`
	RetryResetDays int           `db:"retry_reset_days" json:"retryResetDays"`
	MaxReuse       int           `db:"max_reuse" json:"maxReuse"`
	pattern        *regexp.Regexp
}

func (cp *CredentialPolicy) MatchPattern(pw string) error {
	if cp.pattern == nil {
		var err error
		if cp.pattern, err = regexp.Compile(cp.Pattern); err != nil {
			return err
		}
	}

	if !cp.pattern.MatchString(pw) {
		return errx.Errf(ErrPwPatternMismatch,
			"password does not match pattern defined by the policy")
	}
	return nil
}

type Hasher interface {
	Hash(pw string) (string, error)
	Verify(pw, hash string) error
}

type SecretStorage interface {
	CreatePassword(gtx context.Context, creds *Creds) error
	UpdatePassword(gtx context.Context, creds *Creds) error
	Authenticate(gtx context.Context, creds *Creds) error

	StoreToken(gtx context.Context, token *Token) error
	VerifyToken(gtx context.Context, id, operation, token string) error

	CredentialPolicy(
		gtx context.Context, credType AuthEntity) (*CredentialPolicy, error)
	SetCredentialPolicy(
		gtx context.Context, cp *CredentialPolicy) error
}

func NewToken(uname, operation, assocType string) *Token {
	return &Token{
		UniqueName: uname,
		Operation:  operation,
		AssocType:  assocType,
		Token:      uuid.NewString(),
	}
}

type Group struct {
	DbItem
	ServiceId   int      `db:"service_id" json:"service_id"`
	Name        string   `db:"name" json:"name"`
	DisplayName string   `db:"display_name" json:"displayName"`
	Description string   `db:"description" json:"description"`
	Perms       []string `json:"perms"`
}

type GroupController interface {
	Save(gtx context.Context, group *Group) (int64, error)
	Update(gtx context.Context, group *Group) error
	GetOne(gtx context.Context, id int64) (*Group, error)
	Remove(gtx context.Context, id int64) error
	Get(gtx context.Context, params *data.CommonParams) ([]*Group, error)

	Exists(gtx context.Context, id int64) (bool, error)
	Count(gtx context.Context, filter *data.Filter) (int64, error)

	SetPermissions(gtx context.Context, groupId int64, perms []string) error
	GetPermissions(gtx context.Context, groupId int64) ([]string, error)

	AddToGroups(gtx context.Context, userId int64, groupIds ...int64) error
	RemoveFromGroup(gtx context.Context, userId, groupId int64) error

	// Storage() GroupStorage
	SaveWithPerms(
		gtx context.Context, group *Group, perms []string) (int64, error)
}

type Service struct {
	DbItem
	Name        string              `db:"name" json:"name"`
	OwnerId     int64               `db:"owner_id" json:"ownerId"`
	DisplayName string              `db:"display_name" json:"displayName"`
	Permissions auth.PermissionTree `db:"permissions" json:"permissions"`
}

type ServiceController interface {
	Save(gtx context.Context, service *Service) (int64, error)
	Update(gtx context.Context, service *Service) error
	GetOne(gtx context.Context, id int64) (*Service, error)
	GetByName(gtx context.Context, name string) (*Service, error)
	GetForOwner(gtx context.Context, ownerId string) ([]*Service, error)
	Remove(gtx context.Context, id int64) error
	Get(gtx context.Context, params *data.CommonParams) ([]*Service, error)

	AddAdmin(gtx context.Context, serviceId, userId int64) error
	GetAdmins(gtx context.Context, serviceId int64) ([]*User, error)
	RemoveAdmin(gtx context.Context, serviceId, userId int64) error
	IsAdmin(gtx context.Context, serviceId, userId int64) (bool, error)

	Exists(gtx context.Context, name string) (bool, error)
	Count(gtx context.Context, filter *data.Filter) (int64, error)

	GetPermissionForService(
		gtx context.Context, userId, serviceId int64) ([]string, error)
}
