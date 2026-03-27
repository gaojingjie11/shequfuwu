package model

import "time"

type AIReport struct {
	ID                 int64     `gorm:"primaryKey" json:"id"`
	RepairNewCount     int64     `gorm:"column:repair_new_count;not null;default:0" json:"repair_new_count"`
	RepairPendingCount int64     `gorm:"column:repair_pending_count;not null;default:0" json:"repair_pending_count"`
	VisitorNewCount    int64     `gorm:"column:visitor_new_count;not null;default:0" json:"visitor_new_count"`
	PropertyPaidCount  int64     `gorm:"column:property_paid_count;not null;default:0" json:"property_paid_count"`
	PropertyPaidAmount float64   `gorm:"column:property_paid_amount;type:decimal(10,2);not null;default:0.00" json:"property_paid_amount"`
	ReportSummary      string    `gorm:"column:report_summary;type:varchar(255)" json:"report_summary"`
	Report             string    `gorm:"column:report_markdown;type:longtext" json:"report"`
	GeneratedBy        int64     `gorm:"column:generated_by;not null;default:0" json:"generated_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (AIReport) TableName() string {
	return "cms_ai_report"
}
