package vobj

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go-gin-clean-arch/config"
	"go-gin-clean-arch/packages/cerrors"
	"go-gin-clean-arch/packages/crypto"
)

type RecoveryToken string

func NewRecoveryToken(recoveryToken string) *RecoveryToken {
	value := RecoveryToken(recoveryToken)
	return &value
}

func (p *RecoveryToken) Generate() (time.Duration, time.Time, error) {
	duration := config.RecoveryTokenExpire
	expire := time.Now().Add(duration)
	token, err := crypto.EncryptCTR(
		expire.Format(time.RFC3339),
		config.Env.App.Secret,
	)
	if err != nil {
		return 0, time.Time{}, cerrors.NewUnexpected(err)
	}

	*p = RecoveryToken(token)

	return duration, expire, nil
}

func (p *RecoveryToken) IsValid() bool {
	if len(string(*p)) < 16 {
		return false
	}

	decrypted, err := crypto.DecryptCTR(string(*p), config.Env.App.Secret)
	if err != nil {
		return false
	}

	expire, err := time.Parse(time.RFC3339, decrypted)
	if err != nil {
		return false
	}

	return time.Now().Before(expire)
}

func (p *RecoveryToken) String() string {
	return string(*p)
}

func (p *RecoveryToken) Clear() {
	*p = ""
}

// sql

func (p *RecoveryToken) Scan(value interface{}) error {
	nullString := &sql.NullString{}
	err := nullString.Scan(value)
	if err != nil {
		return cerrors.NewUnexpected(err)
	}

	*p = RecoveryToken(nullString.String)

	return nil
}

func (p *RecoveryToken) Value() (driver.Value, error) {
	return string(*p), nil
}

// GormDataType gorm common data type
func (p *RecoveryToken) GormDataType() string {
	return "recovery_token"
}

// GormDBDataType gorm db data type
func (p *RecoveryToken) GormDBDataType(_ *gorm.DB, _ *schema.Field) string {
	return "varchar(256)"
}

// json

func (p *RecoveryToken) MarshalJSON() ([]byte, error) {
	return []byte("\"" + *p + "\""), nil
}

func (p *RecoveryToken) UnmarshalJSON(b []byte) error {
	*p = RecoveryToken(b)
	return nil
}
