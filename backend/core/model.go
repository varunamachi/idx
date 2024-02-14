package core

import (
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

type DbItem struct {
	Id        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy uint64    `db:"created_by" json:"createdBy"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy uint64    `db:"updated_by" json:"updatedBy"`
}

func GetBuildInfo() *libx.BuildInfo {
	return &bi
}

type CredentialStorage interface {
	SetPassword(itemType, id, password string) error
	UpdatePassword(itemType, id, oldPw, newPw string) error
	Verify(itemType, id, password string) error
}
