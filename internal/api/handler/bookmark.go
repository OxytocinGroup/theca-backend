package handler

import (
	"net/http"

	"github.com/OxytocinGroup/theca-backend/internal/domain"
	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg"
	"github.com/OxytocinGroup/theca-backend/pkg/cerr"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	BookmarkUseCase usecase.BookmarkUseCase
	Logger          logger.Logger
}

func NewBookmarkHandler(usecase usecase.BookmarkUseCase, log logger.Logger) *BookmarkHandler {
	return &BookmarkHandler{
		BookmarkUseCase: usecase,
		Logger:          log,
	}
}

// @CreateBookmark GoDoc
// @Summary Create a bookmark
// @Description This endpoint allows an authenticated user to create a new bookmark.
// @Tags Bookmark
// @Accept  json
// @Produce  json
// @Param request body domain.Bookmark true "Bookmark details"
// @Success 201 {object} pkg.Response "Bookmark created successfully"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 401 {object} pkg.Response "Unauthorized - User not authenticated"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Security CookieAuth
// @Router /api/bookmarks/create [post]
func (bh *BookmarkHandler) CreateBookmark(c *gin.Context) {
	var bookmark domain.Bookmark

	userID, _ := c.Get("user_id")
	bookmark.UserID = userID.(uint)
	if err := c.ShouldBindJSON(&bookmark); err != nil {
		bh.Logger.Info(c, "bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{Code: http.StatusBadRequest, Message: ("Bad request " + err.Error()), Error: cerr.ErrInvalidBody})
		return
	}

	resp := bh.BookmarkUseCase.CreateBookmark(bookmark)
	c.JSON(resp.Code, resp)
}

// GetBookmarks godoc
// @Summary Get bookmarks by user ID
// @Description Fetch all bookmarks associated with the current user
// @Tags Bookmark
// @Produce json
// @Security CookieAuth
// @Success 200 {array} domain.Bookmark "List of bookmarks"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /api/bookmarks/get [get]
func (bh *BookmarkHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetUint("user_id")
	bookmarks, resp := bh.BookmarkUseCase.GetBookmarksByUser(userID)
	c.JSON(resp.Code, bookmarks)
}

// DeleteBookmark godoc
// @Summary Delete a bookmark by ID
// @Description Delete a specific bookmark associated with the user based on the bookmark ID
// @Tags Bookmark
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body domain.Bookmark true "Request body with the bookmark ID to delete"
// @Success 200 {object} pkg.Response "Successfully deleted the bookmark"
// @Failure 400 {object} pkg.Response "Bad request, invalid input"
// @Failure 403 {object} pkg.Response "Forbidden, the user does not have permission to delete this bookmark"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /api/bookmarks/delete [delete]
func (bh *BookmarkHandler) DeleteBookmark(c *gin.Context) {
	userID := c.GetUint("user_id")
	var bookmark domain.Bookmark

	if err := c.ShouldBindJSON(&bookmark); err != nil {
		bh.Logger.Info(c, "bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{Code: http.StatusBadRequest, Message: ("Bad request " + err.Error()), Error: cerr.ErrInvalidBody})
		return
	}

	resp := bh.BookmarkUseCase.DeleteBookmark(userID, bookmark.ID)
	c.JSON(resp.Code, resp)
}

// UpdateBookmark godoc
// @Summary Update a bookmark by ID
// @Description Update a specific bookmark associated with the user based on the bookmark ID
// @Tags Bookmark
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body domain.Bookmark true "Request body with the bookmark ID to update"
// @Success 200 {object} pkg.Response "Successfully updated the bookmark"
// @Failure 400 {object} pkg.Response "Bad request, invalid input"
// @Failure 403 {object} pkg.Response "Forbidden, the user does not have permission to update this bookmark"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /api/bookmarks/update [post]
func (bh *BookmarkHandler) UpdateBookmark(c *gin.Context) {
	userID := c.GetUint("user_id")
	var bookmark domain.Bookmark

	if err := c.ShouldBindJSON(&bookmark); err != nil {
		bh.Logger.Info(c, "bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{Code: http.StatusBadRequest, Message: ("Bad request " + err.Error()), Error: cerr.ErrInvalidBody})
		return
	}

	bookmark.UserID = userID
	resp := bh.BookmarkUseCase.UpdateBookmark(userID, &bookmark)
	c.JSON(resp.Code, resp)
}
