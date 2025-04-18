package models

// Table name on database
type ClientPartner struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}
