package model

type LogfileModel struct {
	Username string `gorm:"size:32" json:"username"`
	Password string `gorm:"size:64" json:"password"`
}
