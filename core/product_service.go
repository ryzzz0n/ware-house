package core

import (
	"warehouse-app/database"
)

// Интерфейс сервиса продуктов (первичный порт)
type ProductService interface {
	CreateSupplier(supplier *database.Supplier) error
	CreateCategory(category *database.Category) error
	GetAllSupplier() ([]database.Supplier, error)
	GetAllCategory() ([]database.Category, error)
	CreateProduct(product *database.Product) error
	GetProductByID(productID uint) (*database.Product, error)
	GetAllProduct() ([]database.Product, error)
	GetAllProductOfCategory(categoryID uint) ([]database.Product, error)
	GetAllProductOfSupplier(supplierID uint) ([]database.Product, error)
	UpdateSupplier(supplier *database.Supplier) error
	UpdateProductByID(productID uint, updatedProduct *database.Product) error
	DeleteProductByID(productID uint) error
	DeleteSupplierByID(supplierID uint) error
	DeleteCategoryByID(categoryID uint) error

}

// Реализация сервиса продуктов
type productServiceImpl struct {
	// Поле repo — это репозиторий продуктов
	repo ProductRepository
}

// Создаёт новый экземпляр ProductService
func NewProductService(repo ProductRepository) ProductService {
	return &productServiceImpl{repo: repo}
}

// Вызывает репозиторий для создания поставщика
func (s *productServiceImpl) CreateSupplier(supplier *database.Supplier) error {
	if err := s.repo.SaveSupplier(supplier); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для создания категории
func (s *productServiceImpl) CreateCategory(category *database.Category) error {
	if err := s.repo.SaveCategory(category); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для получения всех поставщиков
func (s *productServiceImpl) GetAllSupplier() ([]database.Supplier, error) {
	suppliers, err := s.repo.FindAllSupplier()
	if err != nil {
		return nil, err
	}
	return suppliers, err
}

// Вызывает репозиторий для получения всех категорий
func (s *productServiceImpl) GetAllCategory() ([]database.Category, error) {
	categories, err := s.repo.FindAllCategory()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Вызывает репозиторий для создания продукта
func (s *productServiceImpl) CreateProduct(product *database.Product) error {
	if err := s.repo.SaveProduct(product); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для получения продукта по ID
func (s *productServiceImpl) GetProductByID(productID uint) (*database.Product, error) {
	product, err := s.repo.FindProductByID(productID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// Вызывает репозиторий для получения всех продуктов по ID категории
func (s *productServiceImpl) GetAllProductOfCategory(categoryID uint) ([]database.Product, error) {
	products, err := s.repo.FindAllProductOfCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Вызывает репозиторий для получения всех продуктов по ID поставщика
func (s *productServiceImpl) GetAllProductOfSupplier(supplierID uint) ([]database.Product, error) {
	products, err := s.repo.FindAllProductOfSupplier(supplierID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Вызывает репозиторий для удаления поставщика по ID поставщика
func (s *productServiceImpl) DeleteSupplierByID(supplierID uint) error {
	if err := s.repo.DeleteSupplierByID(supplierID); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для удаления категории по ID категории
func (s *productServiceImpl) DeleteCategoryByID(categoryID uint) error {
	if err := s.repo.DeleteCategoryByID(categoryID); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для обновления поставщика
func (s *productServiceImpl) UpdateSupplier(supplier *database.Supplier) error {
	if err := s.repo.UpdateSupplier(supplier); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для получения всех продуктов
func (s *productServiceImpl) GetAllProduct() ([]database.Product, error) {
	products, err := s.repo.FindAllProduct()
	if err != nil {
		return nil, err
	}
	return products, err
}

// Вызывает репозиторий для обновления продукта по ID
func (s *productServiceImpl) UpdateProductByID(productID uint, updatedProduct *database.Product) error {
	if err := s.repo.UpdateProductByID(productID, updatedProduct); err != nil {
		return err
	}
	return nil
}

// Вызывает репозиторий для удаления продукта по ID
func (s *productServiceImpl) DeleteProductByID(productID uint) error {
	if err := s.repo.DeleteProductByID(productID); err != nil {
		return err
	}
	return nil
}