package usecase

import (
	"context"
	"net/http"

	"github.com/OxytocinGroup/theca-backend/internal/domain"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/pkg"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"github.com/OxytocinGroup/theca-backend/pkg/parsers"
)

type BookmarkUseCase interface {
	CreateBookmark(bookmark domain.Bookmark) pkg.Response
	GetBookmarksByUser(userID uint) ([]domain.Bookmark, pkg.Response)
}

type bookmarkUseCase struct {
	bookmarkRepo repository.BookmarkRepository
	log          logger.Logger
}

func NewBookmarkUseCase(bookmarkRepo repository.BookmarkRepository, log logger.Logger) BookmarkUseCase {
	return &bookmarkUseCase{
		bookmarkRepo: bookmarkRepo,
		log:          log,
	}
}

func (buc *bookmarkUseCase) CreateBookmark(bookmark domain.Bookmark) pkg.Response {
	iconURL, err := parsers.FetchFavicon(bookmark.URL)
	if err != nil {
		buc.log.Error(context.Background(), "failed to fetch favicom", map[string]any{"error": err})
	}
	bookmark.IconURL = iconURL
	err = buc.bookmarkRepo.CreateBookmark(&bookmark)
	if err != nil {
		buc.log.Error(context.Background(), "failed to create bookmark", map[string]any{"error": err})
		return pkg.Response{
			Code:    500,
			Message: "failed to create bookmark",
		}
	}

	return pkg.Response{
		Code:    http.StatusCreated,
		Message: "bookmark created successfully",
	}
}

func (buc *bookmarkUseCase) GetBookmarksByUser(userID uint) ([]domain.Bookmark, pkg.Response) {
	bookmarks, err := buc.bookmarkRepo.GetBookmarksByUser(userID)
	if err != nil {
		buc.log.Error(context.Background(), "failed to get bookmarks by user", map[string]any{
			"user_id": userID,
			"error":   err,
		})
		return nil, pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "failed to get bookmarks",
		}
	}
	return bookmarks, pkg.Response{
		Code: http.StatusOK,
	}
}
