package api

import (
	"bookstore/config"
	"bookstore/database"
	"bookstore/model"
	"context"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestMain(m *testing.M) {
	ConnectTestDB()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterRoutes(r)
	Router = r

	// run tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

func ConnectTestDB() {
	ctx := context.Background()
	// set database config
	config.Database.Host = "localhost"
	config.Database.Username = "root"
	config.Database.Password = "password"
	config.Database.Database = "bookstore"

	// start mysql testcontainer
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:9.0",
		mysql.WithDatabase(config.Database.Database),
		mysql.WithUsername(config.Database.Username),
		mysql.WithPassword(config.Database.Password),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// get mapped port
	p, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("failed to get mapped port: %s", err)
	}
	config.Database.Port = p.Port()

	// initialize database
	database.InitializeDB()
}

func ResetBookTable() {
	database.DB.Where("1 = 1").Delete(&model.Book{})
}
