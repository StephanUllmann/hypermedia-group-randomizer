package server

import (
	"fmt"
	"hm-group-randomizer/cmd/utils"
	web "hm-group-randomizer/cmd/web/templates"
	"hm-group-randomizer/internal/database"
	"log"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// Get Batch - search for existing ans sending back editing, otherwise creation form
func (s *FiberServer) GetBatch(c *fiber.Ctx) error {
	// batch := c.FormValue("batch")
	batch := c.Query("batch")
	loweredBatch := strings.ToLower(batch)
	foundBatch := database.Groups{Batch: loweredBatch}

	err := database.DB.Where("is_base = true AND batch = ?", loweredBatch).First(&foundBatch).Error
	if err != nil {
		fmt.Printf("db: %v\n", err)
		handler := adaptor.HTTPHandler(templ.Handler(web.BatchCollector(loweredBatch, nil, nil)))
		return handler(c)
	}
	fmt.Printf("found batch: %v\n", foundBatch)

	names := strings.Split(foundBatch.Names, ",")

	var transitionNames []templ.Attributes

	for _, name := range names {
		transitionNames = append(transitionNames, utils.CreateTransitionName(name))
	}

	handler := adaptor.HTTPHandler(templ.Handler(web.BatchCollector(loweredBatch, names, transitionNames)))
	return handler(c)

}

// Create Batch
func (s *FiberServer) GetCreateBatchForm(c *fiber.Ctx) error {
	handler := adaptor.HTTPHandler(templ.Handler(web.BatchInput()))
	return handler(c)
}

func (s *FiberServer) AddName(c *fiber.Ctx) error {
	name := c.FormValue("next_name")
	transitionName := utils.CreateTransitionName(name)
	handler := adaptor.HTTPHandler(templ.Handler(web.NewName(name, transitionName)))
	return handler(c)
}

func (s *FiberServer) RemoveName(c *fiber.Ctx) error {
	return c.SendString("")
}

// batch=wd049&Puri=Puri&Miro=Miro&Adam=Adam&Timo=Timo&Hani=Hani&Maxine=Maxine&Test=Mali&Blah=Yosra
func (s *FiberServer) CreateBatch(c *fiber.Ctx) error {
	body := string(c.BodyRaw())
	var batch string
	var names []string
	splitBody := strings.Split(body, "&")

	for _, field := range splitBody {
		parts := strings.Split(field, "=")
		if parts[0] != "batch" && !strings.Contains(parts[1], "%3B") {
			names = append(names, strings.ReplaceAll(parts[1], "%20", " "))
		} else if parts[0] == "batch" {
			batch = parts[1]
		}
	}

	// Check if batch already exists
	var foundBatch string
	rows := database.DB.Table("groups").Where("batch = ?", batch).Select("batch").Row()
	rows.Scan(&foundBatch)
	if foundBatch != "" {
		fmt.Printf("found existing batch on creation: %v\n", foundBatch)
		handler := adaptor.HTTPHandler(templ.Handler(web.Message(fmt.Sprintf("Batch %v already exists", batch), "error")))
		return handler(c)
	}

	namesString := strings.Join(names, ",")

	fmt.Printf("Batch: %v\n", batch)
	fmt.Printf("names: %v\n", names)

	newBatch := database.Groups{Batch: batch, Names: namesString, IsBase: true, Project: "", Group1: "", Group2: "", Group3: "", Group4: "", Group5: "", Group6: "", Group7: ""}

	result := database.DB.Create(&newBatch)

	if result.Error != nil {
		log.Println(result.Error)
		handler := adaptor.HTTPHandler(templ.Handler(web.Message("Unable to create Batch. Please try again", "error")))
		return handler(c)
		// return c.Status(500).SendString("Unable to create Batch. Please try again")
	}

	if result.RowsAffected == 1 {
		handler := adaptor.HTTPHandler(templ.Handler(web.MessageAndButtons(fmt.Sprintf("Batch %v created", batch), "success", batch)))
		return handler(c)
	} else {
		handler := adaptor.HTTPHandler(templ.Handler(web.Message("Unable to create Batch. Please try again", "error")))
		return handler(c)
	}

}

// Edit Batch
// func (s *FiberServer) GetFindBatchEditForm(c *fiber.Ctx) error  {
// 	handler := adaptor.HTTPHandler(templ.Handler(web.FindBatchEdit()))
// 	return handler(c)
// }

// func (s *FiberServer) GetEditBatchForm(c *fiber.Ctx) error  {
// 	batch := c.FormValue("batch")
// 	loweredBatch := strings.ToLower(batch)
// 	foundBatch := database.Groups{Batch: loweredBatch}

// 	err := database.DB.Where("is_base = true AND batch = ?", loweredBatch).First(&foundBatch).Error
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		handler := adaptor.HTTPHandler(templ.Handler(web.Message(fmt.Sprintf("Batch %v not found", batch), "error")))
// 		return handler(c)
// 	}
// 	fmt.Printf("found batch: %v\n", foundBatch)

// 	names := strings.Split(foundBatch.Names, ",")

// 	var transitionNames []templ.Attributes

// 	for _, name := range names {
// 		transitionNames = append(transitionNames, utils.CreateTransitionName(name))
// 	}

// 	handler := adaptor.HTTPHandler(templ.Handler(web.BatchEdit(names, transitionNames)))
// 	return handler(c)
// }

func (s *FiberServer) EditBatch(c *fiber.Ctx) error {
	body := string(c.BodyRaw())
	var batch string
	var names []string
	splitBody := strings.Split(body, "&")

	for _, field := range splitBody {
		parts := strings.Split(field, "=")
		if parts[0] != "batch" && !strings.Contains(parts[1], "%3B") {
			names = append(names, strings.ReplaceAll(parts[1], "%20", " "))
		} else if parts[0] == "batch" {
			batch = parts[1]
		}
	}

	// Check if batch actually exists
	var foundBatch database.Groups
	rowsAffected := database.DB.Where("batch = ? AND is_base = ?", batch, true).Find(&foundBatch).RowsAffected
	if rowsAffected != 1 {
		fmt.Printf("To edit Batch doesn't exist: %v\n", foundBatch)
		handler := adaptor.HTTPHandler(templ.Handler(web.Message(fmt.Sprintf("Batch %v doesn't exists", batch), "error")))
		return handler(c)
	}

	namesString := strings.Join(names, ",")
	fmt.Printf("found batch: %v\n", foundBatch)
	result := database.DB.Model(&foundBatch).Where("is_base = ?", true).Update("names", namesString)

	if result.RowsAffected == 1 {
		handler := adaptor.HTTPHandler(templ.Handler(web.MessageAndButtons(fmt.Sprintf("Batch %v updated", batch), "success", batch)))
		return handler(c)
	} else {
		fmt.Printf("error: %v\n", result)
		handler := adaptor.HTTPHandler(templ.Handler(web.MessageAndButtons("Something went wrong :(", "error", batch)))
		return handler(c)
	}
}

// Delete Batch
func (s *FiberServer) DeleteBatch(c *fiber.Ctx) error {
	batch := c.Params("batchName")

	var deletedBatch []database.Groups

	err := database.DB.Unscoped().Where("batch = ?", batch).Delete(&deletedBatch).Error
	if err != nil {
		fmt.Printf("error: %v\n", err)
		handler := adaptor.HTTPHandler(templ.Handler(web.Message(fmt.Sprintf("Batch %v not found", batch), "error")))
		return handler(c)
	}

	fmt.Printf("deleted batch: %v\n", deletedBatch)

	// return c.Redirect("/", http.StatusMovedPermanently)
	handler := adaptor.HTTPHandler(templ.Handler(web.MessageAndManageBatches(fmt.Sprintf("Batch %v deleted", batch), "success")))
	return handler(c)
}
