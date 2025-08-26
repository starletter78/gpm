package model

type TokenBlack struct {
	BaseModel
	TokenUuid string
	Reason    string `gorm:"type:varchar(255)"`
	StarTime  int    `gorm:"type:int(11)" json:"starTime"`
	StopTime  int    `gorm:"type:int(11)" json:"stopTime"`
}
