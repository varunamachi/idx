package core

import (
	"github.com/varunamachi/libx/data"
)

type Group struct {
	DbItem
	Name        string           `db:"name" json:"name"`
	DisplayName string           `db:"display_name" json:"displayName"`
	Description string           `db:"description" json:"description"`
	Perms       data.Vec[string] `db:"perms" json:"perms"`
}

func MergePerms(groups ...Group) []string {
	// for _, g := range groups {

	// }

	return []string{}
}
