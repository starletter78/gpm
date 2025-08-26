package model

type Role struct {
	BaseModel
	Name     string `gorm:"type:varchar(255);not null;comment:角色名称" json:"name"`
	TenantId string `gorm:"type:uuid;not null;comment:所属租户标识" json:"tenantId"`
	Tenant   Tenant `gorm:"foreignkey:TenantId" json:"-"`
}

func (Role) TableName() string {
	return "role"
}
