package services

import (
	"errors"
	"gorm.io/gorm"

	"gin-examples/project/models"
	"gin-examples/project/utils"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
// 实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
// 实现文章的更新功能，只有文章的作者才能更新自己的文章。
// 实现文章的删除功能，只有文章的作者才能删除自己的文章。
func (s *PostService) CreatePost(req models.CreatePostRequest) (*models.Post, error) {
	// 检查用户名是否已存在
	var existingPost models.Post
	if err := s.db.Where("title = ?", req.Title).First(&existingPost).Error; err == nil {
		return nil, utils.NewAppError(409, "Title already exists")
	}

	// 创建用户
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	}

	if err := s.db.Create(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
func (s *PostService) UpdatePost(req models.UpdatePostRequest) (*models.Post, error) {
	var existingPost models.Post
	if err := s.db.Where("id = ?", req.ID).Where("user_id = ?", req.UserID).First(&existingPost).Error; err != nil {
		return nil, utils.NewAppError(409, "文章不存在")
	}
	// 修改
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
	}
	post.ID = req.ID

	if err := s.db.Model(&models.Post{}).Where("id = ?", req.ID).Updates(map[string]any{"content": req.Content, "title": req.Title}).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) ListPosts(userID uint) (*[]models.Post, error) {
	// 获取所有本人文章列表
	posts := []models.Post{}
	if err := s.db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, utils.NewAppError(409, "查询文章列表报错")
	}
	return &posts, nil
}

//func (s *UserService) GetUserByID(id uint) (*models.User, error) {
//	var user models.User
//	if err := s.db.First(&user, id).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, utils.NewAppError(404, "User not found")
//		}
//		return nil, err
//	}
//	return &user, nil
//}

func (s *PostService) GetPost(id uint) (*models.Post, error) {
	// 获取所有本人文章列表
	post := models.Post{}
	if err := s.db.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(404, "Post not found")
		}
		return nil, err
	}
	return &post, nil
}
func (s *PostService) DelPost(ID uint) (*models.Post, error) {
	// 删除单个文章
	post := models.Post{}
	post.ID = ID
	if err := s.db.Where("id = ?", ID).Delete(&post).Error; err != nil {
		return nil, utils.NewAppError(409, "查询单个文章的详细信息报错")
	}
	return &post, nil
}
