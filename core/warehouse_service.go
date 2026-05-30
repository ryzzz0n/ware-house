package core

import (
	"errors"
	"fmt"
	"time"
	"warehouse-app/database"
)

type WarehouseService interface {
	// Зоны
	CreateZone(zone *database.Zone) error
	GetAllZones() ([]database.Zone, error)
	GetZoneByID(id uint) (*database.Zone, error)
	DeleteZone(id uint) error

	// Ячейки
	CreateCell(cell *database.Cell) error
	GetAllCells() ([]database.Cell, error)
	GetCellsByZone(zoneID uint) ([]database.Cell, error)
	DeleteCell(id uint) error

	// Остатки
	GetAllStock() ([]database.Stock, error)
	GetStockByCell(cellID uint) ([]database.Stock, error)
	GetStockByProduct(productID uint) ([]database.Stock, error)

	// Приёмка
	CreateReceipt(receipt *database.Receipt) error
	GetAllReceipts() ([]database.Receipt, error)
	GetReceiptByID(id uint) (*database.Receipt, error)
	ConfirmReceipt(id uint) error
	CancelReceipt(id uint) error

	// Отгрузка
	CreateShipment(shipment *database.Shipment) error
	GetAllShipments() ([]database.Shipment, error)
	GetShipmentByID(id uint) (*database.Shipment, error)
	ConfirmShipment(id uint) error
	CancelShipment(id uint) error

	// Движения
	GetAllMovements() ([]database.Movement, error)
	GetMovementsByProduct(productID uint) ([]database.Movement, error)

	// Инвентаризация
	CreateInventory(inventory *database.Inventory) error
	GetAllInventories() ([]database.Inventory, error)
	GetInventoryByID(id uint) (*database.Inventory, error)
	CloseInventory(id uint) error

	// Обработка
	CreateProcessingOrder(order *database.ProcessingOrder) error
	GetAllProcessingOrders() ([]database.ProcessingOrder, error)
	GetProcessingOrderByID(id uint) (*database.ProcessingOrder, error)
	ConfirmProcessingOrder(id uint) error
}

type warehouseServiceImpl struct {
	repo WarehouseRepository
}

func NewWarehouseService(repo WarehouseRepository) WarehouseService {
	return &warehouseServiceImpl{repo: repo}
}

// ── Зоны ──

func (s *warehouseServiceImpl) CreateZone(zone *database.Zone) error {
	return s.repo.SaveZone(zone)
}

func (s *warehouseServiceImpl) GetAllZones() ([]database.Zone, error) {
	return s.repo.FindAllZones()
}

func (s *warehouseServiceImpl) GetZoneByID(id uint) (*database.Zone, error) {
	return s.repo.FindZoneByID(id)
}

func (s *warehouseServiceImpl) DeleteZone(id uint) error {
	return s.repo.DeleteZoneByID(id)
}

// ── Ячейки ──

func (s *warehouseServiceImpl) CreateCell(cell *database.Cell) error {
	return s.repo.SaveCell(cell)
}

func (s *warehouseServiceImpl) GetAllCells() ([]database.Cell, error) {
	return s.repo.FindAllCells()
}

func (s *warehouseServiceImpl) GetCellsByZone(zoneID uint) ([]database.Cell, error) {
	return s.repo.FindCellsByZoneID(zoneID)
}

func (s *warehouseServiceImpl) DeleteCell(id uint) error {
	return s.repo.DeleteCellByID(id)
}

// ── Остатки ──

func (s *warehouseServiceImpl) GetAllStock() ([]database.Stock, error) {
	return s.repo.FindAllStock()
}

func (s *warehouseServiceImpl) GetStockByCell(cellID uint) ([]database.Stock, error) {
	return s.repo.FindStockByCell(cellID)
}

func (s *warehouseServiceImpl) GetStockByProduct(productID uint) ([]database.Stock, error) {
	return s.repo.FindStockByProduct(productID)
}

// ── Приёмка ──

func (s *warehouseServiceImpl) CreateReceipt(receipt *database.Receipt) error {
	receipt.Status = database.StatusDraft
	receipt.Date = time.Now()
	if receipt.Number == "" {
		receipt.Number = fmt.Sprintf("RCP-%d", time.Now().UnixNano())
	}
	return s.repo.SaveReceipt(receipt)
}

func (s *warehouseServiceImpl) GetAllReceipts() ([]database.Receipt, error) {
	return s.repo.FindAllReceipts()
}

func (s *warehouseServiceImpl) GetReceiptByID(id uint) (*database.Receipt, error) {
	return s.repo.FindReceiptByID(id)
}

func (s *warehouseServiceImpl) ConfirmReceipt(id uint) error {
	receipt, err := s.repo.FindReceiptByID(id)
	if err != nil {
		return err
	}
	if receipt.Status != database.StatusDraft {
		return errors.New("можно подтвердить только черновик")
	}
	if len(receipt.Items) == 0 {
		return errors.New("нельзя подтвердить пустую приёмку")
	}
	// обновляем остатки и пишем движения
	for _, item := range receipt.Items {
		if err := s.repo.UpsertStock(item.ProductID, receipt.CellID, item.Quantity); err != nil {
			return err
		}
		if err := s.repo.SaveMovement(&database.Movement{
			ProductID: item.ProductID,
			ToCellID:  &receipt.CellID,
			Quantity:  item.Quantity,
			Type:      database.MovementIn,
			RefID:     receipt.ID,
			RefType:   "receipt",
			Note:      fmt.Sprintf("Приёмка %s", receipt.Number),
		}); err != nil {
			return err
		}
	}
	return s.repo.ConfirmReceipt(id)
}

func (s *warehouseServiceImpl) CancelReceipt(id uint) error {
	receipt, err := s.repo.FindReceiptByID(id)
	if err != nil {
		return err
	}
	if receipt.Status == database.StatusCancelled {
		return errors.New("приёмка уже отменена")
	}
	// если была проведена — откатываем остатки
	if receipt.Status == database.StatusConfirmed {
		for _, item := range receipt.Items {
			if err := s.repo.UpsertStock(item.ProductID, receipt.CellID, -item.Quantity); err != nil {
				return err
			}
			if err := s.repo.SaveMovement(&database.Movement{
				ProductID:  item.ProductID,
				FromCellID: &receipt.CellID,
				Quantity:   item.Quantity,
				Type:       database.MovementStorno,
				RefID:      receipt.ID,
				RefType:    "receipt",
				Note:       fmt.Sprintf("Сторно приёмки %s", receipt.Number),
			}); err != nil {
				return err
			}
		}
	}
	return s.repo.CancelReceipt(id)
}

// ── Отгрузка ──

func (s *warehouseServiceImpl) CreateShipment(shipment *database.Shipment) error {
	shipment.Status = database.StatusDraft
	shipment.Date = time.Now()
	if shipment.Number == "" {
		shipment.Number = fmt.Sprintf("SHP-%d", time.Now().UnixNano())
	}
	return s.repo.SaveShipment(shipment)
}

func (s *warehouseServiceImpl) GetAllShipments() ([]database.Shipment, error) {
	return s.repo.FindAllShipments()
}

func (s *warehouseServiceImpl) GetShipmentByID(id uint) (*database.Shipment, error) {
	return s.repo.FindShipmentByID(id)
}

func (s *warehouseServiceImpl) ConfirmShipment(id uint) error {
	shipment, err := s.repo.FindShipmentByID(id)
	if err != nil {
		return err
	}
	if shipment.Status != database.StatusDraft {
		return errors.New("можно подтвердить только черновик")
	}
	if len(shipment.Items) == 0 {
		return errors.New("нельзя подтвердить пустую отгрузку")
	}
	for _, item := range shipment.Items {
		// проверяем остаток
		stocks, err := s.repo.FindStockByProduct(item.ProductID)
		if err != nil {
			return err
		}
		total := 0
		for _, st := range stocks {
			total += st.Quantity
		}
		if total < item.Quantity {
			return fmt.Errorf("недостаточно товара ID=%d: есть %d, нужно %d", item.ProductID, total, item.Quantity)
		}
		if err := s.repo.UpsertStock(item.ProductID, shipment.CellID, -item.Quantity); err != nil {
			return err
		}
		if err := s.repo.SaveMovement(&database.Movement{
			ProductID:  item.ProductID,
			FromCellID: &shipment.CellID,
			Quantity:   item.Quantity,
			Type:       database.MovementOut,
			RefID:      shipment.ID,
			RefType:    "shipment",
			Note:       fmt.Sprintf("Отгрузка %s", shipment.Number),
		}); err != nil {
			return err
		}
	}
	return s.repo.ConfirmShipment(id)
}

func (s *warehouseServiceImpl) CancelShipment(id uint) error {
	shipment, err := s.repo.FindShipmentByID(id)
	if err != nil {
		return err
	}
	if shipment.Status == database.StatusCancelled {
		return errors.New("отгрузка уже отменена")
	}
	if shipment.Status == database.StatusConfirmed {
		for _, item := range shipment.Items {
			if err := s.repo.UpsertStock(item.ProductID, shipment.CellID, item.Quantity); err != nil {
				return err
			}
			if err := s.repo.SaveMovement(&database.Movement{
				ProductID: item.ProductID,
				ToCellID:  &shipment.CellID,
				Quantity:  item.Quantity,
				Type:      database.MovementStorno,
				RefID:     shipment.ID,
				RefType:   "shipment",
				Note:      fmt.Sprintf("Сторно отгрузки %s", shipment.Number),
			}); err != nil {
				return err
			}
		}
	}
	return s.repo.CancelShipment(id)
}

// ── Движения ──

func (s *warehouseServiceImpl) GetAllMovements() ([]database.Movement, error) {
	return s.repo.FindAllMovements()
}

func (s *warehouseServiceImpl) GetMovementsByProduct(productID uint) ([]database.Movement, error) {
	return s.repo.FindMovementsByProduct(productID)
}

// ── Инвентаризация ──

func (s *warehouseServiceImpl) CreateInventory(inventory *database.Inventory) error {
	inventory.Status = database.InventoryOpen
	inventory.Date = time.Now()
	if inventory.Number == "" {
		inventory.Number = fmt.Sprintf("INV-%d", time.Now().UnixNano())
	}
	// заполняем ожидаемые остатки по ячейкам зоны
	stocks, err := s.repo.FindAllStock()
	if err != nil {
		return err
	}
	cells, err := s.repo.FindCellsByZoneID(inventory.ZoneID)
	if err != nil {
		return err
	}
	cellIDs := make(map[uint]bool)
	for _, c := range cells {
		cellIDs[c.ID] = true
	}
	for _, st := range stocks {
		if cellIDs[st.CellID] {
			inventory.Items = append(inventory.Items, database.InventoryItem{
				ProductID: st.ProductID,
				CellID:    st.CellID,
				Expected:  st.Quantity,
				Actual:    st.Quantity, // по умолчанию совпадает
			})
		}
	}
	return s.repo.SaveInventory(inventory)
}

func (s *warehouseServiceImpl) GetAllInventories() ([]database.Inventory, error) {
	return s.repo.FindAllInventories()
}

func (s *warehouseServiceImpl) GetInventoryByID(id uint) (*database.Inventory, error) {
	return s.repo.FindInventoryByID(id)
}

func (s *warehouseServiceImpl) CloseInventory(id uint) error {
	inventory, err := s.repo.FindInventoryByID(id)
	if err != nil {
		return err
	}
	if inventory.Status == database.InventoryClosed {
		return errors.New("инвентаризация уже закрыта")
	}
	// применяем расхождения — корректируем остатки
	for _, item := range inventory.Items {
		diff := item.Actual - item.Expected
		if diff != 0 {
			if err := s.repo.UpsertStock(item.ProductID, item.CellID, diff); err != nil {
				return err
			}
			moveType := database.MovementIn
			if diff < 0 {
				moveType = database.MovementOut
				diff = -diff
			}
			if err := s.repo.SaveMovement(&database.Movement{
				ProductID: item.ProductID,
				ToCellID:  &item.CellID,
				Quantity:  diff,
				Type:      moveType,
				RefID:     inventory.ID,
				RefType:   "inventory",
				Note:      fmt.Sprintf("Корректировка инвентаризации %s", inventory.Number),
			}); err != nil {
				return err
			}
		}
	}
	return s.repo.CloseInventory(id)
}

// ── Обработка (сканирование + упаковка) ──

func (s *warehouseServiceImpl) CreateProcessingOrder(order *database.ProcessingOrder) error {
	order.Status = database.StatusDraft
	order.Date = time.Now()
	if order.Number == "" {
		order.Number = fmt.Sprintf("PRC-%d", time.Now().UnixNano())
	}
	return s.repo.SaveProcessingOrder(order)
}

func (s *warehouseServiceImpl) GetAllProcessingOrders() ([]database.ProcessingOrder, error) {
	return s.repo.FindAllProcessingOrders()
}

func (s *warehouseServiceImpl) GetProcessingOrderByID(id uint) (*database.ProcessingOrder, error) {
	return s.repo.FindProcessingOrderByID(id)
}

func (s *warehouseServiceImpl) ConfirmProcessingOrder(id uint) error {
	order, err := s.repo.FindProcessingOrderByID(id)
	if err != nil {
		return err
	}
	if order.Status != database.StatusDraft {
		return errors.New("можно подтвердить только черновик")
	}
	// считаем потери и пишем движения
	for i := range order.Items {
		item := &order.Items[i]
		item.Lost = item.Scanned - item.Packed
		if item.Lost < 0 {
			return fmt.Errorf("упаковано больше чем отсканировано для товара ID=%d", item.ProductID)
		}
		if item.Lost > 0 {
			if err := s.repo.SaveMovement(&database.Movement{
				ProductID: item.ProductID,
				Quantity:  item.Lost,
				Type:      database.MovementStorno,
				RefID:     order.ID,
				RefType:   "processing",
				Note:      fmt.Sprintf("Потери при обработке %s: отсканировано %d, упаковано %d", order.Number, item.Scanned, item.Packed),
			}); err != nil {
				return err
			}
		}
	}
	return s.repo.ConfirmProcessingOrder(id)
}
