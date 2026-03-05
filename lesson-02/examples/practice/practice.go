package main

import (
	"fmt"
	"gorm.io/gorm"
	"lesson02examples/testutil"
	"testing"
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:64;not null"`
	Email       string    `gorm:"size:128;uniqueIndex;not null"`
	Age         uint8     `gorm:"not null"`
	Status      string    `gorm:"size:16;default:active;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Phone       string    `gorm:"size:20;uniqueIndex"`
	LastLoginAt time.Time `gorm:""`
}

func main() {
	t := testing.T{}
	db := testutil.NewTestDB(&t, "testdb.db")

	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	type UserInfo struct {
		Name  string
		Email string
	}

	var users []UserInfo
	db.Model(&User{}).Where("status IN ?", []string{"active", "pending"}).Find(&users)
	fmt.Println(users)

}

// 新增用户：创建用户并默认开启激活状态
func CreateUser(db *gorm.DB, name, email string) (*User, error) {
	user := User{Name: name, Email: email}
	if err := db.Omit("Status").Create(&user).Error; err != nil {
		fmt.Println("CreateUser error: ", err)
		return nil, err
	}
	return &user, nil
}

// 模糊查询：根据邮箱模糊查询用户列表（支持分页）
func SearchUsersByEmail(db *gorm.DB, emailPattern string, page, size int) ([]User, error) {
	users := []User{}
	if err := db.Model(&User{}).Where("email LIKE ?", "%"+emailPattern+"%").Scopes(paginate(page, size)).Find(&users).Error; err != nil {
		fmt.Println("SearchUsersByEmail error: ", err)
		return nil, err
	}
	return users, nil
}

//批量更新状态：批量更新用户状态

func UpdateUserStatus(db *gorm.DB, ids []uint, status string) error {
	if err := db.Model(&User{}).Where("id in ?", ids).Updates(map[string]interface{}{"status": status}).Error; err != nil {
		fmt.Println("UpdateUserStatus error: ", err)
		return err
	}
	return nil
}

// 删除过期用户：删除超过 30 天未登录的用户
func DeleteInactiveUsers(db *gorm.DB) error {
	// 你的实现（注意：软删除将在进阶模块讲解）
	if err := db.Where("last_login_at < ?", time.Now().Add(-30*24*time.Hour)).Delete(&User{}).Error; err != nil {
		fmt.Println("DeleteInactiveUsers error: ", err)
		return err
	}
	return nil
}

func paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 10
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
