package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

type Info struct {
	id    int64
	name  string
	age   int8
	sex   int8
	phone string
}

func QueryOne(singleId int) (string, error) {
	var name string

	err := DB.QueryRow("select name from info where id = ?", singleId).Scan(&name)

	//if err != nil  {
	//	log.Fatal(err)
	//}
	//fmt.Println(name)

	//switch {
	//	case err == sql.ErrNoRows:
	//		name = fmt.Sprintf("id is %d has no name.\n", singleId)
	//	case err != nil:
	//		log.Fatal(err)
	//	default:
	//		return name
	//}
	//return name

	if err != nil {
		if err == sql.ErrNoRows {
			//没有行，但也没有错误发生
			err = errors.Wrap(err, "sql: no rows in result set")
		} else {
			err = errors.Wrap(err, "sql query has wrong")
		}
	}

	return name, err

}

// Query 查询操作
func Query() {
	res, err := DB.Query("SELECT * FROM info where id >10")
	//sql.ErrNoRows
	defer res.Close()

	if err != nil {
		fmt.Printf("sql info %T,%v\n", errors.Cause(err), errors.Cause(err))
		log.Fatal(err)
	}

	for res.Next() {
		var fo Info
		err := res.Scan(&fo.id, &fo.name, &fo.age, &fo.sex, &fo.phone)

		if err != nil {
			fmt.Printf("sql info %T,%v\n", errors.Cause(err), errors.Cause(err))
			log.Fatal(err)
		}

		fmt.Printf("%v\n", fo)
	}
	//返回json
	//jsons, err := json.Marshal(res)
	//if err != nil {
	//	fmt.Println("err",err)
	//}
	//fmt.Println(string(jsons))

}

func DeleteUser() bool {

	//delete

	stmt, e := DB.Prepare("delete from info where id =?")
	ErrorCheck(e)

	res, e := stmt.Exec("24")
	ErrorCheck(e)
	// affected rows
	a, e := res.RowsAffected()
	ErrorCheck(e)

	fmt.Println(a) // 1
	return true
}

func InsertUser() bool {
	sql := "insert into info(name,age,sex,phone) values('bao',10,1,1121)"

	res, err := DB.Exec(sql)

	if err != nil {
		panic(err.Error())
	}
	lastId, err := res.LastInsertId()
	ErrorCheck(err)

	fmt.Printf("this is id: %d\n", lastId)

	return true
}
func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
