package model

type UserBlack struct {
	BaseModel
	UserID   uint   `json:"userId"`
	User     User   `gorm:"foreignkey:UserID" json:"_"`
	Type     int8   `gorm:"type:tinyint" json:"type"` //1.封禁2.风控3.异地登录
	Reason   string `gorm:"size:255"`
	StarTime int    `gorm:"type:int(11)" json:"starTime"`
	StopTime int    `gorm:"type:int(11)" json:"stopTime"`
}
