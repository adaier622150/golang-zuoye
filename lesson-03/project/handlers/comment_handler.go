package handlers

import (
	"gin-examples/project/models"
	"gin-examples/project/services"
	"gin-examples/project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentHandler struct {
	commentService *services.CommentService
}

// 实现评论的创建功能，已认证的用户可以对文章发表评论。
// 实现评论的读取功能，支持获取某篇文章的所有评论列表。
func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	req.UserID = userID.(uint)

	comment, err := h.commentService.CreateComment(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, models.CreateCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	})
}

func (h *CommentHandler) ListComments(c *gin.Context) {

	var req models.ListCommentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}
	commentList := make([]models.CommentResponse, 0, 10)

	comments, err := h.commentService.ListComments(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	for _, comment := range *comments {
		commentList = append(commentList, models.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			Post:      comment.Post,
			PostID:    comment.PostID,
			User:      comment.User,
			UserID:    comment.UserID,
		})
	}

	utils.Success(c, models.ListCommentResponse{
		List: commentList,
	})
}
