package context

import (
	"gorm.io/gorm"
)

func (c *ctx) DB() *gorm.DB {
	if c.db != nil {
		return c.db
	}
	return c.getDB()
}

func (c *ctx) Transaction(fn func(ctx Context) error) error {
	return c.getDB().Transaction(
		func(tx *gorm.DB) error {
			c.db = tx

			return fn(c)
		},
	)
}
