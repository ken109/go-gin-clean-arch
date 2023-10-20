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
	ID        xid.ID         `json:"id" gorm:"type:varchar(20);primaryKey;autoIncrement:false"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (base *SoftDeleteModel) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("ID", xid.New())
	return nil
}

type HardDeleteModel struct {
	ID        xid.ID    `json:"id" gorm:"type:varchar(20);primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (base *HardDeleteModel) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("ID", xid.New())
	return nil
}
