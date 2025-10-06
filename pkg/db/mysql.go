package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQLOptions ...
type MySQLOptions struct {
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Host            string        `mapstructure:"host"`
	Port            uint16        `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	CreateBatchSize int           `mapstructure:"createBatchSize"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime"`
}

// NewMySQLOptions ...
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Port:            3306,
		Database:        "vetes",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		CreateBatchSize: 1000,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Second,
	}
}

// Validate ...
func (o *MySQLOptions) Validate() error {
	return nil
}

// AddFlags ...
func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Username, "mysql-username", o.Username, "mysql db username")
	fs.StringVar(&o.Password, "mysql-password", o.Password, "mysql db password")
	fs.StringVar(&o.Host, "mysql-host", "", "mysql db host")
	fs.Uint16Var(&o.Port, "mysql-port", o.Port, "mysql db port")
	fs.StringVar(&o.Database, "mysql-database", o.Database, "mysql database name")
	fs.IntVar(&o.MaxIdleConns, "mysql-max-idle-conns", o.MaxIdleConns, "mysql max idle conns")
	fs.IntVar(&o.MaxOpenConns, "mysql-max-open-conns", o.MaxOpenConns, "mysql max open conns")
	fs.IntVar(&o.CreateBatchSize, "mysql-create-batch-size", o.CreateBatchSize, "mysql create batch size")
	fs.DurationVar(&o.ConnMaxIdleTime, "mysql-conn-max-idle-time", o.ConnMaxIdleTime, "mysql conn max idle time")
	fs.DurationVar(&o.ConnMaxLifetime, "mysql-conn-max-life-time", o.ConnMaxLifetime, "mysql conn max life time")
}

var (
	once sync.Once
	db   *gorm.DB
)

// GetGORMInstance ...
func (o *MySQLOptions) GetGORMInstance() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		db, err = o.newGormInstance()
	})
	return db, err
}

func (o *MySQLOptions) newGormInstance() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		o.Username, o.Password, o.Host, o.Port, o.Database,
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: o.CreateBatchSize,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	// set connection pool
	sqlDB.SetMaxIdleConns(o.MaxIdleConns)
	sqlDB.SetMaxOpenConns(o.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(o.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(o.ConnMaxIdleTime)

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return gormDB, nil
}
