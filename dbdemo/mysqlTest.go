package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DbConn struct {
	Dsn string
	Db  *sql.DB
}

type UserTable struct {
	Uid        int
	Username   string
	Department string
	Created    string
}

func main() {
	var err error
	dbConn := DbConn{
		Dsn: "root:root@tcp(127.0.0.1:3306)/test?charset=utf8",
		Db:  nil,
	}
	dbConn.Db, err = sql.Open("mysql", dbConn.Dsn)
	if err != nil {
		panic(err)
		return
	}
	defer dbConn.Db.Close()

	//execData(&dbConn)

	//preExecData(&dbConn)

	result := dbConn.QueryRowData("select * from user_info where uid=(select max(uid) from user_info);")
	fmt.Println(result)
}

func execData(dbConn *DbConn) {
	result, err := dbConn.Db.Exec("insert user_info(username,departname,created) values ('Josh','bussiness group','2018-07-03')")
	if err != nil {
		fmt.Println(err)
	} else {
		count, _ := result.RowsAffected()
		fmt.Println("受影响行数：", count)
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("新添加数据的id：", id)
		}

	}
}

func preExecData(dbConn *DbConn) {
	stmt, err := dbConn.Db.Prepare("insert user_info(username,departname,created) values ('Josh','bussiness group','2018-07-03')")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		result, err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
		} else {
			count, _ := result.RowsAffected()
			fmt.Println("受影响行数：", count)
			id, _ := result.LastInsertId()
			fmt.Println("新添加数据的id：", id)
		}
	}
}

func (dbConn *DbConn) QueryRowData(sqlString string) (data UserTable) {
	user := new(UserTable)
	err := dbConn.Db.QueryRow(sqlString).Scan(&user.Uid, &user.Username, &user.Department, &user.Created)
	if err != nil {
		panic(err)
		return
	}
	return *user
}
