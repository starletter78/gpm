package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID       string `gorm:"type:uuid;primary_key;"`
	CreateAt int    `gorm:"autoUpdateTime;" json:"create_at"`
	UpdateAt int    `gorm:"autoCreateTime;" json:"update_at"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	b.ID = uid.String()
	return
}
