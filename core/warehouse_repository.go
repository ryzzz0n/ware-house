package core

import "warehouse-app/database"

type WarehouseRepository interface {
	// Зоны
	SaveZone(zone *database.Zone) error
	FindAllZones() ([]database.Zone, error)
	FindZoneByID(id uint) (*database.Zone, error)
	DeleteZoneByID(id uint) error

	// Ячейки
	SaveCell(cell *database.Cell) error
	FindAllCells() ([]database.Cell, error)
	FindCellByID(id uint) (*database.Cell, error)
	FindCellsByZoneID(zoneID uint) ([]database.Cell, error)
	DeleteCellByID(id uint) error

	// Остатки
	FindAllStock() ([]database.Stock, error)
	FindStockByCell(cellID uint) ([]database.Stock, error)
	FindStockByProduct(productID uint) ([]database.Stock, error)
	UpsertStock(productID uint, cellID uint, delta int) error

	// Приёмка
	SaveReceipt(receipt *database.Receipt) error
	FindAllReceipts() ([]database.Receipt, error)
	FindReceiptByID(id uint) (*database.Receipt, error)
	ConfirmReceipt(id uint) error
	CancelReceipt(id uint) error

	// Отгрузка
	SaveShipment(shipment *database.Shipment) error
	FindAllShipments() ([]database.Shipment, error)
	FindShipmentByID(id uint) (*database.Shipment, error)
	ConfirmShipment(id uint) error
	CancelShipment(id uint) error

	// Движения
	SaveMovement(movement *database.Movement) error
	FindAllMovements() ([]database.Movement, error)
	FindMovementsByProduct(productID uint) ([]database.Movement, error)

	// Инвентаризация
	SaveInventory(inventory *database.Inventory) error
	FindAllInventories() ([]database.Inventory, error)
	FindInventoryByID(id uint) (*database.Inventory, error)
	CloseInventory(id uint) error

	// Обработка (сканирование + упаковка)
	SaveProcessingOrder(order *database.ProcessingOrder) error
	FindAllProcessingOrders() ([]database.ProcessingOrder, error)
	FindProcessingOrderByID(id uint) (*database.ProcessingOrder, error)
	ConfirmProcessingOrder(id uint) error
}
