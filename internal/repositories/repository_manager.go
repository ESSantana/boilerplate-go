package repositories

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/go-sql-driver/mysql"
)

type repositoryManager struct {
	conn *sql.DB
}

func NewRepositoryManager(ctx context.Context) interfaces.RepositoryManager {
	timeLoc, _ := time.LoadLocation("America/Sao_Paulo")
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
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

func (rm *repositoryManager) NewUserRepository() interfaces.UserRepository {
	return newUserRepository(rm.conn)
}
