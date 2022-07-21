package data

import (
	"fmt"

	"gorm.io/gorm"
)

//构建View模型Belongs To Users
//constraint级联更新级联删除
type View struct {
	gorm.Model
	UserID int    `form:"userid"`
	View   string `form:"view"`
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

//获取看法信息
func CreateViews(db *gorm.DB, view *View) (err error) {
	err = db.Create(view).Error
	if err != nil {
		return err
	}
	return nil
}

//获取某一用户的所有看法
func GetViewsbyUserId(db *gorm.DB, views *[]View, id string) (err error) {
	err = db.Where("user_id = ?", id).Find(views).Error
	if err != nil {
		return err
	}
	return nil
}

//根据viewid删除view,每个用户只能删除自己所发的view
func DeleteViewsbyId(db *gorm.DB, view *View, id string, userid string) (err error, count int64) {
	db.Model(&View{}).Where("id = ? AND user_id = ?", id, userid).Count(&count)
	fmt.Println(count)
	err = db.Where("id = ? AND user_id = ?", id, userid).Delete(view).Error
	fmt.Println(count)
	if err != nil {
		return err, count
	}
	return nil, count
}

//验证用户密码是否正确
func TryLogin(db *gorm.DB, id string, password string) int {
	var count int64
	db.Model(&User{}).Where("ID = ? AND Password = ?", id, password).Count(&count)
	fmt.Printf("password:%s, ID: %s", password, id)
	return int(count)
}
