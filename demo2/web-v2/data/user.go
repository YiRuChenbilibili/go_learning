package data

import (
	"gorm.io/gorm"
)

type User struct {
	//注意开头字母要大写！！否则无法bind
	ID       uint   `gorm:"AUTO_INCREMENT=1" form:"id"`
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email"`
	Phonenum string `form:"phonenum"`
	Password string `gorm:"default:'123456'" form:"password"`
}

func CreateUser(db *gorm.DB, user *User) (err error) {
	err = db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers(db *gorm.DB, users *[]User) (err error) {
	err = db.Find(users).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserById(db *gorm.DB, user *User, id string) (err error) {
	err = db.Where("ID = ?", id).First(user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByName(db *gorm.DB, users *[]User, name string) (err error) {
	err = db.Where("Name = ?", name).Find(users).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(db *gorm.DB, user *User) (err error) {
	db.Save(user)
	return nil
}

func DeleteUsers(db *gorm.DB, user *User, id string) (err error) {
	db.Where("ID = ?", id).Delete(user)
	return nil

}
