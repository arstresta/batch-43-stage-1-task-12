package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnection() {
	var err error
	databaseURL := "postgres://postgres:police321@localhost:5432/personal_web"
	// databaseURL := "postgres://{user}:{password}@{serverName:port}}/{namaDatabase}"
	Conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v", err)
		os.Exit(1)
	}
	fmt.Println("Database Connected")
}
