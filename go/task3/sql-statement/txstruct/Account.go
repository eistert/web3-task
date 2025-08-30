package txstruct

import "time"

/************ 数据模型（金额用分） ************/
type Account struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Balance   int64     `gorm:"not null;default:0"       json:"balance"` // 分
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
