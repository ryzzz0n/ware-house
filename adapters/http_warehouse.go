package adapters

import (
	"strconv"
	"warehouse-app/core"
	"warehouse-app/database"

	"github.com/gofiber/fiber/v2"
)

type httpWarehouseHandler struct {
	service core.WarehouseService
}

func NewHttpWarehouseHandler(service core.WarehouseService) *httpWarehouseHandler {
	return &httpWarehouseHandler{service: service}
}

func (h *httpWarehouseHandler) RegisterRoutes(app *fiber.App) {
	// Зоны
	app.Post("/zone", h.CreateZone)
	app.Get("/zone", h.GetAllZones)
	app.Get("/zone/:id", h.GetZoneByID)
	app.Delete("/zone/:id", h.DeleteZone)

	// Ячейки
	app.Post("/cell", h.CreateCell)
	app.Get("/cell", h.GetAllCells)
	app.Get("/zone/:id/cell", h.GetCellsByZone)
	app.Delete("/cell/:id", h.DeleteCell)

	// Остатки
	app.Get("/stock", h.GetAllStock)
	app.Get("/stock/cell/:id", h.GetStockByCell)
	app.Get("/stock/product/:id", h.GetStockByProduct)

	// Приёмка
	app.Post("/receipt", h.CreateReceipt)
	app.Get("/receipt", h.GetAllReceipts)
	app.Get("/receipt/:id", h.GetReceiptByID)
	app.Post("/receipt/:id/confirm", h.ConfirmReceipt)
	app.Post("/receipt/:id/cancel", h.CancelReceipt)

	// Отгрузка
	app.Post("/shipment", h.CreateShipment)
	app.Get("/shipment", h.GetAllShipments)
	app.Get("/shipment/:id", h.GetShipmentByID)
	app.Post("/shipment/:id/confirm", h.ConfirmShipment)
	app.Post("/shipment/:id/cancel", h.CancelShipment)

	// Движения
	app.Get("/movement", h.GetAllMovements)
	app.Get("/movement/product/:id", h.GetMovementsByProduct)

	// Инвентаризация
	app.Post("/inventory", h.CreateInventory)
	app.Get("/inventory", h.GetAllInventories)
	app.Get("/inventory/:id", h.GetInventoryByID)
	app.Post("/inventory/:id/close", h.CloseInventory)

	// Обработка
	app.Post("/processing", h.CreateProcessingOrder)
	app.Get("/processing", h.GetAllProcessingOrders)
	app.Get("/processing/:id", h.GetProcessingOrderByID)
	app.Post("/processing/:id/confirm", h.ConfirmProcessingOrder)
}

// ── helpers ──

func parseID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func badID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
}

func notFound(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": msg})
}

func internalErr(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

// ── Зоны ──

func (h *httpWarehouseHandler) CreateZone(c *fiber.Ctx) error {
	var zone database.Zone
	if err := c.BodyParser(&zone); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if zone.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Название обязательно"})
	}
	if err := h.service.CreateZone(&zone); err != nil {
		return internalErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(zone)
}

func (h *httpWarehouseHandler) GetAllZones(c *fiber.Ctx) error {
	zones, err := h.service.GetAllZones()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(zones)
}

func (h *httpWarehouseHandler) GetZoneByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	zone, err := h.service.GetZoneByID(id)
	if err != nil {
		return notFound(c, "Участок не найден")
	}
	return c.JSON(zone)
}

func (h *httpWarehouseHandler) DeleteZone(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.DeleteZone(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Участок удалён"})
}

// ── Ячейки ──

func (h *httpWarehouseHandler) CreateCell(c *fiber.Ctx) error {
	var cell database.Cell
	if err := c.BodyParser(&cell); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if cell.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Код ячейки обязателен"})
	}
	if cell.ZoneID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Участок обязателен"})
	}
	if err := h.service.CreateCell(&cell); err != nil {
		return internalErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(cell)
}

func (h *httpWarehouseHandler) GetAllCells(c *fiber.Ctx) error {
	cells, err := h.service.GetAllCells()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(cells)
}

func (h *httpWarehouseHandler) GetCellsByZone(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	cells, err := h.service.GetCellsByZone(id)
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(cells)
}

func (h *httpWarehouseHandler) DeleteCell(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.DeleteCell(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Ячейка удалена"})
}

// ── Остатки ──

func (h *httpWarehouseHandler) GetAllStock(c *fiber.Ctx) error {
	stocks, err := h.service.GetAllStock()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(stocks)
}

func (h *httpWarehouseHandler) GetStockByCell(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	stocks, err := h.service.GetStockByCell(id)
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(stocks)
}

func (h *httpWarehouseHandler) GetStockByProduct(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	stocks, err := h.service.GetStockByProduct(id)
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(stocks)
}

// ── Приёмка ──

func (h *httpWarehouseHandler) CreateReceipt(c *fiber.Ctx) error {
	var receipt database.Receipt
	if err := c.BodyParser(&receipt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if receipt.SupplierID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Поставщик обязателен"})
	}
	if receipt.CellID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ячейка обязательна"})
	}
	if err := h.service.CreateReceipt(&receipt); err != nil {
		return internalErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(receipt)
}

func (h *httpWarehouseHandler) GetAllReceipts(c *fiber.Ctx) error {
	receipts, err := h.service.GetAllReceipts()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(receipts)
}

func (h *httpWarehouseHandler) GetReceiptByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	receipt, err := h.service.GetReceiptByID(id)
	if err != nil {
		return notFound(c, "Приёмка не найдена")
	}
	return c.JSON(receipt)
}

func (h *httpWarehouseHandler) ConfirmReceipt(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.ConfirmReceipt(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Приёмка проведена"})
}

func (h *httpWarehouseHandler) CancelReceipt(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.CancelReceipt(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Приёмка отменена"})
}

// ── Отгрузка ──

func (h *httpWarehouseHandler) CreateShipment(c *fiber.Ctx) error {
	var shipment database.Shipment
	if err := c.BodyParser(&shipment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if shipment.Recipient == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Получатель обязателен"})
	}
	if shipment.CellID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ячейка обязательна"})
	}
	if err := h.service.CreateShipment(&shipment); err != nil {
		return internalErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(shipment)
}

func (h *httpWarehouseHandler) GetAllShipments(c *fiber.Ctx) error {
	shipments, err := h.service.GetAllShipments()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(shipments)
}

func (h *httpWarehouseHandler) GetShipmentByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	shipment, err := h.service.GetShipmentByID(id)
	if err != nil {
		return notFound(c, "Отгрузка не найдена")
	}
	return c.JSON(shipment)
}

func (h *httpWarehouseHandler) ConfirmShipment(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.ConfirmShipment(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Отгрузка проведена"})
}

func (h *httpWarehouseHandler) CancelShipment(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.CancelShipment(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Отгрузка отменена"})
}

// ── Движения ──

func (h *httpWarehouseHandler) GetAllMovements(c *fiber.Ctx) error {
	movements, err := h.service.GetAllMovements()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(movements)
}

func (h *httpWarehouseHandler) GetMovementsByProduct(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	movements, err := h.service.GetMovementsByProduct(id)
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(movements)
}

// ── Инвентаризация ──

func (h *httpWarehouseHandler) CreateInventory(c *fiber.Ctx) error {
	var inventory database.Inventory
	if err := c.BodyParser(&inventory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if inventory.ZoneID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Участок обязателен"})
	}
	if err := h.service.CreateInventory(&inventory); err != nil {
		return internalErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(inventory)
}

func (h *httpWarehouseHandler) GetAllInventories(c *fiber.Ctx) error {
	inventories, err := h.service.GetAllInventories()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(inventories)
}

func (h *httpWarehouseHandler) GetInventoryByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	inventory, err := h.service.GetInventoryByID(id)
	if err != nil {
		return notFound(c, "Инвентаризация не найдена")
	}
	return c.JSON(inventory)
}

func (h *httpWarehouseHandler) CloseInventory(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.CloseInventory(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Инвентаризация закрыта"})
}

// ── Обработка ──

func (h *httpWarehouseHandler) CreateProcessingOrder(c *fiber.Ctx) error {
	var order database.ProcessingOrder
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}
	if order.ZoneID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Участок обязателен"})
	}
	if err := h.service.CreateProcessingOrder(&order); err != nil {
		return internalErr(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *httpWarehouseHandler) GetAllProcessingOrders(c *fiber.Ctx) error {
	orders, err := h.service.GetAllProcessingOrders()
	if err != nil {
		return internalErr(c, err)
	}
	return c.JSON(orders)
}

func (h *httpWarehouseHandler) GetProcessingOrderByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	order, err := h.service.GetProcessingOrderByID(id)
	if err != nil {
		return notFound(c, "Задание не найдено")
	}
	return c.JSON(order)
}

func (h *httpWarehouseHandler) ConfirmProcessingOrder(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badID(c)
	}
	if err := h.service.ConfirmProcessingOrder(id); err != nil {
		return internalErr(c, err)
	}
	return c.JSON(fiber.Map{"message": "Задание подтверждено"})
}
