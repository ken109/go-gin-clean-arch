package config

type ctxKey int

const (
	DBKey ctxKey = iota
	UIDKey
	ErrorKey
)
