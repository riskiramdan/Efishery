package main

import (
	"log"

	"github.com/riskiramdan/efishery/golang/databases"
	"github.com/riskiramdan/efishery/golang/seeder"
)

func main() {
	databases.MigrateUp()
	err := seeder.SeedUp()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
}
