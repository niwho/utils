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
	if err := c.initdb(); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

func (db *DBClient) initdb() error {
	var err error
	db.RealDb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=100ms&readTimeout=50ms&writeTimeout=50ms", db.User, db.Password, db.Server, db.Name))
	return err
}
