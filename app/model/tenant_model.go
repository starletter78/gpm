package model

type Tenant struct {
	BaseModel
	Name string `gorm:"type:varchar(255);not null;comment:租户名称" json:"name"`
}

func (Tenant) TableName() string {
	return "tenant"
}
