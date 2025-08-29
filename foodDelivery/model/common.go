package model

type COMMONMODEL struct {
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at" gorm:"timestamptz"`
	UpDatesAt int64  `json:"updated_at" gorm:"timestamptz"`
	TIMESTAMP int64  `json:"timestamp"`
}
