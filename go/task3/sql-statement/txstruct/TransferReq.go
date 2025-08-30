package txstruct

/************ 传参结构 ************/
type TransferReq struct {
	FromID uint64 `json:"from_id" binding:"required,gt=0"`
	ToID   uint64 `json:"to_id"   binding:"required,gt=0,nefield=FromID"`
	Amount int64  `json:"amount"  binding:"required,gt=0"` // 分
}
