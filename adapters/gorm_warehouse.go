package adapters

import (
	"warehouse-app/core"
	"warehouse-app/database"

	"gorm.io/gorm"
)

type gormWarehouseRepository struct {
	db *gorm.DB
}

func NewGormWarehouseRepository(db *gorm.DB) core.WarehouseRepository {
	return &gormWarehouseRepository{db: db}
}

// ── Зоны ──

func (r *gormWarehouseRepository) SaveZone(zone *database.Zone) error {
	return r.db.Create(zone).Error
}

func (r *gormWarehouseRepository) FindAllZones() ([]database.Zone, error) {
	var zones []database.Zone
	return zones, r.db.Preload("Cells").Find(&zones).Error
}

func (r *gormWarehouseRepository) FindZoneByID(id uint) (*database.Zone, error) {
	var zone database.Zone
	return &zone, r.db.Preload("Cells").First(&zone, id).Error
}

func (r *gormWarehouseRepository) DeleteZoneByID(id uint) error {
	return r.db.Unscoped().Delete(&database.Zone{}, id).Error
}

// ── Ячейки ──

func (r *gormWarehouseRepository) SaveCell(cell *database.Cell) error {
	return r.db.Create(cell).Error
}

func (r *gormWarehouseRepository) FindAllCells() ([]database.Cell, error) {
	var cells []database.Cell
	return cells, r.db.Preload("Zone").Find(&cells).Error
}

func (r *gormWarehouseRepository) FindCellByID(id uint) (*database.Cell, error) {
	var cell database.Cell
	return &cell, r.db.Preload("Zone").First(&cell, id).Error
}

func (r *gormWarehouseRepository) FindCellsByZoneID(zoneID uint) ([]database.Cell, error) {
	var cells []database.Cell
	return cells, r.db.Preload("Zone").Where("zone_id = ?", zoneID).Find(&cells).Error
}

func (r *gormWarehouseRepository) DeleteCellByID(id uint) error {
	return r.db.Unscoped().Delete(&database.Cell{}, id).Error
}

// ── Остатки ──

func (r *gormWarehouseRepository) FindAllStock() ([]database.Stock, error) {
	var stocks []database.Stock
	return stocks, r.db.Preload("Product").Preload("Cell").Preload("Cell.Zone").Find(&stocks).Error
}

func (r *gormWarehouseRepository) FindStockByCell(cellID uint) ([]database.Stock, error) {
	var stocks []database.Stock
	return stocks, r.db.Preload("Product").Preload("Cell").Where("cell_id = ?", cellID).Find(&stocks).Error
}

func (r *gormWarehouseRepository) FindStockByProduct(productID uint) ([]database.Stock, error) {
	var stocks []database.Stock
	return stocks, r.db.Preload("Cell").Where("product_id = ?", productID).Find(&stocks).Error
}

func (r *gormWarehouseRepository) UpsertStock(productID uint, cellID uint, delta int) error {
	var stock database.Stock
	err := r.db.Where("product_id = ? AND cell_id = ?", productID, cellID).First(&stock).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(&database.Stock{
			ProductID: productID,
			CellID:    cellID,
			Quantity:  delta,
		}).Error
	}
	if err != nil {
		return err
	}
	return r.db.Model(&stock).Update("quantity", stock.Quantity+delta).Error
}

// ── Приёмка ──

func (r *gormWarehouseRepository) SaveReceipt(receipt *database.Receipt) error {
	return r.db.Create(receipt).Error
}

func (r *gormWarehouseRepository) FindAllReceipts() ([]database.Receipt, error) {
	var receipts []database.Receipt
	return receipts, r.db.Preload("Supplier").Preload("Cell").Preload("Items").Preload("Items.Product").Find(&receipts).Error
}

func (r *gormWarehouseRepository) FindReceiptByID(id uint) (*database.Receipt, error) {
	var receipt database.Receipt
	return &receipt, r.db.Preload("Supplier").Preload("Cell").Preload("Items").Preload("Items.Product").First(&receipt, id).Error
}

func (r *gormWarehouseRepository) ConfirmReceipt(id uint) error {
	return r.db.Model(&database.Receipt{}).Where("id = ?", id).Update("status", database.StatusConfirmed).Error
}

func (r *gormWarehouseRepository) CancelReceipt(id uint) error {
	return r.db.Model(&database.Receipt{}).Where("id = ?", id).Update("status", database.StatusCancelled).Error
}

// ── Отгрузка ──

func (r *gormWarehouseRepository) SaveShipment(shipment *database.Shipment) error {
	return r.db.Create(shipment).Error
}

func (r *gormWarehouseRepository) FindAllShipments() ([]database.Shipment, error) {
	var shipments []database.Shipment
	return shipments, r.db.Preload("Cell").Preload("Items").Preload("Items.Product").Find(&shipments).Error
}

func (r *gormWarehouseRepository) FindShipmentByID(id uint) (*database.Shipment, error) {
	var shipment database.Shipment
	return &shipment, r.db.Preload("Cell").Preload("Items").Preload("Items.Product").First(&shipment, id).Error
}

func (r *gormWarehouseRepository) ConfirmShipment(id uint) error {
	return r.db.Model(&database.Shipment{}).Where("id = ?", id).Update("status", database.StatusConfirmed).Error
}

func (r *gormWarehouseRepository) CancelShipment(id uint) error {
	return r.db.Model(&database.Shipment{}).Where("id = ?", id).Update("status", database.StatusCancelled).Error
}

// ── Движения ──

func (r *gormWarehouseRepository) SaveMovement(movement *database.Movement) error {
	return r.db.Create(movement).Error
}

func (r *gormWarehouseRepository) FindAllMovements() ([]database.Movement, error) {
	var movements []database.Movement
	return movements, r.db.Preload("Product").Preload("FromCell").Preload("ToCell").
		Order("created_at desc").Find(&movements).Error
}

func (r *gormWarehouseRepository) FindMovementsByProduct(productID uint) ([]database.Movement, error) {
	var movements []database.Movement
	return movements, r.db.Preload("Product").Preload("FromCell").Preload("ToCell").
		Where("product_id = ?", productID).Order("created_at desc").Find(&movements).Error
}

// ── Инвентаризация ──

func (r *gormWarehouseRepository) SaveInventory(inventory *database.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *gormWarehouseRepository) FindAllInventories() ([]database.Inventory, error) {
	var inventories []database.Inventory
	return inventories, r.db.Preload("Zone").Preload("Items").Preload("Items.Product").Preload("Items.Cell").Find(&inventories).Error
}

func (r *gormWarehouseRepository) FindInventoryByID(id uint) (*database.Inventory, error) {
	var inventory database.Inventory
	return &inventory, r.db.Preload("Zone").Preload("Items").Preload("Items.Product").Preload("Items.Cell").First(&inventory, id).Error
}

func (r *gormWarehouseRepository) CloseInventory(id uint) error {
	return r.db.Model(&database.Inventory{}).Where("id = ?", id).Update("status", database.InventoryClosed).Error
}

// ── Обработка ──

func (r *gormWarehouseRepository) SaveProcessingOrder(order *database.ProcessingOrder) error {
	return r.db.Create(order).Error
}

func (r *gormWarehouseRepository) FindAllProcessingOrders() ([]database.ProcessingOrder, error) {
	var orders []database.ProcessingOrder
	return orders, r.db.Preload("Zone").Preload("Items").Preload("Items.Product").Find(&orders).Error
}

func (r *gormWarehouseRepository) FindProcessingOrderByID(id uint) (*database.ProcessingOrder, error) {
	var order database.ProcessingOrder
	return &order, r.db.Preload("Zone").Preload("Items").Preload("Items.Product").First(&order, id).Error
}

func (r *gormWarehouseRepository) ConfirmProcessingOrder(id uint) error {
	order, err := r.FindProcessingOrderByID(id)
	if err != nil {
		return err
	}
	for _, item := range order.Items {
		if err := r.db.Model(&database.ProcessingItem{}).
			Where("id = ?", item.ID).
			Update("lost", item.Scanned-item.Packed).Error; err != nil {
			return err
		}
	}
	return r.db.Model(&database.ProcessingOrder{}).Where("id = ?", id).Update("status", database.StatusConfirmed).Error
}
