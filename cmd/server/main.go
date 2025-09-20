package main

import (
	"context"

	"github.com/joaopdias/blog-server/internal/api"
	"github.com/joaopdias/blog-server/internal/config"
	"github.com/joaopdias/blog-server/internal/database"
)

func main() {
	
	
	config := config.Load()
	pool, err := database.NewPool(context.Background(), config.DSN)
	
	if err != nil {
		panic("failed to connect with database: " + err.Error())
	}
	
	api := api.NewRouter(pool)

	api.Run(":"+config.Port)
}
