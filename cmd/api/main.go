package main

import (
	"fmt"
	"hm-group-randomizer/internal/server"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()

	server.Use(compress.New())
	server.Use(recover.New())

	server.Get("/metrics", monitor.New())
	server.Use(logger.New())

	server.RegisterFiberRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
