package adapters

import (
	"strconv"
	"warehouse-app/core"
	"warehouse-app/database"

	"github.com/gofiber/fiber/v2"
)

// HTTP-обработчик для продуктов
type httpProductHandler struct {
	service core.ProductService
}

// Создаёт новый HTTP-обработчик
func NewHttpProductHandler(service core.ProductService) *httpProductHandler {
	// Возвращаем указатель на новый экземпляр обработчика
	return &httpProductHandler{service: service}
}

// Создаёт поставщика
func (h *httpProductHandler) CreateSupplierFiber(c *fiber.Ctx) error {
	// Создаём переменную поставщика
	var supplier database.Supplier
	// Парсим тело запроса в структуру поставщика
	if err := c.BodyParser(&supplier); err != nil {
		// Если ошибка — возвращаем 400 и сообщение
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Вызываем метод сервиса для создания поставщика
	if err := h.service.CreateSupplier(&supplier); err != nil {
		// Если ошибка — возвращаем 500 и сообщение
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 201 и созданного поставщика
	return c.Status(fiber.StatusCreated).JSON(supplier)
}

// Создаёт категорию
func (h *httpProductHandler) CreateCategoryFiber(c *fiber.Ctx) error {
	// Создаём переменную категории
	var category database.Category
	// Парсим тело запроса в структуру категории
	if err := c.BodyParser(&category); err != nil {
		// Если ошибка — возвращаем 400 и сообщение
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Вызываем метод сервиса для создания категории
	if err := h.service.CreateCategory(&category); err != nil {
		// Если ошибка — возвращаем 500 и сообщение
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 201 и созданную категорию
	return c.Status(fiber.StatusCreated).JSON(category)
}

// Возвращает всех поставщиков
func (h *httpProductHandler) GetAllSupplierFiber(c *fiber.Ctx) error {
	// Вызываем метод сервиса для получения всех поставщиков
	suppliers, err := h.service.GetAllSupplier()
	// Если ошибка — возвращаем 500 и сообщение
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и список поставщиков
	return c.Status(fiber.StatusOK).JSON(suppliers)
}

// Возвращает все категории
func (h *httpProductHandler) GetAllCategoryFiber(c *fiber.Ctx) error {
	// Вызываем метод сервиса для получения всех категорий
	categories, err := h.service.GetAllCategory()
	// Если ошибка — возвращаем 500 и сообщение
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и список категорий
	return c.Status(fiber.StatusOK).JSON(categories)
}

// Создаёт продукт
func (h *httpProductHandler) CreateProductFiber(c *fiber.Ctx) error {
	// Создаём переменную продукта
	var product database.Product
	// Парсим тело запроса в структуру продукта
	if err := c.BodyParser(&product); err != nil {
		// Если ошибка — возвращаем 400 и сообщение
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	// Вызываем метод сервиса для создания продукта
	if err := h.service.CreateProduct(&product); err != nil {
		// Если ошибка — возвращаем 500 и сообщение
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 201 и созданный продукт
	return c.Status(fiber.StatusCreated).JSON(product)
}

// Возвращает продукт по ID
func (h *httpProductHandler) GetProductByIDFiber(c *fiber.Ctx) error {
	// Конвертируем параметр ID из строки в число
	productID, err := strconv.Atoi(c.Params("id"))
	// Если ошибка — возвращаем 400
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// Вызываем метод сервиса для получения продукта по ID
	product, err := h.service.GetProductByID(uint(productID))
	// Если ошибка — возвращаем 500 и сообщение
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и продукт
	return c.Status(fiber.StatusOK).JSON(product)
}

// Возвращает все продукты по ID категории
func (h *httpProductHandler) GetAllProductOfCategoryFiber(c *fiber.Ctx) error {
	// Конвертируем параметр ID из строки в число
	categoryID, err := strconv.Atoi(c.Params("id"))
	// Если ошибка — возвращаем 400
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// Вызываем метод сервиса для получения продуктов по ID категории
	products, err := h.service.GetAllProductOfCategory(uint(categoryID))
	// Если ошибка — возвращаем 500 и сообщение
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и список продуктов
	return c.Status(fiber.StatusOK).JSON(products)
}

// Возвращает все продукты по ID поставщика
func (h *httpProductHandler) GetAllProductOfSupplierFiber(c *fiber.Ctx) error {
	// Конвертируем параметр ID из строки в число
	supplierID, err := strconv.Atoi(c.Params("id"))
	// Если ошибка — возвращаем 400
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
	}
	// Вызываем метод сервиса для получения продуктов по ID поставщика
	products, err := h.service.GetAllProductOfSupplier(uint(supplierID))
	// Если ошибка — возвращаем 500 и сообщение
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и список продуктов
	return c.Status(fiber.StatusOK).JSON(products)
}

// Обновляет поставщика
func (h *httpProductHandler) UpdateSupplierFiber(c *fiber.Ctx) error {
	// Конвертируем параметр ID из строки в число
	supplierID, err := strconv.Atoi(c.Params("id"))
	// Если ошибка — возвращаем 400 и сообщение
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
	}
	// Создаём переменную поставщика
	var supplier database.Supplier
	// Парсим тело запроса в структуру поставщика
	if err := c.BodyParser(&supplier); err != nil {
		// Если ошибка — возвращаем 400 и сообщение
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// Присваиваем ID поставщику
	supplier.ID = uint(supplierID)
	// Вызываем метод сервиса для обновления поставщика
	if err := h.service.UpdateSupplier(&supplier); err != nil {
		// Если ошибка — возвращаем 500 и сообщение
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и обновлённого поставщика
	return c.Status(fiber.StatusOK).JSON(supplier)
}

// Возвращает все продукты
func (h *httpProductHandler) GetAllProductFiber(c *fiber.Ctx) error {
	// Вызываем метод сервиса для получения всех продуктов
	products, err := h.service.GetAllProduct()
	// Если ошибка — возвращаем 500 и сообщение
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и список продуктов
	return c.Status(fiber.StatusOK).JSON(products)
}

// Обновляет продукт по ID
func (h *httpProductHandler) UpdateProductByIDFiber(c *fiber.Ctx) error {
	// Конвертируем параметр ID из строки в число
	productID, err := strconv.Atoi(c.Params("id"))
	// Если ошибка — возвращаем 400 и сообщение
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	// Создаём переменную обновлённого продукта
	var updatedProduct database.Product
	// Парсим тело запроса в структуру продукта
	if err := c.BodyParser(&updatedProduct); err != nil {
		// Если ошибка — возвращаем 400 и сообщение
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// Вызываем метод сервиса для обновления продукта по ID
	if err := h.service.UpdateProductByID(uint(productID), &updatedProduct); err != nil {
		// Если ошибка — возвращаем 500 и сообщение
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и обновлённый продукт
	return c.Status(fiber.StatusOK).JSON(updatedProduct)
}

// Удаляет продукт по ID
func (h *httpProductHandler) DeleteProductByIDFiber(c *fiber.Ctx) error {
	// Конвертируем параметр ID из строки в число
	productID, err := strconv.Atoi(c.Params("id"))
	// Если ошибка — возвращаем 400 и сообщение
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	// Вызываем метод сервиса для удаления продукта по ID
	if err := h.service.DeleteProductByID(uint(productID)); err != nil {
		// Если ошибка — возвращаем 500 и сообщение
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Возвращаем 200 и сообщение об успешном удалении
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully deleted product by ID"})
}

func (h *httpProductHandler) DeleteSupplierByIDFiber(c *fiber.Ctx) error {
	supplierID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
	}
	if err := h.service.DeleteSupplierByID(uint(supplierID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully deleted supplier by ID"})
}

func (h *httpProductHandler) DeleteCategoryByIDFiber(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}
	if err := h.service.DeleteCategoryByID(uint(categoryID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully deleted category by ID"})
}
