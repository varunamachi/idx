package core

import "github.com/varunamachi/libx/auth"

type App struct {
	DbItem
	Name        string
	DisplayName string
	Permissions auth.PermissionTree
}
