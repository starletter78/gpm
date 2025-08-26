package model

type Api struct {
	BaseModel
	Name     string `gorm:"type:varchar(128);not null;comment:菜单名称（唯一标识）" json:"name"`
	Path     string `gorm:"type:varchar(255);default:'';comment:API路径（用于权限匹配）" json:"path"`
	Method   string `gorm:"type:varchar(10);default:'';comment:请求方法（GET/POST/PUT/DELETE等）" json:"method"`
	TenantID string `gorm:"type:uuid;not null;comment:所属租户标识" json:"tenantId"`
	Tenant   Tenant `gorm:"foreignkey:TenantID" json:"-"`
	Auth     bool   `gorm:"default:false;comment:是否需要鉴权" json:"auth"`
	Status   bool   `gorm:"comment:状态（1=启用，2=禁用）" json:"status"`
	MenuID   string `gorm:"type:uuid;comment:父级菜单ID（0=顶级菜单）" json:"menuId"`
	Menu     Menu   `gorm:"foreignkey:MenuID" json:"-"`
}

func (Api) TableName() string {
	return "api"
}
