package cfg

import (
	"github.com/xbitgo/components/database"
)

type DB struct {
	Type            string `json:"type" yaml:"type"`
	DSN             string `json:"dsn" yaml:"dsn"`
	MaxOpenConn     int    `json:"max_open_conn" yaml:"max_open_conn"`
	MaxIdleConn     int    `json:"max_idle_conn" yaml:"max_idle_conn"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`   //sec
	ConnMaxIdleTime int    `json:"conn_max_idle_time" yaml:"conn_max_idle_time"` //sec
}

// CreateInstance DI生成器默认调用
func (d *DB) CreateInstance() (*database.Database, error) {
	return database.NewDB(database.Config{
		Type:            d.Type,
		DSN:             d.DSN,
		MaxOpenConn:     d.MaxOpenConn,
		MaxIdleConn:     d.MaxIdleConn,
		ConnMaxLifetime: d.ConnMaxLifetime,
		ConnMaxIdleTime: d.ConnMaxIdleTime,
	})
}
