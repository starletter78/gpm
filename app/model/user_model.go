package model

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;type:varchar(255);comment:用户名（唯一标识）" json:"username"`
	Nickname string `gorm:"type:varchar(255);comment:昵称" json:"nickname"`
	Password string `gorm:"type:varchar(255);comment:密码（建议使用哈希加密存储）" json:"password"`
	Sex      bool   `gorm:"default:false;comment:性别 false=男 true=女" json:"sex"`
	Email    string `gorm:"uniqueIndex;type:varchar(255);comment:邮箱地址" json:"email"`
	Phone    string `gorm:"type:varchar(255);comment:手机号码" json:"phone"`
	Address  string `gorm:"type:varchar(255);comment:地址" json:"address"`
	Avatar   string `gorm:"type:varchar(255);comment:头像URL" json:"avatar"`
	Status   bool   `gorm:"not null;default:false;comment:状态 0=禁用 1=启用" json:"status"`
	Salt     string `gorm:"not null;type:uuid;comment:加密盐值（可选）" json:"salt"`
}

func (User) TableName() string {
	return "user"
}
