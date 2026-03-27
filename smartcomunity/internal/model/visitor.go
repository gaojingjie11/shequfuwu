package model

import "time"

type Visitor struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	UserID int64 `json:"user_id"`
	// 必须有 column:visitor_name
	Name string `json:"name" gorm:"column:visitor_name"`

	// 必须有 column:visitor_phone
	Mobile string `json:"mobile" gorm:"column:visitor_phone"`

	Reason    string    `json:"reason"`     // 来访原因
	VisitTime time.Time `json:"visit_time"` // 预计来访时间
	Status    int       `json:"status"`     // 0:待审核 1:通过 2:拒绝
	// 保留这个字段，用于存审核意见
	AuditRemark string `json:"audit_remark"`

	CreatedAt time.Time `json:"created_at"`
}

func (Visitor) TableName() string {
	return "cms_visitor"
}
