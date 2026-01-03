package model

type Parking struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	ParkingNo string `json:"parking_no"` // 车位号 (如 A-101)
	Status    int    `json:"status"`     // 0:空闲 1:已占用
	UserID    int64  `json:"user_id"`    // 归属业主ID (0表示未绑定)
	CarPlate  string `json:"car_plate"`  // 绑定车牌号 (如 辽A88888)
	//UpdatedAt time.Time `json:"updated_at"` // 绑定时间
}

func (Parking) TableName() string {
	return "cms_parking"
}
