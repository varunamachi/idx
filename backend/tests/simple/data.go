package simple

import (
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
)

var users = []*core.User{
	{
		DbItem: core.DbItem{
			Id:        0,
			CreatedAt: time.Time{},
			CreatedBy: 0,
			UpdatedAt: time.Time{},
			UpdatedBy: 0,
		},
		UserId:    "user1",
		EmailId:   "user1@example.com",
		AuthzRole: auth.Super,
		State:     core.None,
		FirstName: "User",
		LastName:  "One",
		Title:     "Mr",
		Props: map[string]any{
			"test": "test",
		},
	},
	{
		DbItem: core.DbItem{
			Id:        0,
			CreatedAt: time.Time{},
			CreatedBy: 0,
			UpdatedAt: time.Time{},
			UpdatedBy: 0,
		},
		UserId:    "user2",
		EmailId:   "user2@example.com",
		AuthzRole: auth.Super,
		State:     core.None,
		FirstName: "User",
		LastName:  "One",
		Title:     "Ms",
		Props: map[string]any{
			"test": "test",
		},
	},
}

var services = []*core.Service{}

var groups = []*core.Group{}
