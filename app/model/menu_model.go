package model

type Menu struct {
	BaseModel
	Name       string `gorm:"type:varchar(255);not null;comment:菜单名称（唯一标识）" json:"name"`
	RouterPath string `gorm:"type:varchar(255);comment:前端路由路径" json:"routerPath"`
	Method     string `gorm:"type:varchar(10);comment:请求方法（GET/POST/PUT/DELETE等）" json:"method"`
	TenantID   string `gorm:"type:uuid;not null;comment:所属租户标识" json:"tenantId"`
	Tenant     Tenant `gorm:"foreignkey:TenantID" json:"-"`
	Auth       bool   `gorm:"default:false;comment:是否需要鉴权（0=不需要，1=需要）" json:"auth"`
	Icon       string `gorm:"type:varchar(255);default:'';comment:菜单图标（如：fa-user）" json:"icon"`
	Status     bool   `gorm:"default:false;comment:状态（1=启用，2=禁用）" json:"status"`
	ParentID   string `gorm:"type:uuid;comment:父级菜单ID（0=顶级菜单）" json:"parent_id"`
	Sort       int    `gorm:"not null;default:0;comment:排序（数字越小越靠前）" json:"sort"`
}

func (Menu) TableName() string {
	return "menu"
}
