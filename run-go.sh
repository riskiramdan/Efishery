cd golang/databases

rice embed-go 

cd ..

go run cmd/efishery-migrate/main.go
go run cmd/efishery-seeder/main.go
go run cmd/efishery/main.go