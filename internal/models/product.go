package models

type Product struct {
	ID          uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string   `json:"name" gorm:"type:text;not null"`
	Description string   `json:"description" gorm:"type:text;not null"`
	Price       float64  `json:"price" gorm:"type:float;not null"`
	CategoryID  uint     `json:"category"`
	Category    Category `json:"category_detail" gorm:"foreignKey:CategoryID;references:ID"`

	Available bool   `json:"available" gorm:"type:bool"`
	Visible   bool   `json:"visible" gorm:"type:bool"`
	Image     string `json:"image" gorm:"type:text;not null"`
}
