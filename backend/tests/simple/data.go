package simple

import (
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
)

var superUser = &core.User{
	DbItem:    core.DbItem{},
	UserId:    "super",
	EmailId:   "super@example.com",
	AuthzRole: auth.Super,
	State:     core.Active,
	FirstName: "Super",
	LastName:  "User",
	Title:     "Dr",
	Props: map[string]any{
		"initialUser": true,
	},
}

var users = []*core.User{
	{
		DbItem: core.DbItem{
			Id:        0,
			CreatedAt: time.Time{},
			CreatedBy: 0,
			UpdatedAt: time.Time{},
			UpdatedBy: 0,
		},
		UserId:    "admin_1",
		EmailId:   "admin1@example.com",
		AuthzRole: auth.Admin,
		State:     core.None,
		FirstName: "Admin",
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
		UserId:    "",
		EmailId:   "user2@example.com",
		AuthzRole: auth.Admin,
		State:     core.None,
		FirstName: "User",
		LastName:  "One",
		Title:     "Ms",
		Props: map[string]any{
			"test": "test",
		},
	},
}

var services = []*core.Service{
	{
		DbItem: core.DbItem{
			Id:        0,
			CreatedAt: time.Time{},
			CreatedBy: 0,
			UpdatedAt: time.Time{},
			UpdatedBy: 0,
		},
		Name:        "svc1",
		OwnerId:     0,
		DisplayName: "service1",
		Permissions: auth.PermissionTree{
			Permissions: []*auth.PermissionNode{
				{
					Id:         0,
					PermId:     "svc1_perm_0_0",
					Name:       "Perm 0.0",
					Predefined: false,
					Children: []*auth.PermissionNode{
						{
							Id:         0,
							PermId:     "svc1_perm_0_1",
							Name:       "Perm 0.1",
							Predefined: false,
							Children: []*auth.PermissionNode{
								{
									Id:         0,
									PermId:     "svc1_perm_0_1_1",
									Name:       "Perm 0.1.1",
									Predefined: false,
								},
							},
						},
					},
				},
				{
					Id:         0,
					PermId:     "svc1_perm_1_0",
					Name:       "Perm 1.0",
					Predefined: false,
					Children: []*auth.PermissionNode{
						{
							Id:         0,
							PermId:     "svc1_perm_1_1",
							Name:       "Perm 1.1",
							Predefined: false,
							Children:   nil,
						},
					},
				},
			},
		},
	},
}

var groups = []*core.Group{
	{},
}
