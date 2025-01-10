package requests

type BookmarkRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	URL      string `json:"url" binding:"required"`
	ShowText bool   `json:"show_text" binding:"required"`
}
