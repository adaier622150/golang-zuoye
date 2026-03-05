package main

import (
	"fmt"
	"gorm.io/gorm"
	"lesson02examples/testutil"
	"testing"
	"time"
)

//题目1：模型定义
//假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
//要求 ：
//使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
//Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
//编写Go代码，使用Gorm创建这些模型对应的数据库表。
//题目2：关联查询
//基于上述博客系统的模型定义。
//要求 ：
//编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
//编写Go代码，使用Gorm查询评论数量最多的文章信息。
//题目3：钩子函数
//继续使用博客系统的模型。
//要求 ：
//为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
//为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
//

type UserInfo struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64;not null"`
	Email     string `gorm:"size:128"`
	Age       uint8
	PostNum   uint8
	Status    string `gorm:"size:16;default:active;index"`
	Posts     []Post
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
type Post struct {
	ID            uint `gorm:"primaryKey"`
	UserInfoID    uint // Foreign key to user
	PostName      string
	Content       string
	Comments      []Comment
	CommentStatus string    `gorm:"default:我是默认值"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
func (a *Post) AfterCreate(tx *gorm.DB) error {
	UserInfoID := a.UserInfoID
	var count int64
	tx.Model(&Post{}).Where("user_info_id = ?", UserInfoID).Count(&count)
	tx.Model(&UserInfo{}).Where("id = ?", UserInfoID).Update("post_num", count)
	return nil
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint // Foreign key to Post
	Content   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (a *Comment) AfterDelete(tx *gorm.DB) error {
	PostID := a.PostID
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", PostID).Count(&count)
	if count <= 1 {
		tx.Model(&Post{}).Where("id = ?", PostID).Update("comment_status", "无评论")
	} else {
		tx.Model(&Post{}).Where("id = ?", PostID).Update("comment_status", "")
	}
	return nil
}

func main() {
	t := testing.T{}
	db := testutil.NewTestDB(&t, "testdb.db")

	initDb(&t, db)
	title2(&t, db, "Alice3")
	title3(&t, db)

}
func title2(t *testing.T, db *gorm.DB, userName string) {
	user := UserInfo{}
	if err := db.Preload("Posts.Comments").Where("name = ?", userName).First(&user).Error; err != nil {
		t.Fatalf("title2 : %v", err)
	}
	fmt.Println("查询某个用户发布的所有文章及其对应的评论信息:", user)
	type StatusCount struct {
		PostID uint
		Total  int64
	}
	var count StatusCount
	if err := db.Model(&Comment{}).Select("post_id, COUNT(*) as total").Group("post_id").Order("total desc").First(&count).Error; err != nil {
		t.Fatalf("统计查询最多评论文章 : %v", err)
	}
	post := Post{}
	if err := db.First(&post, count.PostID).Error; err != nil {
		t.Fatalf("查询文章 : %v", err)
	}
	fmt.Println("查询评论数量最多的文章信息。:", post)

}
func title3(t *testing.T, db *gorm.DB) {
	post := Post{}
	if err := db.Model(&Post{}).Take(&post).Error; err != nil {
		t.Fatalf("随便获取一条文章记录 : %v", err)
	}
	comment := Comment{PostID: post.ID}
	if err := db.Where("post_id = ?", post.ID).Delete(&comment).Error; err != nil {
		t.Fatalf("删除文章评论 : %v", err)
	}

}
func initDb(t *testing.T, db *gorm.DB) {

	if err := db.AutoMigrate(&UserInfo{}, &Post{}, &Comment{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// Clean up existing data
	db.Exec("PRAGMA foreign_keys = OFF")
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM user_infos")
	db.Exec("PRAGMA foreign_keys = ON")

	userInfos := []UserInfo{
		{
			Name:  "Alice1",
			Email: "alice@example.com",
			Age:   18,

			// Has Many: Orders association (nested with Has Many OrderItems)
			Posts: []Post{
				{
					PostName: "文章1",
					Content:  "文章内容1",
					// Has Many: OrderItems within Order
					Comments: []Comment{
						{
							Content: "评论1",
						},
						{
							Content: "评论2",
						},
					},
				},
			},
		},
		{
			Name:  "Alice2",
			Email: "alice@example.com",
			Age:   18,

			// Has Many: Orders association (nested with Has Many OrderItems)
			Posts: []Post{
				{
					PostName: "文章2",
					Content:  "文章内容2",
					// Has Many: OrderItems within Order
					Comments: []Comment{
						{
							Content: "评论1",
						},
						{
							Content: "评论2",
						},
					},
				},
			},
		},
		{
			Name:  "Alice3",
			Email: "alice@example.com",
			Age:   18,

			// Has Many: Orders association (nested with Has Many OrderItems)
			Posts: []Post{
				{
					PostName: "文章3",
					Content:  "文章内容3",
					// Has Many: OrderItems within Order
					Comments: []Comment{
						{
							Content: "评论1",
						},
						{
							Content: "评论2",
						},
					},
				},
				{
					PostName: "文章4",
					Content:  "文章内容4",
					// Has Many: OrderItems within Order
					Comments: []Comment{
						{
							Content: "评论1",
						},
						{
							Content: "评论2",
						},
						{
							Content: "评论3",
						},
					},
				},
			},
		},
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&userInfos).Error; err != nil {
		t.Fatalf("CreateUser : %v", err)
	}
}
