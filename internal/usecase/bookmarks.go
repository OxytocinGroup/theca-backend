package usecase

import (
	"context"
	"net/http"

	"github.com/OxytocinGroup/theca-backend/internal/domain"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/pkg"
	"github.com/OxytocinGroup/theca-backend/pkg/cerr"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"github.com/OxytocinGroup/theca-backend/pkg/parsers"
)

type BookmarkUseCase interface {
	CreateBookmark(bookmark domain.Bookmark) pkg.Response
	GetBookmarksByUser(userID uint) ([]domain.Bookmark, pkg.Response)
	DeleteBookmark(userID, bookmarkID uint) pkg.Response
	UpdateBookmark(userID uint, bookmark *domain.Bookmark) pkg.Response
}

type bookmarkUseCase struct {
	bookmarkRepo repository.BookmarkRepository
	userRepo     repository.UserRepository
	log          logger.Logger
}

func NewBookmarkUseCase(bookmarkRepo repository.BookmarkRepository, userRepo repository.UserRepository, log logger.Logger) BookmarkUseCase {
	return &bookmarkUseCase{
		bookmarkRepo: bookmarkRepo,
		userRepo:     userRepo,
		log:          log,
	}
}

func (buc *bookmarkUseCase) CreateBookmark(bookmark domain.Bookmark) pkg.Response {
	user, err := buc.userRepo.GetByID(bookmark.UserID)
	if err != nil {
		buc.log.Error(context.Background(), "Create bookmark: error while get user by id", map[string]any{"error": err, "user_id": bookmark.UserID})
		return pkg.Response{Code: http.StatusInternalServerError, Message: "Failed to get user"}
	}

	if user.AmountOfBookmarks >= 25 {
		buc.log.Info(context.Background(), "Create bookmark: limit of bookmarks for user", map[string]any{"user_id": user.ID})
		return pkg.Response{Code: http.StatusConflict, Message: "Limit of bookmarks: 25", Error: cerr.ErrLimitOfBookmarks}
	}

	user.AmountOfBookmarks += 1
	if err := buc.userRepo.Update(&user); err != nil {
		buc.log.Error(context.Background(), "Create bookmark: failed to update user", map[string]any{"error": err})
		return pkg.Response{Code: http.StatusInternalServerError, Message: "Failed to update user"}
	}

	err = buc.bookmarkRepo.CreateBookmark(&bookmark)
	if err != nil {
		buc.log.Error(context.Background(), "Create bookmark: failed to create bookmark", map[string]any{"error": err})
		return pkg.Response{
			Code:    500,
			Message: "Failed to create bookmark",
		}
	}

	go func() {
		iconURL, err := parsers.FetchFavicon(bookmark.URL)
		if err != nil {
			buc.log.Error(context.Background(), "Create bookmark: failed to fetch favicon", map[string]any{"error": err})
			return
		}
		bookmark.IconURL = iconURL
		if err := buc.bookmarkRepo.UploadBookmarkFavicon(bookmark.ID, iconURL); err != nil {
			buc.log.Error(context.Background(), "Create bookmark: failed to upload favicon url to bookmark", map[string]any{"error": err})
			return
		}
	}()

	buc.log.Info(context.Background(), "Create bookmark: created succesfully", map[string]any{})
	return pkg.Response{
		Code:    http.StatusCreated,
		Message: "bookmark created successfully",
	}
}

func (buc *bookmarkUseCase) GetBookmarksByUser(userID uint) ([]domain.Bookmark, pkg.Response) {
	bookmarks, err := buc.bookmarkRepo.GetBookmarksByUser(userID)
	if err != nil {
		buc.log.Error(context.Background(), "Get bookmarks by user: failed to get bookmarks by user", map[string]any{
			"user_id": userID,
			"error":   err,
		})
		return nil, pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "failed to get bookmarks",
		}
	}

	buc.log.Info(context.Background(), "Get bookmarks by user: success", map[string]any{})
	return bookmarks, pkg.Response{
		Code: http.StatusOK,
	}
}

func (buc *bookmarkUseCase) DeleteBookmark(userID, bookmarkID uint) pkg.Response {
	user, err := buc.userRepo.GetByID(userID)
	if err != nil {
		buc.log.Error(context.Background(), "Delete bookmark: error while get user by id", map[string]any{"error": err, "user_id": userID})
		return pkg.Response{Code: http.StatusInternalServerError, Message: "Failed to get user"}
	}

	user.AmountOfBookmarks -= 1
	if err := buc.userRepo.Update(&user); err != nil {
		buc.log.Error(context.Background(), "Delete bookmark: failed to update user", map[string]any{"error": err})
		return pkg.Response{Code: http.StatusInternalServerError, Message: "Failed to update user"}
	}

	bookmarkOwner, err := buc.bookmarkRepo.GetBookmarkOwner(bookmarkID)
	if err != nil {
		buc.log.Error(context.Background(), "Delete bookmark: failed to get bookmark owner", map[string]any{
			"bookmarkID": bookmarkID,
			"error":      err,
		})
		return pkg.Response{
			Code:    500,
			Message: "failed to get bookmark owner",
		}
	}

	if userID != bookmarkOwner {
		buc.log.Info(context.Background(), "Delete bookmark: bookmark belongs to another user", map[string]any{
			"userID":     userID,
			"ownerID":    bookmarkOwner,
			"bookmarkID": bookmarkID,
		})
		return pkg.Response{
			Code:    http.StatusForbidden,
			Message: "bookmark belongs to another user",
			Error:   cerr.BelongsToAnotherUser,
		}
	}

	err = buc.bookmarkRepo.DeleteBookmarkByID(bookmarkID)
	if err != nil {
		buc.log.Error(context.Background(), "Delete bookmark: failed to delete bookmark", map[string]any{
			"bookmarkID": bookmarkID,
			"error":      err,
		})
		return pkg.Response{
			Code:    500,
			Message: "failed to delete bookmark",
		}
	}

	buc.log.Info(context.Background(), "Delete bookmark: success", map[string]any{})
	return pkg.Response{
		Code:    200,
		Message: "Bookmark deleted",
	}
}

func (buc *bookmarkUseCase) UpdateBookmark(userID uint, bookmark *domain.Bookmark) pkg.Response {
	go func() {
		iconURL, err := parsers.FetchFavicon(bookmark.URL)
		if err != nil {
			buc.log.Error(context.Background(), "Update bookmark: failed to fetch favicon", map[string]any{"error": err})
			return
		}
		if iconURL == "" {
			buc.log.Warn(context.Background(), "Update bookmark: empty icon url", map[string]any{
				"bookmark url": bookmark.URL,
			})
			return
		}

		bookmark.IconURL = iconURL
		if err := buc.bookmarkRepo.UploadBookmarkFavicon(bookmark.ID, iconURL); err != nil {
			buc.log.Error(context.Background(), "Update bookmark: failed to upload favicon url to bookmark", map[string]any{"error": err})
			return
		}
	}()

	err := buc.bookmarkRepo.UpdateBookmark(bookmark)
	if err != nil {
		buc.log.Error(context.Background(), "Update bookmark: failed to update bookmark", map[string]any{
			"bookmarkID": bookmark.ID,
			"error":      err,
		})
		return pkg.Response{
			Code:    500,
			Message: "failed to update bookmark",
		}
	}
	return pkg.Response{
		Code: 200,
	}
}
