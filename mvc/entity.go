package mvc

import "time"

type BaseEntity struct {
	ID        int64     `gorm:"column:id; type:bigint; primaryKey; not null" json:"id"`
	Remark    string    `gorm:"column:remark; type:text; default: null" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at; type:datetime; default now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at; type:datetime; default now()" json:"updatedAt"`
}
