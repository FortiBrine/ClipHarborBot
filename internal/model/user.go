package model

type User struct {
	ID       int64 `gorm:"primaryKey"`
	Language string
}
