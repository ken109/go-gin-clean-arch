package domain

import (
	"context"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type SoftDeleteModel struct {
	ID        uint           `json:"-"`
	XID       xid.ID         `json:"id" gorm:"column:xid"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (base *SoftDeleteModel) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("XID", xid.New())
	return nil
}

type HardDeleteModel struct {
	ID        uint      `json:"-"`
	XID       xid.ID    `json:"id" gorm:"column:xid"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (base *HardDeleteModel) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("XID", xid.New())
	return nil
}
