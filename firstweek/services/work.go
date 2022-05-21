package services

import (
	"encoding/json"
	"firstweek/models"
	"fmt"
)

type JsonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func WorkQueryOne() {
	singleId := 5

	var info JsonResult
	info.Code = 200

	data, err := models.QueryOne(singleId)

	if err != nil {
		info.Code = 500
		info.Msg = data
		jsons, _ := json.Marshal(info)
		fmt.Printf("res%v ,%v\n", string(jsons), err)
		return
	}

	info.Msg = data
	jsons, _ := json.Marshal(info)
	fmt.Printf("res%v\n", string(jsons))
	return
}

func WorkQuery() {
	models.Query()
}

func WorkInsertUser() {
	models.InsertUser()

}

func WorkDeleteUser() {
	models.DeleteUser()
}
