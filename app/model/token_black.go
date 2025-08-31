package model

type TokenBlack struct {
	BaseModel
	TokenUuid string
	Reason    string `gorm:"type:varchar(255)"`
	StarTime  int    `gorm:"" json:"starTime"`
	StopTime  int    `gorm:"" json:"stopTime"`
}
