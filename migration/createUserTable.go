package migration

import (
	"context"
	"fmt"
)

// CreateRoleTable ...
func CreateUserTable() {

	r, err := DbPool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
        id serial PRIMARY KEY,
        name VARCHAR ( 50 ) UNIQUE NOT NULL,
        email VARCHAR ( 255 ) UNIQUE NOT NULL,
				encrypted_password text,
				role_id INTEGER REFERENCES roles (id),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)
`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

}
