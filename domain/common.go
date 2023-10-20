package domain

import (
	"time"

	"context"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type SoftDeleteModel struct {
	ID        xid.ID         `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type HardDeleteModel struct {
	ID        xid.ID    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (base *SoftDeleteModel) BeforeCreate(tx *gorm.DB) error {
	guid := xid.New()

	tx.Statement.SetColumn("ID", guid)
	return nil
}
