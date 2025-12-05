package database

import (
	"gorm.io/gorm"
)

// Модель товара
type Product struct {
	// Встраиваем gorm.Model — добавляет ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	SupplierID  int        `json:"supplier_id"`                                     // Связь: один товар принадлежит одному поставщику
	Supplier    Supplier   `json:"-"`                                               // "-" исключает поле из JSON
	Categories  []Category `json:"categories" gorm:"many2many:product_categories;"` // many2many — промежуточная таблица
}

// Модель поставщика
type Supplier struct {
	// Встраиваем gorm.Model — добавляет ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model
	Name    string `json:"name" gorm:"unique"`
	Contact string `json:"contact"`
}

// Модель категории
type Category struct {
	// Встраиваем gorm.Model — добавляет ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model
	Name     string    `json:"name" gorm:"unique"`
	Products []Product `gorm:"many2many:product_categories;"` // many2many — промежуточная таблица
}

// Промежуточная таблица для связи многие-ко-многим
type ProductCategory struct {
	ProductID  int
	Product    Product
	CategoryID int
	Category   Category
}

// Модель пользователя
type User struct {
	// Встраиваем gorm.Model — добавляет ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}
