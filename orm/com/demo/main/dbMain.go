package main

import(
	"../kdb"
	"fmt"
	_"github.com/go-sql-driver/mysql"
)
func main(){
	kConf :=new(kdb.KConfig)
	dbConfig := new(kdb.DBConfig)
	dbConfig.Driver = "mysql"
	dbConfig.Dsn = "root:root@tcp(127.0.0.1:3306)/kdb?carset=utf8&parseTime=true"
	dbConfig.IsMaster = true
	kConf.DBConfigList = []kdb.DBConfig{*dbConfig}
	kdb.RegisterDataBase(*kConf)

	data,_:=kdb.Table("user").Where("name","nopsky").Get().ToArray()

	fmt.Println("data:",data)
}

