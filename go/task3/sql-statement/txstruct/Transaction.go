package txstruct

import "time"

type Transaction struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FromAccountID uint64    `gorm:"index;not null"            json:"from_account_id"`
	ToAccountID   uint64    `gorm:"index;not null"            json:"to_account_id"`
	Amount        int64     `gorm:"not null"                  json:"amount"` // 分
	CreatedAt     time.Time `json:"created_at"`
}
