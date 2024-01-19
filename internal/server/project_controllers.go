package server

import (
	"fmt"
	"hm-group-randomizer/cmd/utils"
	web "hm-group-randomizer/cmd/web/templates"
	"hm-group-randomizer/internal/database"
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// Create Project
func (s *FiberServer) GetCreateProjectForm(c *fiber.Ctx) error {
	handler := adaptor.HTTPHandler(templ.Handler(web.CreateProject()))
	return handler(c)
}

func (s *FiberServer) CreateProject(c *fiber.Ctx) error {
	batch := c.FormValue("batch")
	project := c.FormValue("project")
	numOfGroups, err := strconv.Atoi(c.FormValue("number"))
	if err != nil {
		fmt.Println("Error: ", err)
		numOfGroups = 0
		// handler := adaptor.HTTPHandler(templ.Handler(web.Message("Creating Project failed", "error")))
		// return handler(c)
	}

	var foundGroups []database.Groups

	err = database.DB.Table("groups").Where("batch = ?", batch).Find(&foundGroups).Error
	if err != nil {
		fmt.Println("Error: ", err)
		handler := adaptor.HTTPHandler(templ.Handler(web.Message("Creating Project failed", "error")))
		return handler(c)
	}
	// Batch not found <- empty slice
	if len(foundGroups) == 0 {
		handler := adaptor.HTTPHandler(templ.Handler(web.Message(fmt.Sprintf("Batch %v doesn't exist  ¯\\_(ツ)_/¯", batch), "error")))
		return handler(c)
	}

	// Check if project already exists	
	projectInd := slices.IndexFunc(foundGroups, func(group database.Groups) bool {
		return strings.EqualFold(group.Project, project)
	})

	// Respond with already sorted groups
	if projectInd != -1 {
		group := foundGroups[projectInd]
		displayGroups := utils.GroupsStringsToDisplay(group)
		// delays := utils.AddAnimationDelay(displayGroups)
		handler := adaptor.HTTPHandler(templ.Handler(web.ExistingProject(group.Batch, group.Project, displayGroups, nil)))
		return handler(c)
	}

	newEntry := utils.GroupsToEntry(foundGroups, numOfGroups, batch, project)

	result := database.DB.Create(&newEntry)
	if result.Error != nil {
		fmt.Println(result.Error)
		handler := adaptor.HTTPHandler(templ.Handler(web.Message("Unable to create Project. Please try again", "error")))
		return handler(c)
	}

	displayGroups := utils.GroupsStringsToDisplay(newEntry)
	delays := utils.AddAnimationDelay(displayGroups)
	handler := adaptor.HTTPHandler(templ.Handler(web.ExistingProject(newEntry.Batch, newEntry.Project, displayGroups, delays)))
	return handler(c)
}

// Shuffle Project
func (s *FiberServer) ShuffleProject(c *fiber.Ctx) error {
	batchName := c.Params("batchName")
	projectName := c.Params("projectName")

	var foundGroups []database.Groups

	err := database.DB.Table("groups").Where("batch = ?", batchName).Find(&foundGroups).Error
	if err != nil {
		fmt.Println("Error: ", err)
		handler := adaptor.HTTPHandler(templ.Handler(web.Message("Creating Project failed", "error")))
		return handler(c)
	}

	idx := slices.IndexFunc(foundGroups, func (group database.Groups) bool {
		return strings.EqualFold(group.Project, projectName)
	})

	var numOfGroups int
	groups := []string{foundGroups[idx].Group1, foundGroups[idx].Group2, foundGroups[idx].Group3, foundGroups[idx].Group4, foundGroups[idx].Group5, foundGroups[idx].Group6, foundGroups[idx].Group7}

	for _, group := range groups {
		if group != "" {
			numOfGroups++
		}
	}

	patchedEntry := utils.GroupsToEntry(foundGroups, numOfGroups, batchName, projectName)
	patchedEntry.ID = foundGroups[idx].ID

	result := database.DB.Save(&patchedEntry)
	fmt.Println(result)


	displayGroups := utils.GroupsStringsToDisplay(patchedEntry)
	delays := utils.AddAnimationDelay(displayGroups)
	handler := adaptor.HTTPHandler(templ.Handler(web.ExistingProject(patchedEntry.Batch, patchedEntry.Project, displayGroups, delays)))
	return handler(c)
}