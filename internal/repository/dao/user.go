package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail      = errors.New("邮箱重复")
	InvalidEmailOrPassword = errors.New("邮箱或密码不正确")
	ErrUserNotFound        = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

type User struct {
	Id       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"unique;not_null" json:"email"`
	Password string `json:"-"`
	Ctime    int64  `json:"ctime"`
	Utime    int64  `json:"utime"`
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	return u, dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
}
