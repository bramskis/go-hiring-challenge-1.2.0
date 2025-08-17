package models

// Category represents a category to which products could belong.
// It includes a unique id, unique code, and a name.
type Category struct {
	ID       uint      `gorm:"primaryKey"`
	Code     string    `gorm:"uniqueIndex;not null"`
	Name     string    `gorm:"not null"`
	Products []Product `gorm:"foreignKey:ProductID"`
}

func (p *Category) TableName() string {
	return "categories"
}
