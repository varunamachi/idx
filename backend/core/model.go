package core

import (
	"context"
	"time"

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
	Id        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy uint64    `db:"created_by" json:"createdBy"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy uint64    `db:"updated_by" json:"updatedBy"`
}

type Token struct {
	Token     string `db:"token" json:"token"`
	Id        string `db:"id" json:"id"`
	AssocType string `db:"assoc_type" json:"assocType"`
	Operation string `db:"operation" json:"operation"`
	Created   string `db:"created" json:"created"`
}

type Creds struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	Type     string `json:"type"`
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
	VerifyToken(gtx context.Context, token, id, operation string) error
}
