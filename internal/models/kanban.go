package models

import (
	"time"
)

// Board represents a Kanban board
type Board struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	UserID    string `gorm:"type:uuid;not null"` // 新增，關聯 User
	Position  int    `gorm:"not null;default:0"` // 新增 position 欄位
	CreatedAt time.Time
	UpdatedAt time.Time
	Lists     []List `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// List represents a list in a Kanban board
type List struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	BoardID   uint   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Position  int    `gorm:"not null;default:0"`
	Cards     []Card `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Card represents a card in a Kanban list
type Card struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string
	ListID    uint `gorm:"not null"`
	BoardID   uint `gorm:"not null"` // 新增 BoardID 欄位
	CreatedAt time.Time
	UpdatedAt time.Time
	Position  int `gorm:"not null;default:0"`
}
