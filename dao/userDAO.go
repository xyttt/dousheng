package dao

import (
	"log"
)

// TableName 修改表名映射
func (tableUser users) TableName() string {
	return "users"
}

func WriteTable(username string, password string) (users, error) {
	newUser := users{

		Name:     username,
		Password: password,
	}
	return newUser, nil
}

// GetTableUserList 获取全部TableUser对象
func GetTableUserList() ([]users, error) {
	tableUsers := []users{}
	if err := DB.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

// GetTableUserByUsername 根据username获得TableUser对象
func GetTableUserByUsername(name string) (users, error) {
	tableUser := users{}
	if err := DB.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// GetTableUserById 根据user_id获得TableUser对象
func GetTableUserById(id int64) (users, error) {
	tableUser := users{}
	if err := DB.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// InsertTableUser 将tableUser插入表内
func InsertTableUser(tableUser *users) bool {
	if err := DB.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
