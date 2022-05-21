package main

import (
	"database/sql"
	"firstweek/models"
	_ "firstweek/models"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"strings"
)

//数据库配置
const (
	userName = "root"
	password = "yoyo1033"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "test"
)

// InitDB 初始化
func InitDB() error {

	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	db, _ := sql.Open("mysql", path)

	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil {
		// 从根本上获取堆栈的错误信息
		fmt.Printf("db info %T,%v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
		// wrapper 嵌套使用
		return errors.Wrap(err, "db init failed")
	}
	models.DB = db
	return nil
}
