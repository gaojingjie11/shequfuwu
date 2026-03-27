package model

type Parking struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	ParkingNo string `gorm:"type:varchar(32);index" json:"parking_no"`
	Status    int    `gorm:"not null;default:0" json:"status"`
	UserID    int64  `gorm:"not null;default:0" json:"user_id"`
	CarPlate  string `gorm:"type:varchar(32)" json:"car_plate"`
}

func (Parking) TableName() string {
	return "cms_parking"
}
