package migration

import (
	"context"
	"fmt"
)

// CreateRoleTable ...
func CreateRoleTable() {

	r, err := DbPool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS roles (
        id serial PRIMARY KEY,
        name VARCHAR ( 50 ) UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)
`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

}
