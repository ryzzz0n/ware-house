package core

import (
	"warehouse-app/database"
)

// Интерфейс репозитория продуктов (вторичный порт)
type ProductRepository interface {
	// Создаёт поставщика
	SaveSupplier(supplier *database.Supplier) error
	// Создаёт категорию
	SaveCategory(category *database.Category) error
	// Возвращает всех поставщиков
	FindAllSupplier() ([]database.Supplier, error)
	// Возвращает все категории
	FindAllCategory() ([]database.Category, error)
	// Создаёт продукт
	SaveProduct(product *database.Product) error
	// Находит продукт по ID
	FindProductByID(productID uint) (*database.Product, error)
	// Возвращает все продукты
	FindAllProduct() ([]database.Product, error)
	// Возвращает все продукты по ID категории
	FindAllProductOfCategory(categoryID uint) ([]database.Product, error)
	// Возвращает все продукты по ID поставщика
	FindAllProductOfSupplier(supplierID uint) ([]database.Product, error)
	// Обновляет поставщика
	UpdateSupplier(updatedSupplier *database.Supplier) error
	// Обновляет продукт по ID
	UpdateProductByID(productID uint, updatedProduct *database.Product) error
	// Удаляет продукт по ID
	DeleteProductByID(productID uint) error
	// Удаляет поставщика по ID
	DeleteSupplierByID(supplierID uint) error
	// Удаляет категорию по ID
	DeleteCategoryByID(categoryID uint) error
}
