package user

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(ID int) (User, error)
	Update(user User) (User, error)
	FindAll() ([]User, error)

	// contract for web
	GetUsers(user DTJson) ([]UserOnWeb, error)
	GetTotalUser(user DTJson) (int, error)
	DeleteUserByID(ID int) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil

}
func (r *repository) FindByID(ID int) (User, error) {
	var user User
	err := r.db.Where("id = ?", ID).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil

}
func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
func (r *repository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}

func (r *repository) GetUsers(user DTJson) ([]UserOnWeb, error) {

	var usersOnWeb []UserOnWeb
	sql := "SELECT id,name,occupation,email,avatar_file_name from users WHERE 1=1 "

	if search := user.Search.Value; search != "" {
		sql = fmt.Sprintf("%s AND (name LIKE '%%%s%%' OR occupation LIKE '%%%s%%' OR email LIKE '%%%s%%') ", sql, search, search, search)
	}

	start := user.Start
	length := user.Length

	sql = fmt.Sprintf("%s LIMIT %d, %d", sql, start, length)

	err := r.db.Raw(sql).Scan(&usersOnWeb).Error
	if err != nil {
		return usersOnWeb, err
	}
	return usersOnWeb, nil
}

func (r *repository) GetTotalUser(user DTJson) (int, error) {
	var users []UserOnWeb

	sql := "SELECT id,name,occupation,email,avatar_file_name from users WHERE 1=1 "

	if search := user.Search.Value; search != "" {
		sql = fmt.Sprintf("%s AND (name LIKE '%%%s%%' OR occupation LIKE '%%%s%%' OR email LIKE '%%%s%%') ", sql, search, search, search)
	}

	err := r.db.Raw(sql).Scan(&users).Error
	if err != nil {
		return len(users), err
	}
	return len(users), nil
}

func (r *repository) DeleteUserByID(ID int) (User, error) {
	var user User
	err := r.db.Where("id = ?", ID).Delete(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
