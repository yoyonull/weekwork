package main

import (
	_ "firstweek/models"
	"firstweek/services"
)

func main() {
	err := InitDB()
	if err != nil {
		panic(err)
	}

	services.WorkQueryOne()
	//services.WorkQuery()
	//services.WorkInsertUser()
	//services.WorkDeleteUser()
	//defer DB.Close()
}
