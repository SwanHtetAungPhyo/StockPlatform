package models


type User struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `gorm:"size:100;not null" json:"name"`
	Email    string  `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password string  `gorm:"size:255;not null" json:"password"`
	Balance  Balance `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"balance"`
}

type Balance struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Amount float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	UserID uint    `gorm:"not null;index" json:"user_id"`
}
