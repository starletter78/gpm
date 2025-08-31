package model

type ActionLog struct {
	BaseModel
	LogID        string  `gorm:"type:uuid;not null;comment:日志唯一标识" json:"log_id"`
	UserID       string  `gorm:"type:uuid;comment:操作用户ID" json:"user_id"`
	User         User    `gorm:"foreignkey:UserID" json:"-"`
	IP           string  `gorm:"type:varchar(45);default:'';comment:IP地址" json:"ip"`
	UA           string  `gorm:"type:varchar(1024);default:'';comment:用户代理" json:"ua"`
	Action       string  `gorm:"type:varchar(255);default:'';comment:操作描述" json:"action"`
	Path         string  `gorm:"type:varchar(255);default:'';comment:请求路径" json:"path"`
	Method       string  `gorm:"type:varchar(10);default:'';comment:请求方法" json:"method"`
	Tenant       string  `gorm:"type:varchar(255);default:default;comment:所属租户" json:"tenant"`
	Header       *string `gorm:"type:text;comment:请求头信息" json:"header,omitempty"`
	RequestBody  *string `gorm:"type:text;comment:请求体" json:"request_body,omitempty"`
	ResponseBody *string `gorm:"type:text;comment:响应体" json:"response_body,omitempty"`
	Status       int     `gorm:"default:0;comment:HTTP状态码" json:"status"`
	Duration     string  `gorm:"type:varchar(11);default:'0';comment:请求耗时（毫秒）" json:"duration"`
}

func (ActionLog) TableName() string {
	return "action_log"
}
