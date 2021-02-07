package main

import (
	"fmt"
	"log"
	testUser "os/user"

	"github.com/go-redis/redis/v8"
	"github.com/riskiramdan/efishery/golang/config"
	"github.com/riskiramdan/efishery/golang/databases"
	"github.com/riskiramdan/efishery/golang/internal/concurrency"
	"github.com/riskiramdan/efishery/golang/internal/data"
	"github.com/riskiramdan/efishery/golang/internal/hosts"
	internalhttp "github.com/riskiramdan/efishery/golang/internal/http"
	"github.com/riskiramdan/efishery/golang/internal/user"
	userPg "github.com/riskiramdan/efishery/golang/internal/user/postgres"
	redisManager "github.com/riskiramdan/efishery/golang/redis"
	"github.com/riskiramdan/efishery/golang/util"

	"github.com/jmoiron/sqlx"
)

// InternalServices represents all the internal domain services
type InternalServices struct {
	userService        user.ServiceInterface
	concurrencyService concurrency.ServiceInterface
}

func buildInternalServices(db *sqlx.DB, config *config.Config, redisManager *redis.Client) *InternalServices {
	userPostgresStorage := userPg.NewPostgresStorage(
		data.NewPostgresStorage(db, "users", user.Users{}),
	)
	userService := user.NewService(userPostgresStorage)

	concurrencyService := concurrency.NewService(&hosts.HTTPManager{}, redisManager)
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
	redis, err := redisManager.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to connect to redis x: ", err)
	}

	util := &util.Utility{}
	httpManager := &hosts.HTTPManager{}
	defer db.Close()
	dataManager := data.NewManager(db)
	internalServices := buildInternalServices(db, config, redis)
	// Migrate the db
	databases.MigrateUp()

	s := internalhttp.NewServer(
		internalServices.userService,
		internalServices.concurrencyService,
		dataManager,
		config,
		util,
		httpManager,
		redis,
	)
	s.Serve()
}
