package models

// Product 产品
type Product struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price uint64
	Image string
}

func (*Product) TableName() string {
	return "product"
}
