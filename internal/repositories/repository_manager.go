package repositories

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/ESSantana/boilerplate-backend/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/utils"
	"github.com/go-sql-driver/mysql"
)

type repositoryManager struct {
	conn *sql.DB
}

type Scanner interface {
	Scan(dest ...any) error
	Err() error
}

func NewRepositoryManager(ctx context.Context) interfaces.RepositoryManager {
	timeLoc, _ := time.LoadLocation("America/Sao_Paulo")

	mysqlUser := utils.RetrieveSecretValue("DB_USER_FILE")
	mysqlPass := utils.RetrieveSecretValue("DB_PASS_FILE")
	if os.Getenv("ENV") == "development" {
		mysqlUser = os.Getenv("DB_USER")
		mysqlPass = os.Getenv("DB_PASS")
	}

	cfg := mysql.Config{
		User:                 mysqlUser,
		Passwd:               mysqlPass,
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOST"),
		DBName:               os.Getenv("DB_NAME"),
		Loc:                  timeLoc,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	mysqlConn, err := sql.Open("mysql", cfg.FormatDSN())
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
