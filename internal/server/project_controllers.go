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

	var groupArr [][]string
	for _, group := range foundGroups {
		groupArr = append(groupArr, strings.Split(group.Names, ","))
	}
	newGroup := utils.ShuffleGroups(groupArr)
	groups := utils.SortToGroups(newGroup, numOfGroups)
	fmt.Println(groups)
	fmt.Printf("numOfGroups: %v\n", numOfGroups)
	fmt.Println("groupArr: ", groupArr)
	fmt.Println("project: ", project)

	newGroupString := strings.Join(newGroup, ",")
	group1 := strings.Join(groups[0], ",")
	group2 := strings.Join(groups[1], ",")
	group3 := strings.Join(groups[2], ",")
	group4 := strings.Join(groups[3], ",")
	group5 := strings.Join(groups[4], ",")
	group6 := strings.Join(groups[5], ",")
	group7 := strings.Join(groups[6], ",")

	newEntry := database.Groups{Batch: batch, Names: newGroupString, Group1: group1, Group2: group2, Group3: group3, Group4: group4, Group5: group5, Group6: group6, Group7: group7, Project: project, IsBase: false}

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
