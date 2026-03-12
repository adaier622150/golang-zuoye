package services

import (
	"gin-examples/project/models"
	"gin-examples/project/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

// 实现评论的创建功能，已认证的用户可以对文章发表评论。
// 实现评论的读取功能，支持获取某篇文章的所有评论列表。
func (s *CommentService) CreateComment(req models.CreateCommentRequest) (*models.Comment, error) {

	// 添加评论
	comment := models.Comment{
		Content: req.Content,
		PostID:  req.PostID,
		UserID:  req.UserID,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (s *CommentService) ListComments(req models.ListCommentRequest) (*[]models.Comment, error) {
	// 获取所评论
	comments := []models.Comment{}
	if err := s.db.Preload(clause.Associations).Where("post_id = ?", req.PostID).Find(&comments).Error; err != nil {
		return nil, utils.NewAppError(409, "查询文章所有评论列表报错")
	}
	return &comments, nil
}
