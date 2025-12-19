// adapters/http_adapter.go
package adapters

import (
	"strconv"
	"warehouse-app/core"
	"warehouse-app/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HTTP-обработчик для продуктов
type httpProductHandler struct {
	service core.ProductService
}

func NewHttpProductHandler(service core.ProductService) *httpProductHandler {
	return &httpProductHandler{service: service}
}

// --- Роуты ---

func (h *httpProductHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/supplier", h.CreateSupplierFiber)
	app.Get("/supplier", h.GetAllSupplierFiber)
	app.Put("/supplier/:id", h.UpdateSupplierFiber)
	app.Delete("/supplier/:id", h.DeleteSupplierByIDFiber)

	app.Post("/category", h.CreateCategoryFiber)
	app.Get("/category", h.GetAllCategoryFiber)
	app.Delete("/category/:id", h.DeleteCategoryByIDFiber)

	app.Post("/product", h.CreateProductFiber)
	app.Get("/product", h.GetAllProductFiber)
	app.Get("/product/:id", h.GetProductByIDFiber)
	app.Put("/product/:id", h.UpdateProductByIDFiber)
	app.Delete("/product/:id", h.DeleteProductByIDFiber)

	app.Get("/category/:id/product", h.GetAllProductOfCategoryFiber)
	app.Get("/supplier/:id/product", h.GetAllProductOfSupplierFiber)
}

// --- Обработчики (оставляем как есть, только убираем комментарии для краткости) ---

func (h *httpProductHandler) CreateSupplierFiber(c *fiber.Ctx) error {
	var supplier database.Supplier
	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.service.CreateSupplier(&supplier); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(supplier)
}

func (h *httpProductHandler) CreateCategoryFiber(c *fiber.Ctx) error {
	var category database.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.service.CreateCategory(&category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(category)
}

func (h *httpProductHandler) GetAllSupplierFiber(c *fiber.Ctx) error {
	suppliers, err := h.service.GetAllSupplier()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(suppliers)
}

func (h *httpProductHandler) GetAllCategoryFiber(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(categories)
}

func (h *httpProductHandler) CreateProductFiber(c *fiber.Ctx) error {
	var product database.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.service.CreateProduct(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *httpProductHandler) GetProductByIDFiber(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	product, err := h.service.GetProductByID(uint(productID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *httpProductHandler) GetAllProductOfCategoryFiber(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	products, err := h.service.GetAllProductOfCategory(uint(categoryID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *httpProductHandler) GetAllProductOfSupplierFiber(c *fiber.Ctx) error {
	supplierID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
	}
	products, err := h.service.GetAllProductOfSupplier(uint(supplierID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *httpProductHandler) UpdateSupplierFiber(c *fiber.Ctx) error {
	supplierID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid supplier ID"})
	}
	var supplier database.Supplier
	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	supplier.ID = uint(supplierID)
	if err := h.service.UpdateSupplier(&supplier); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(supplier)
}

func (h *httpProductHandler) GetAllProductFiber(c *fiber.Ctx) error {
	products, err := h.service.GetAllProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *httpProductHandler) UpdateProductByIDFiber(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	var updatedProduct database.Product
	if err := c.BodyParser(&updatedProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := h.service.UpdateProductByID(uint(productID), &updatedProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(updatedProduct)
}

func (h *httpProductHandler) DeleteProductByIDFiber(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	if err := h.service.DeleteProductByID(uint(productID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
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

// --- НОВОЕ: Автомиграция ---

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&database.Product{},
		&database.Category{},
		&database.Supplier{},
		&database.ProductCategory{},
		&database.User{},
	)
}
