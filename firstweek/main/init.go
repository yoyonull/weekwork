package main

import (
	"database/sql"
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
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		// 从根本上获取堆栈的错误信息
		fmt.Printf("db info %T,%v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
		//return fmt.Errorf("db info %v:%w",err,err)
		// wrapper 嵌套使用
		return errors.Wrap(err, "db init failed")
	}
	return nil
}
