package adapters

import (
	"warehouse-app/core"
	"warehouse-app/database"

	"gorm.io/gorm"
)

// GORM-реализация репозитория продуктов
type gormProductRepository struct {
	db *gorm.DB
}

// Создаёт новый экземпляр GORM-репозитория
func NewGormProductRepository(db *gorm.DB) core.ProductRepository {
	// Возвращаем указатель на новую структуру с подключением к БД
	return &gormProductRepository{db: db}
}

// Сохраняет поставщика в БД
func (r *gormProductRepository) SaveSupplier(supplier *database.Supplier) error {
	// Вызываем метод Create GORM, передаём поставщика
	if result := r.db.Create(supplier); result.Error != nil {
		// Если была ошибка — возвращаем её
		return result.Error
	}
	return nil
}

// Сохраняет категорию в БД
func (r *gormProductRepository) SaveCategory(category *database.Category) error {
	// Вызываем метод Create GORM, передаём категорию
	if result := r.db.Create(category); result.Error != nil {
		// Если была ошибка — возвращаем её
		return result.Error
	}
	return nil
}

// Возвращает всех поставщиков
func (r *gormProductRepository) FindAllSupplier() ([]database.Supplier, error) {
	// Создаём пустой срез поставщиков
	var suppliers []database.Supplier
	// Выполняем запрос Find — получаем всех поставщиков
	if result := r.db.Find(&suppliers); result.Error != nil {
		// Если ошибка — возвращаем nil и ошибку
		return nil, result.Error
	}
	// Иначе — возвращаем список поставщиков
	return suppliers, nil
}

// Возвращает все категории
func (r *gormProductRepository) FindAllCategory() ([]database.Category, error) {
	// Создаём пустой срез категорий
	var categories []database.Category
	// Выполняем запрос Find — получаем все категории
	if result := r.db.Find(&categories); result.Error != nil {
		// Если ошибка — возвращаем nil и ошибку
		return nil, result.Error
	}
	// Иначе — возвращаем список категорий
	return categories, nil
}

// Сохраняет продукт в БД
func (r *gormProductRepository) SaveProduct(product *database.Product) error {
	// Вызываем метод Create GORM, передаём продукт
	if result := r.db.Create(product); result.Error != nil {
		// Если была ошибка — возвращаем её
		return result.Error
	}
	return nil
}

// Находит продукт по ID с поставщиком и категориями
func (r *gormProductRepository) FindProductByID(productID uint) (*database.Product, error) {
	// Создаём переменную продукта
	var product database.Product
	// Выполняем запрос First с Preload — загружаем поставщика и категории
	if result := r.db.Preload("Supplier").Preload("Categories").First(&product, productID); result.Error != nil {
		// Если ошибка — возвращаем nil и ошибку
		return nil, result.Error
	}
	// Иначе — возвращаем указатель на продукт
	return &product, nil
}

// Находит все продукты по ID категории
func (r *gormProductRepository) FindAllProductOfCategory(categoryID uint) ([]database.Product, error) {
	// Создаём пустой срез продуктов
	var products []database.Product
	// Выполняем JOIN-запрос — ищем продукты по ID категории
	result := r.db.Preload("Supplier").Preload("Categories").Joins("JOIN product_categories on product_categories.product_id = products.id").
		Where("product_categories.category_id = ?", categoryID).
		Find(&products)
	// Если ошибка — возвращаем nil и ошибку
	if result.Error != nil {
		return nil, result.Error
	}
	// Иначе — возвращаем список продуктов
	return products, nil
}

// Находит все продукты по ID поставщика
func (r *gormProductRepository) FindAllProductOfSupplier(supplierID uint) ([]database.Product, error) {
	// Создаём пустой срез продуктов
	var products []database.Product
	// Выполняем JOIN-запрос — ищем продукты по ID поставщика
	result := r.db.Preload("Supplier").Preload("Categories").Joins("Join suppliers on suppliers.id = products.supplier_id").
		Where("suppliers.id = ?", supplierID).
		Find(&products)
	// Если ошибка — возвращаем nil и ошибку
	if result.Error != nil {
		return nil, result.Error
	}
	// Иначе — возвращаем список продуктов
	return products, nil
}

// Обновляет поставщика
func (r *gormProductRepository) UpdateSupplier(supplier *database.Supplier) error {
	// Вызываем метод Updates GORM — обновляем поля поставщика
	result := r.db.Model(&supplier).Updates(supplier)
	// Если ошибка — возвращаем её
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Возвращает все продукты
func (r *gormProductRepository) FindAllProduct() ([]database.Product, error) {
	// Создаём пустой срез продуктов
	var products []database.Product
	// Выполняем запрос Find — получаем все продукты с поставщиками и категориями
	if result := r.db.Preload("Supplier").Preload("Categories").Find(&products); result.Error != nil {
		// Если ошибка — возвращаем nil и ошибку
		return nil, result.Error
	}
	// Иначе — возвращаем список продуктов
	return products, nil
}

// Обновляет продукт по ID
func (r *gormProductRepository) UpdateProductByID(productID uint, updatedProduct *database.Product) error {
	// Создаём переменную продукта
	var product database.Product
	// Находим продукт по ID с поставщиком и категориями
	if err := r.db.Preload("Categories").Preload("Supplier").First(&product, productID).Error; err != nil {
		// Если ошибка — возвращаем её
		return err
	}

	// Обновляем поля продукта
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	product.SupplierID = updatedProduct.SupplierID

	// Удаляем старые связи с категориями
	if err := r.db.Model(&product).Association("Categories").Clear(); err != nil {
		// Если ошибка — возвращаем её
		return err
	}

	// Если есть новые категории — добавляем их
	if len(updatedProduct.Categories) > 0 {
		// Создаём срез категорий
		var categories []database.Category
		// Проходим по новым категориям
		for _, cat := range updatedProduct.Categories {
			// Создаём переменную категории
			var category database.Category
			// Находим категорию по ID
			if err := r.db.First(&category, cat.ID).Error; err != nil {
				// Если ошибка — возвращаем её
				return err
			}
			// Добавляем категорию в срез
			categories = append(categories, category)
		}
		// Присваиваем продукту обновлённые категории
		product.Categories = categories
	}

	// Сохраняем изменения в БД
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&product).Error
}

// Удаляет продукт по ID
func (r *gormProductRepository) DeleteProductByID(productID uint) error {
	// Создаём переменную продукта
	var product database.Product
	// Находим продукт по ID с поставщиком и категориями
	if err := r.db.Preload("Categories").Preload("Supplier").First(&product, productID).Error; err != nil {
		// Если ошибка — возвращаем её
		return err
	}

	// Удаляем связи с категориями
	if err := r.db.Model(&product).Association("Categories").Clear(); err != nil {
		// Если ошибка — возвращаем её
		return err
	}

	// Удаляем продукт из БД (Hard Delete)
	if err := r.db.Unscoped().Delete(&product).Error; err != nil {
		// Если ошибка — возвращаем её
		return err
	}
	return nil
}

// Удаляет поставщика по ID
func (r *gormProductRepository) DeleteSupplierByID(supplierID uint) error {
	// Создаём переменную поставщика
	var supplier database.Supplier
	// Находим поставщика по ID
	if err := r.db.First(&supplier, supplierID).Error; err != nil {
		// Если ошибка — возвращаем её
		return err
	}
	// Удаляем поставщика из БД (Hard Delete)
	if err := r.db.Unscoped().Delete(&supplier).Error; err != nil {
		// Если ошибка — возвращаем её
		return err
	}
	return nil
}

// Удаляет категорию по ID
func (r *gormProductRepository) DeleteCategoryByID(categoryID uint) error {
	// Создаём переменную категории
	var category database.Category
	// Находим категорию по ID
	if err := r.db.First(&category, categoryID).Error; err != nil {
		// Если ошибка — возвращаем её
		return err
	}
	// Удаляем категорию из БД (Hard Delete)
	if err := r.db.Unscoped().Delete(&category).Error; err != nil {
		return err
	}
	return nil
}
