package core

import (
	"time"

	"github.com/varunamachi/libx/data"
)

type Group struct {
	Name        string           `db:"name" json:"name"`
	DisplayName string           `db:"display_name" json:"displayName"`
	Description string           `db:"description" json:"description"`
	Perms       data.Vec[string] `db:"perms" json:"perms"`
	CreatedOn   time.Time        `db:"created_on" json:"createdOn"`
	UpdatedOn   time.Time        `db:"updated_on" json:"updatedOn"`
}

func MergePerms(groups ...Group) []string {
	return []string{}
}
