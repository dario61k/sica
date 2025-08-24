package models

type Auth struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Password string `json:"password" gorm:"type:text;not null"`
}
