package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"testing"
)

var host = "127.0.0.1"
var port = 3306
var username = "root"
var password = "xxx"
var database = "xxxx"

func TestFindOne(t *testing.T) {
	conn, err := newConn()
	if err != nil {
		t.Errorf("get conn failed, error: %v", err)
		return
	}

	model := NewCollectionUserModel(conn)
	one, err := model.FindOne(10166)
	if err != nil {
		t.Errorf("find collection failed, error: %v", err)
		return
	}

	fmt.Printf("one: %#v\n", one)
	fmt.Println(one.Cert)
}

// 创建数据库连接对象
func newConn() (*gorm.DB, error) {
	// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		username,
		password,
		host+":"+strconv.Itoa(port),
		database)

	conf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(mysql.Open(dsn), conf)
	if err != nil {
		return nil, err
	}

	return db, err
}
