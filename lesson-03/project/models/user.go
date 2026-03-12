package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=50"`
	Content string `json:"content" binding:"required,min=1,max=500"`
	UserID  uint
}
type UpdatePostRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required,min=1,max=50"`
	Content string `json:"content" binding:"required,min=1,max=500"`
	UserID  uint
}

type GetPostRequest struct {
	ID uint `uri:"id" binding:"required"`
}
type DelPostRequest struct {
	ID uint `uri:"id" binding:"required"`
}
type ListPostRequest struct {
	UserID uint `json:"id" binding:"required"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
type ListPostResponse struct {
	List []PostResponse
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
	PostID  uint   `json:"postId" binding:"required"`
	UserID  uint
}
type CreateCommentResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListCommentRequest struct {
	PostID uint `form:"postId" binding:"required"`
}

type CommentResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	UserID    uint
	User      User
	PostID    uint
	Post      Post
	CreatedAt time.Time `json:"createdAt"`
}
type ListCommentResponse struct {
	List []CommentResponse
}
