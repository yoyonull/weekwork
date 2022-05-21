package main

import (
	"database/sql"
	"fmt"
	"weekwork/firstweek/services"
)

// DB Db数据库连接池
var DB *sql.DB

func main() {
	err := InitDB()
	if err != nil {
		panic(err)
	}

	data := services.WorkQueryOne()
	fmt.Println(string(data))
	//services.WorkQuery()
	//services.WorkInsertUser()
	//services.WorkDeleteUser()

	defer DB.Close()
}
