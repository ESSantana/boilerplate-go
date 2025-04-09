package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/application-ellas/ella-backend/internal/repositories/interfaces"
	"github.com/go-sql-driver/mysql"
)

type repositoryManager struct {
	conn *sql.DB
}

type Scanner interface{ Scan(dest ...any) error }

func NewRepositoryManager(ctx context.Context) interfaces.RepositoryManager {
	fmt.Println("Connecting to MySQL database...")
	data, err := os.ReadFile(os.Getenv("DB_PASS"))
	if err != nil {
		return &repositoryManager{}
	}
	fmt.Println("DB_PASS:", string(data))

	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))

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

func (rm *repositoryManager) NewCustomerRepository() interfaces.CustomerRepository {
	return newCustomerRepository(rm.conn)
}
