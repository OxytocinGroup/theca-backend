package repository

import (
	"github.com/OxytocinGroup/theca-backend/internal/domain"
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	CreateBookmark(bookmark *domain.Bookmark) error
	GetBookmarksByUser(userID uint) ([]domain.Bookmark, error)
	UpdateBookmark(bookmark *domain.Bookmark) error
	DeleteBookmarkByID(bookmarkID uint) error
	GetBookmarkOwner(bookmarkID uint) (uint, error)
}

type bookmarkDatabase struct {
	DB *gorm.DB
}

func NewBookmarkRepository(DB *gorm.DB) BookmarkRepository {
	return &bookmarkDatabase{DB}
}

func (bdb *bookmarkDatabase) CreateBookmark(bookmark *domain.Bookmark) error {
	return bdb.DB.Model(&domain.Bookmark{}).Create(bookmark).Error
}

func (bdb *bookmarkDatabase) GetBookmarksByUser(userID uint) ([]domain.Bookmark, error) {
	var results []domain.Bookmark
	err := bdb.DB.Model(&domain.Bookmark{}).Where("user_id = ?", userID).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (bdb *bookmarkDatabase) UpdateBookmark(bookmark *domain.Bookmark) error {
	return bdb.DB.Model(&domain.Bookmark{}).Where("id = ?", bookmark.ID).Save(bookmark).Error
}

func (bdb *bookmarkDatabase) DeleteBookmarkByID(bookmarkID uint) error {
	return bdb.DB.Model(&domain.Bookmark{}).Where("id = ?", bookmarkID).Delete(&domain.Bookmark{}).Error
}

func (bdb *bookmarkDatabase) GetBookmarkOwner(bookmarkID uint) (uint, error) {
	var userID uint
	err := bdb.DB.Model(&domain.Bookmark{}).Where("id = ?", bookmarkID).Pluck("user_id", &userID).Error
	return userID, err
}
