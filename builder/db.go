package builder

import (
	"time"

	_ "github.com/go-sql-driver/mysql" // default use mysql

	"database/sql"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
)

// support native db, xorm, gorm
type DBBuilder struct {
	Driver          string
	Url             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (o DBBuilder) BuildDB() (*sql.DB, error) {
	db, err := sql.Open(o.Driver, o.Url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(o.MaxOpenConns)
	db.SetMaxIdleConns(o.MaxIdleConns)
	db.SetConnMaxLifetime(o.ConnMaxLifetime)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (o DBBuilder) BuildXORM() (*xorm.Engine, error) {
	db, err := xorm.NewEngine(o.Driver, o.Url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(o.MaxOpenConns)
	db.SetMaxIdleConns(o.MaxIdleConns)
	db.SetConnMaxLifetime(o.ConnMaxLifetime)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// recommended
func (o DBBuilder) BuildGORM() (*gorm.DB, error) {
	db, err := gorm.Open(o.Driver, o.Url)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxOpenConns(o.MaxOpenConns)
	db.DB().SetMaxIdleConns(o.MaxIdleConns)
	db.DB().SetConnMaxLifetime(o.ConnMaxLifetime)
	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
