package kdb

import (
	"context"
	"database/sql"
)

var kdb *engine

type engine struct {
	tablePrefix string
	structTag string
}

func RegisterDataBase(kConfig KConfig){
	for _,dbConfig :=range kConfig.DBConfigList{
		db,err := sql.Open(dbConfig.Driver,dbConfig.Dsn)
		if err!=nil{
			panic(err)
		}
		if dbConfig.MaxLifetime>0{
			db.SetConnMaxLifetime(dbConfig.MaxLifetime)
		}

		if dbConfig.MaxIdleConns>0{
			db.SetMaxIdleConns(dbConfig.MaxIdleConns)
		}

		if dbConfig.Name ==""{
			dbConfig.Name = defaultGroupName
		}

		m.addDB(dbConfig.Name,dbConfig.IsMaster,db)
	}
	kdb = new(engine)
	kdb.tablePrefix = kConfig.TablePrefix
	kdb.structTag = "db"
	if kConfig.StructTag!=""{
		kdb.structTag = kConfig.StructTag
	}
}

func Select(query string,bindings ...interface{}) *Rows{
	return newConnection().Select(query,bindings)
}

func Insert(query string,bindings ...interface{})(LastInsertId int64,err error){
	return newConnection().Insert(query,bindings)
}

func MultiInsert(query string,bindingsArr [][]interface{})(LastInsertId []int64,err error){
	return newConnection().MultiInsert(query,bindingsArr)
}

func Update(query string,bindings ...interface{})(RowsAffected int64,err error){
	return newConnection().Update(query,bindings)
}

func Delete(query string,bindings ...interface{})(RowsAffected int64,err error){
	return newConnection().Delete(query,bindings)
}

func WithDB(name string) *Connection{
	return newConnection().WithDB(name)
}

func WithContext(ctx context.Context) *Connection{
	return newConnection().WithContext(ctx)
}

func BeginTransaction()(conn *Connection,err error){
	conn = newConnection()
	err = conn.BeginTransaction()
	if err!=nil{
		return nil,err
	}
	return conn,nil
}

func Table(table string) *Builder{
	return newConnection().Table(table)
}