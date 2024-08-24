package simple

import (
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
)

type userAndPassword struct {
	user        *core.User
	password    string
	newPassword string
}

var super = userAndPassword{

	user: &core.User{
		DbItem:    core.DbItem{},
		UName:     "super",
		EmailId:   "super@example.com",
		AuthzRole: auth.Super,
		State:     core.Active,
		FirstName: "Super",
		LastName:  "User",
		Title:     "Dr",
		Props: map[string]any{
			"initialUser": true,
		},
	},
	password:    "onetwothree",
	newPassword: "threetwoone",
}

var users = []userAndPassword{
	{
		user: &core.User{
			DbItem: core.DbItem{
				Id:        0,
				CreatedAt: time.Time{},
				CreatedBy: 0,
				UpdatedAt: time.Time{},
				UpdatedBy: 0,
			},
			UName:     "admin_1",
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
		password:    "onetwothree",
		newPassword: "threetwoone",
	},
	{
		user: &core.User{
			DbItem: core.DbItem{
				Id:        0,
				CreatedAt: time.Time{},
				CreatedBy: 0,
				UpdatedAt: time.Time{},
				UpdatedBy: 0,
			},
			UName:     "normal_1",
			EmailId:   "normal_1@example.com",
			AuthzRole: auth.Admin,
			State:     core.None,
			FirstName: "Normal",
			LastName:  "One",
			Title:     "Ms",
			Props: map[string]any{
				"test": "test",
			},
		},
		password:    "onetwothree",
		newPassword: "threetwoone",
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
	{
		DbItem: core.DbItem{
			Id:        0,
			CreatedAt: time.Time{},
			CreatedBy: 0,
			UpdatedAt: time.Time{},
			UpdatedBy: 0,
		},
		ServiceId:   0,
		Name:        "svc1_group_1",
		DisplayName: "S1G1",
		Description: "Service group 1",
		Perms: []string{
			"svc1_perm_0_0",
			"svc1_perm_0_1_1",
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
		ServiceId:   0,
		Name:        "svc1_group_2",
		DisplayName: "S2G2",
		Description: "Service group 2",
		Perms: []string{
			"svc1_perm_1_0",
			"svc1_perm_1_1",
		},
	},
}
