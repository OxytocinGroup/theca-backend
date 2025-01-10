package domain

type Bookmark struct {
	ID uint `json:"id" gorm:"primaryKey;not null;unique"`
	UserID uint `json:"user_id"`
	Title string `json:"title" gorm:"size:128"`
	URL string `json:"url" gorm:"size:255"`
	IconURL string `json:"icon_url" gorm:"size:255"`
	ShowText bool `json:"show_text" gorm:"default:false"`
}