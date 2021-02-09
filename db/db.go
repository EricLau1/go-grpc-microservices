package db

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

type Connection interface {
	Close()
	DB() *mgo.Database
}

type conn struct {
	session  *mgo.Session
	database *mgo.Database
}

func NewConnection(cfg Config) (Connection, error) {
	fmt.Println("database url:", cfg.Dsn())
	session, err := mgo.Dial(cfg.Dsn())
	if err != nil {
		return nil, err
	}
	return &conn{session: session, database: session.DB(cfg.DbName())}, nil
}

func (c *conn) Close() {
	c.session.Close()
}

func (c *conn) DB() *mgo.Database {
	return c.database
}
