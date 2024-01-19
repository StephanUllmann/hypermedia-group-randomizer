package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Static("/static", "./cmd/web/static")

	// Create Batch
	s.App.Get("/", s.GetCreateBatchForm)
	s.App.Get("/get-batch", s.GetBatch)
	s.App.Get("/add-name", func(c *fiber.Ctx) error {
		return c.Redirect("/", http.StatusMovedPermanently)
	})
	s.App.Post("/add-name", s.AddName)
	s.App.Delete("/name", s.RemoveName)
	s.App.Post("/create-batch", s.CreateBatch)

	// Edit Batch
	// s.App.Get("/edit", s.GetFindBatchEditForm)
	// s.App.Get("/edit/form", s.GetEditBatchForm)
	s.App.Put("/edit-batch", s.EditBatch)

	// Delete Batch
	s.App.Delete("/batch/:batchName", s.DeleteBatch)

	// Create Project
	s.App.Get("/project", s.GetCreateProjectForm)
	s.App.Get("/project/q", s.CreateProject)

}
