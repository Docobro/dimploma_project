package driver

import (
	"fmt"
	"sync"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type DatabaseConnectionParams struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}
type Driver struct {
	once   sync.Once
	dbLock sync.Mutex
	conn   *gorm.DB
}

func New() *Driver {
	return &Driver{}
}

func (d *Driver) Connect(connectionParams DatabaseConnectionParams) (c *gorm.DB, err error) {
	d.once.Do(func() {
		d.dbLock.Lock()
		defer d.dbLock.Unlock()

		if d.conn != nil {
			return
		}

		dsn := fmt.Sprintf("clickhouse://%v:%v@%v:%v/%v?dial_timeout=10s&read_timeout=20s",
			connectionParams.Username,
			connectionParams.Password,
			connectionParams.Host,
			connectionParams.Port,
			connectionParams.Database,
		)

		db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		d.conn = db
	})
	return d.conn, err
}
