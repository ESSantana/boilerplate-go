package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/internal/repositories/interfaces"
	"github.com/go-sql-driver/mysql"
)

type repositoryManager struct {
	conn *sql.DB
}

type Scanner interface {
	Scan(dest ...any) error
	Err() error
}

func NewRepositoryManager(ctx context.Context, cfg *config.Config) interfaces.RepositoryManager {
	timeLoc, _ := time.LoadLocation("America/Sao_Paulo")

	dbCfg := mysql.Config{
		User:                 cfg.Database.User,
		Passwd:               cfg.Database.Password,
		Net:                  "tcp",
		Addr:                 cfg.Database.Host + ":" + cfg.Database.Port,
		DBName:               cfg.Database.Name,
		Loc:                  timeLoc,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	mysqlConn, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	mysqlConn.SetConnMaxLifetime(time.Minute * 3)
	mysqlConn.SetMaxOpenConns(10)
	mysqlConn.SetMaxIdleConns(10)

	return &repositoryManager{
		conn: mysqlConn,
	}
}

func (rm *repositoryManager) DatabaseHealthCheck() error {
	return rm.conn.Ping()
}

func (rm *repositoryManager) NewCustomerRepository() interfaces.CustomerRepository {
	return newCustomerRepository(rm.conn)
}
