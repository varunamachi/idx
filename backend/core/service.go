package core

import "github.com/varunamachi/libx/auth"

type Service struct {
	DbItem
	Name        string              `db:"name" json:"name"`
	DisplayName string              `db:"display_name" json:"displayName"`
	Permissions auth.PermissionTree `db:"permissions" json:"permissions"`
}
