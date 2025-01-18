package util

import (
	"context"

	"github.com/rs/xid"

	"go-gin-clean-arch/config"
)

func UID(ctx context.Context) xid.ID {
	db, _ := ctx.Value(config.UIDKey).(xid.ID)
	return db
}

func Authenticated(ctx context.Context) bool {
	return !UID(ctx).IsNil()
}
