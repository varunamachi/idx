package core

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/varunamachi/libx"
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
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy int64     `db:"created_by" json:"createdBy"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy int64     `db:"updated_by" json:"updatedBy"`
}

type Token struct {
	Token      string `db:"token" json:"token"`
	UniqueName string `db:"uniqueName" json:"unique_name"`
	AssocType  string `db:"assoc_type" json:"assocType"`
	Operation  string `db:"operation" json:"operation"`
	Created    string `db:"created" json:"created"`
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

type Hasher interface {
	Hash(pw string) (string, error)
	Verify(pw, hash string) (bool, error)
}

type SecretStorage interface {
	SetPassword(gtx context.Context, creds *Creds) error
	UpdatePassword(gtx context.Context, creds *Creds, newPw string) error
	Verify(gtx context.Context, creds *Creds) error

	StoreToken(gtx context.Context, token *Token) error
	VerifyToken(gtx context.Context, id, operation, token string) error
}

func NewToken(uname, operation, assocType string) *Token {
	return &Token{
		UniqueName: uname,
		Operation:  operation,
		AssocType:  assocType,
		Token:      uuid.NewString(),
	}
}
