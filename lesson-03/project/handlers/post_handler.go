package handlers

import (
	"gin-examples/project/models"
	"gin-examples/project/services"
	"gin-examples/project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostHandler struct {
	postService *services.PostService
}

// 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
// 实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
// 实现文章的更新功能，只有文章的作者才能更新自己的文章。
// 实现文章的删除功能，只有文章的作者才能删除自己的文章。
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}
	req.UserID = c.GetUint("userID")

	post, err := h.postService.CreatePost(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	})
}

func (h *PostHandler) ListPosts(c *gin.Context) {

	postList := make([]models.PostResponse, 0, 10)

	userID, _ := c.Get("userID")

	posts, err := h.postService.ListPosts(userID.(uint))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	for _, post := range *posts {
		postList = append(postList, models.PostResponse{ID: post.ID, Title: post.Title, Content: post.Content, CreatedAt: post.CreatedAt})
	}

	utils.Success(c, models.ListPostResponse{
		List: postList,
	})
}
func (h *PostHandler) GetPost(c *gin.Context) {
	var req models.GetPostRequest
	if err := c.ShouldBindUri(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}
	post, err := h.postService.GetPost(req.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}
	req.UserID = userID.(uint)

	postOld, err := h.postService.GetPost(req.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	if postOld.UserID != userID {
		utils.Error(c, 409, "没有权限查看这个文章")
		return
	}

	post, err := h.postService.UpdatePost(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	})
}

func (h *PostHandler) DelPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var req models.DelPostRequest
	if err := c.ShouldBindUri(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}

	postOld, err := h.postService.GetPost(req.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	if postOld.UserID != userID {
		utils.Error(c, 409, "没有权限删除这个文章")
		return
	}

	post, err := h.postService.DelPost(req.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	})
}
