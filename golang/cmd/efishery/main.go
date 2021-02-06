package main

import (
	"fmt"
	"log"
	testUser "os/user"
	"time"

	"github.com/riskiramdan/efishery/golang/config"
	"github.com/riskiramdan/efishery/golang/databases"
	"github.com/riskiramdan/efishery/golang/internal/concurrency"
	"github.com/riskiramdan/efishery/golang/internal/data"
	"github.com/riskiramdan/efishery/golang/internal/hosts"
	internalhttp "github.com/riskiramdan/efishery/golang/internal/http"
	"github.com/riskiramdan/efishery/golang/internal/user"
	userPg "github.com/riskiramdan/efishery/golang/internal/user/postgres"
	"github.com/riskiramdan/efishery/golang/util"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
)

// InternalServices represents all the internal domain services
type InternalServices struct {
	userService        user.ServiceInterface
	concurrencyService concurrency.ServiceInterface
}

func buildInternalServices(db *sqlx.DB, config *config.Config) *InternalServices {
	userPostgresStorage := userPg.NewPostgresStorage(
		data.NewPostgresStorage(db, "users", user.Users{}),
	)
	userService := user.NewService(userPostgresStorage)
	cache := cache.New(time.Hour*12, time.Hour*12)
	concurrencyService := concurrency.NewService(&hosts.HTTPManager{}, cache)
	return &InternalServices{
		userService:        userService,
		concurrencyService: concurrencyService,
	}
}

func main() {

	usr, err := testUser.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.Username)

	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}
	db, err := sqlx.Open("postgres", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	util := &util.Utility{}
	httpManager := &hosts.HTTPManager{}
	defer db.Close()
	dataManager := data.NewManager(db)
	internalServices := buildInternalServices(db, config)
	// Migrate the db
	databases.MigrateUp()

	s := internalhttp.NewServer(
		internalServices.userService,
		internalServices.concurrencyService,
		dataManager,
		config,
		util,
		httpManager,
	)
	s.Serve()
}
