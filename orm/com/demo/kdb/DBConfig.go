package kdb

import "time"

type DBConfig struct {
	Name string
	IsMaster bool
	Driver string
	Dsn string
	MaxLifetime time.Duration
	MaxIdleConns int
	MaxOpenConns int
}

type KConfig struct {
	TablePrefix string
	StructTag string
	DBConfigList []DBConfig
}
