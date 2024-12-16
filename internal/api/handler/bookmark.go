package handler

import (
	"net/http"

	"github.com/OxytocinGroup/theca-backend/internal/domain"
	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg"
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
// @Summary Create a new user bookmark
// @Tags Bookmark
// @Accept json
// @Produce json
// @Param user body domain.Bookmark true "Bookmark"
// @Success 201 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /api/create-bookmark [post]
// @Security ApiKeyAuth
func (bh *BookmarkHandler) CreateBookmark(c *gin.Context) {
	var bookmark domain.Bookmark

	userID, _ := c.Get("user_id")
	bookmark.UserID = userID.(uint)
	if err := c.ShouldBindJSON(&bookmark); err != nil {
		bh.Logger.Info(c, "bad request", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	resp := bh.BookmarkUseCase.CreateBookmark(bookmark)
	c.JSON(resp.Code, resp)
}
