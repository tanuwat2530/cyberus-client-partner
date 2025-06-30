package models

type TmvhSubscriptionLog struct {
	ID            string `gorm:"primaryKey;column:id"`
	Action        string `gorm:"column:action"`
	Code          string `gorm:"column:code"`
	CyberusReturn string `gorm:"column:cyberus_return"`
	Description   string `gorm:"column:description"`
	Media         string `gorm:"column:media"`
	Msisdn        string `gorm:"column:msisdn"`
	Operator      string `gorm:"column:operator"`
	RefID         string `gorm:"column:ref_id"`
	ShortCode     string `gorm:"column:short_code"`
	Timestamp     int    `gorm:"column:timestamp"` // Can use int64 if needed
	Token         string `gorm:"column:token"`
	TranRef       string `gorm:"column:tran_ref"`
}
