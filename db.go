package utils

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DBClient struct {
	Name     string
	Server   string
	User     string
	Password string

	RealDb *gorm.DB
}

func NewDBClient(name, server, user, passwd string) *DBClient {
	c := &DBClient{
		Name:     name,
		Server:   server,
		User:     user,
		Password: passwd,
	}
	if err := c.initdb(10, 250); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

// timeout ms
func NewDBClientV2(name, server, user, passwd string, maxConn int, timeout int) *DBClient {
	c := &DBClient{
		Name:     name,
		Server:   server,
		User:     user,
		Password: passwd,
	}
	if err := c.initdb(maxConn, timeout); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

func (db *DBClient) initdb(maxConn, timeout int) error {
	var err error
	db.RealDb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=100ms&readTimeout=%dms&writeTimeout=%dms", db.User, db.Password, db.Server, db.Name, timeout, timeout))
	db.RealDb.DB().SetMaxOpenConns(maxConn)
	return err
}
