package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gin-examples/project/config"
	"gin-examples/project/handlers"
	"gin-examples/project/middleware"
	"gin-examples/project/models"
	"gin-examples/project/services"
	"gin-examples/project/utils"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := newMySQLDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	//fmt.Println("Database clear data")
	//// Clean up existing data
	//db.Exec("PRAGMA foreign_keys = OFF")
	//db.Exec("DELETE FROM comments")
	//db.Exec("DELETE FROM posts")
	//db.Exec("DELETE FROM users")
	//db.Exec("PRAGMA foreign_keys = ON")
	//fmt.Println("Database clear data completed")

	// 初始化服务
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService, []byte(cfg.JWT.Secret))

	postService := services.NewPostService(db)
	postHandler := handlers.NewPostHandler(postService)
	commentService := services.NewCommentService(db)
	commentHandler := handlers.NewCommentHandler(commentService)

	// 创建 Gin 引擎
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status": "ok",
		})
	})

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/users/register", userHandler.Register)
		public.POST("/users/login", userHandler.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth([]byte(cfg.JWT.Secret)))
	{

		user := protected.Group("/users")
		user.GET("/me", userHandler.GetProfile)
		user.PUT("/me", userHandler.UpdateProfile)

		post := protected.Group("/post")

		post.POST("", postHandler.CreatePost)
		post.GET("", postHandler.ListPosts)
		post.GET("/:id", postHandler.GetPost)
		post.PUT("", postHandler.UpdatePost)
		post.DELETE("/:id", postHandler.DelPost)

		comment := protected.Group("/comment")

		comment.POST("", commentHandler.CreateComment)
		comment.GET("", commentHandler.ListComments)
	}

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
func newMySQLDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: Set to logger.Info to see all SQL queries in development
		Logger: logger.Default.LogMode(logger.Info),

		// NamingStrategy: Customize table and column naming
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
			NoLowerCase:   false,
		},
	})
}
