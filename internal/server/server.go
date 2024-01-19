package server

import (
	"hm-group-randomizer/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberServer struct {
	*fiber.App
	db *gorm.DB
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(),
		db:  database.InitDB(),
	}

	return server
}
