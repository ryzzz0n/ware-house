package database

import (
	"time"

	"gorm.io/gorm"
)

// ── Пользователи ──

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `json:"-"`
}

// ── Справочники ──

type Supplier struct {
	gorm.Model
	Name    string `json:"name" gorm:"unique"`
	Contact string `json:"contact"`
}

type Category struct {
	gorm.Model
	Name     string    `json:"name" gorm:"unique"`
	Products []Product `json:"-" gorm:"many2many:product_categories;"`
}

type Product struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	SupplierID  uint       `json:"supplier_id"`
	Supplier    Supplier   `json:"supplier,omitempty"`
	Categories  []Category `json:"categories" gorm:"many2many:product_categories;"`
}

type ProductCategory struct {
	ProductID  uint
	CategoryID uint
}

// ── Склад: участки и ячейки ──

// ZoneType — тип участка
type ZoneType string

const (
	ZoneInbound  ZoneType = "inbound"  // приёмка
	ZoneScanning ZoneType = "scanning" // сканирование
	ZoneStorage  ZoneType = "storage"  // хранение (СХ)
	ZonePacking  ZoneType = "packing"  // упаковка
	ZoneOutbound ZoneType = "outbound" // отгрузка
)

type Zone struct {
	gorm.Model
	Name  string   `json:"name" gorm:"unique"`
	Type  ZoneType `json:"type"`
	Cells []Cell   `json:"cells,omitempty"`
}

type Cell struct {
	gorm.Model
	Code   string `json:"code" gorm:"unique"` // A1-01, B2-03
	ZoneID uint   `json:"zone_id"`
	Zone   Zone   `json:"zone,omitempty"`
}

// ── Остатки ──

// Stock — текущий остаток товара в ячейке
type Stock struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty"`
	CellID    uint    `json:"cell_id"`
	Cell      Cell    `json:"cell,omitempty"`
	Quantity  int     `json:"quantity"`
}

// ── Документы ──

type DocumentStatus string

const (
	StatusDraft     DocumentStatus = "draft"     // черновик
	StatusConfirmed DocumentStatus = "confirmed" // проведён
	StatusCancelled DocumentStatus = "cancelled" // отменён
)

// Receipt — приёмка товара
type Receipt struct {
	gorm.Model
	Number     string         `json:"number" gorm:"unique"`
	SupplierID uint           `json:"supplier_id"`
	Supplier   Supplier       `json:"supplier,omitempty"`
	CellID     uint           `json:"cell_id"`
	Cell       Cell           `json:"cell,omitempty"`
	Status     DocumentStatus `json:"status" gorm:"default:'draft'"`
	Date       time.Time      `json:"date"`
	Note       string         `json:"note"`
	Items      []ReceiptItem  `json:"items,omitempty"`
}

type ReceiptItem struct {
	gorm.Model
	ReceiptID uint    `json:"receipt_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Shipment — отгрузка товара
type Shipment struct {
	gorm.Model
	Number    string         `json:"number" gorm:"unique"`
	CellID    uint           `json:"cell_id"`
	Cell      Cell           `json:"cell,omitempty"`
	Status    DocumentStatus `json:"status" gorm:"default:'draft'"`
	Date      time.Time      `json:"date"`
	Note      string         `json:"note"`
	Recipient string         `json:"recipient"`
	Items     []ShipmentItem `json:"items,omitempty"`
}

type ShipmentItem struct {
	gorm.Model
	ShipmentID uint    `json:"shipment_id"`
	ProductID  uint    `json:"product_id"`
	Product    Product `json:"product,omitempty"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

// Movement — история движения товара (пишется автоматически)
type MovementType string

const (
	MovementIn       MovementType = "in"       // приход
	MovementOut      MovementType = "out"      // расход
	MovementTransfer MovementType = "transfer" // перемещение
	MovementStorno   MovementType = "storno"   // сторно
)

type Movement struct {
	gorm.Model
	ProductID  uint         `json:"product_id"`
	Product    Product      `json:"product,omitempty"`
	FromCellID *uint        `json:"from_cell_id"`
	FromCell   *Cell        `json:"from_cell,omitempty"`
	ToCellID   *uint        `json:"to_cell_id"`
	ToCell     *Cell        `json:"to_cell,omitempty"`
	Quantity   int          `json:"quantity"`
	Type       MovementType `json:"type"`
	RefID      uint         `json:"ref_id"`   // ID документа (receipt/shipment)
	RefType    string       `json:"ref_type"` // "receipt" / "shipment"
	Note       string       `json:"note"`
}

// ── Инвентаризация ──

type InventoryStatus string

const (
	InventoryOpen   InventoryStatus = "open"
	InventoryClosed InventoryStatus = "closed"
)

type Inventory struct {
	gorm.Model
	Number string          `json:"number" gorm:"unique"`
	ZoneID uint            `json:"zone_id"`
	Zone   Zone            `json:"zone,omitempty"`
	Status InventoryStatus `json:"status" gorm:"default:'open'"`
	Date   time.Time       `json:"date"`
	Note   string          `json:"note"`
	Items  []InventoryItem `json:"items,omitempty"`
}

// InventoryItem — позиция инвентаризации
// Expected — сколько должно быть по системе
// Actual   — сколько реально насчитали
// Diff     — расхождение (Actual - Expected)
type InventoryItem struct {
	gorm.Model
	InventoryID uint    `json:"inventory_id"`
	ProductID   uint    `json:"product_id"`
	Product     Product `json:"product,omitempty"`
	CellID      uint    `json:"cell_id"`
	Cell        Cell    `json:"cell,omitempty"`
	Expected    int     `json:"expected"`
	Actual      int     `json:"actual"`
	Diff        int     `json:"diff"`
}

// ProcessingOrder — задание на обработку (сканирование + упаковка)
type ProcessingOrder struct {
	gorm.Model
	Number string           `json:"number" gorm:"unique"`
	Status DocumentStatus   `json:"status" gorm:"default:'draft'"`
	Date   time.Time        `json:"date"`
	ZoneID uint             `json:"zone_id"`
	Zone   Zone             `json:"zone,omitempty"`
	Items  []ProcessingItem `json:"items,omitempty"`
}

// ProcessingItem — позиция задания
// Lost = Scanned - Packed, считается автоматически при подтверждении
type ProcessingItem struct {
	gorm.Model
	ProcessingOrderID uint    `json:"processing_order_id"`
	ProductID         uint    `json:"product_id"`
	Product           Product `json:"product,omitempty"`
	Scanned           int     `json:"scanned"`
	Packed            int     `json:"packed"`
	Lost              int     `json:"lost"`
}
