package core

import "time"

type DbItem struct {
	Id        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy uint64    `db:"created_by" json:"createdBy"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy uint64    `db:"updated_by" json:"updatedBy"`
}
