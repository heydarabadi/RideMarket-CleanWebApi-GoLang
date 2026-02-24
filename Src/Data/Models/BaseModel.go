package Models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time `gorm:"index"`

	CreatedBy uint          `gorm:"not null;default:0"`
	UpdatedBy sql.NullInt64 `gorm:"default:null"`
	DeletedBy sql.NullInt64 `gorm:"default:null"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")

	var userId int = 0
	if value != nil {
		switch v := value.(type) {
		case int:
			userId = v
		case int64:
			userId = int(v)
		}
	}

	if userId > 0 {
		m.CreatedBy = uint(userId)
	} else {
		m.CreatedBy = 0
	}

	return nil
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")

	var userId sql.NullInt64
	if value != nil {
		switch v := value.(type) {
		case int:
			userId = sql.NullInt64{Valid: true, Int64: int64(v)}
		case int64:
			userId = sql.NullInt64{Valid: true, Int64: v}
		}
	}

	m.UpdatedBy = userId

	return nil
}

func (m *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")

	var userId sql.NullInt64
	if value != nil {
		switch v := value.(type) {
		case int:
			userId = sql.NullInt64{Valid: true, Int64: int64(v)}
		case int64:
			userId = sql.NullInt64{Valid: true, Int64: v}
		}
	}

	m.DeletedBy = userId

	return nil
}
