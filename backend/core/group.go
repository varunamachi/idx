package core

type Group struct {
	DbItem
	Name        string   `db:"name" json:"name"`
	DisplayName string   `db:"display_name" json:"displayName"`
	Description string   `db:"description" json:"description"`
	Perms       []string `json:"perms"`
}

type GroupStorage struct {
}
