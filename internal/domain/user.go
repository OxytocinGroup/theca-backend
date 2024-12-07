package domain

type User struct {
	ID               uint   `json:"id" gorm:"primary_key;unique;not null"`
	Email            string `json:"email" gorm:"size:255;unique;not null"`
	Username         string `json:"username" gorm:"size:255;unique;not null"`
	Password         string `json:"password" gorm:"size:255;not null"`
	IsVerified       bool   `json:"is_verified" gorm:"default:false"`
	VerificationCode string `json:"verification_code" gorm:"size:255"`
}
