package services

import (
	"encoding/json"
	"weekwork/firstweek/models"
)

type JsonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// WorkQueryOne 作业案例
func (r *JsonResult) WorkQueryOne() []byte {
	singleId := 5

	var info JsonResult
	info.Code = 200

	data, err := models.QueryOne(singleId)
	if err != nil {
		info.Code = 500
		info.Msg = data
		jsons, _ := json.Marshal(info)
		return jsons
	}

	info.Msg = data
	jsons, _ := json.Marshal(info)
	return jsons
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
