package seeder

import (
	"database/sql"
	"fmt"

	"github.com/riskiramdan/efishery/golang/config"
)

// SeedUp seeding the database
func SeedUp() error {
	cfg, err := config.GetConfiguration()
	if err != nil {
		return fmt.Errorf("error when getting configuration: %s", err)
	}

	db, err := sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		return fmt.Errorf("error when open postgres connection: %s", err)
	}
	defer db.Close()

	//Roles
	_, err = db.Exec(`
	INSERT INTO public."role"
	(id, "name", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(1, 'Admin', '2021-02-04 17:22:00.028991', 'admin', '2021-02-04 17:22:00.028991', 'admin', NULL, NULL);
	INSERT INTO public."role"
	(id, "name", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(2, 'Operator', '2021-02-04 17:22:21.016457', 'admin', '2021-02-04 17:22:21.016457', 'admin', NULL, NULL);
	INSERT INTO public."role"
	(id, "name", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(3, 'Guest', '2021-02-04 17:22:37.551864', 'admin', '2021-02-04 17:22:37.551864', 'admin', NULL, NULL);	
	`)
	if err != nil {
		return err
	}

	//User
	_, err = db.Exec(`
	INSERT INTO public."user"
	(id, "roleId", "name", phone, "password", "token", "tokenExpiredAt", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(1, 1, 'admin', '081212341234', '$2a$10$TareP8UfLwipzK5jZUQHd.Hu9AA5kFrWxV51qvx4jFnDVPiTb0J.S', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTI1NzA1MzUsImlhdCI6MTYxMjUyNzMzNSwibmFtZSI6ImFkbWluIiwicGhvbmUiOiIwODEyMTIzNDEyMzQiLCJyb2xlSWQiOjEsInRpbWVzdGFtcCI6IjIwMjEtMDItMDZUMDc6MTU6MzUuNzM0NzAyKzA3OjAwIn0.jaTdnoZ26Y-6v5QnXDyMmgqB_jOqVUNJkp1kzRiZ1Dc', '2021-02-06 07:15:35.734702', '2021-02-05 19:15:16.522801', 'admin', '2021-02-05 19:15:35.734701', 'admin', NULL, NULL);	
	`)
	if err != nil {
		return err
	}

	return nil
}
