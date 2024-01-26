package utils

import (
	"fmt"
	"hm-group-randomizer/internal/database"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/a-h/templ"
)

func CreateTransitionName(name string) templ.Attributes {
	name = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(name, "")
	transitionName := fmt.Sprintf("view-transition-name:%s;", name)
	transitionStyle := templ.Attributes{"style": transitionName}
	return transitionStyle
}

func shuffle(toShuffle []string) []string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := range toShuffle {
		randInd := r.Intn(len(toShuffle))
		toShuffle[i], toShuffle[randInd] = toShuffle[randInd], toShuffle[i]
	}
	return toShuffle
}

func ShuffleGroups(prevGroups [][]string) []string {
	shuffled := make([]string, len(prevGroups[0]))
	copy(shuffled, prevGroups[0])

	isSimilar := true
	count := 0
	maxIteration := 1000

	for isSimilar && count < maxIteration {
		shuffled = shuffle(shuffled)
		for _, arr := range prevGroups {
			if reflect.DeepEqual(arr, shuffled) {
				count++
				if count == 1000 {
					fmt.Println("Reached max shuffle iteration")
				}
				break
			}
			isSimilar = false
		}
	}

	return shuffled
}

func SortToGroups(all []string, numGroups int) [][]string {
	lenAll := len(all)
	if numGroups == 0 {
		if lenAll < 10 {
			numGroups = 3
		} else if lenAll > 9 && lenAll < 15 {
			numGroups = 4
		} else if lenAll > 14 && lenAll < 21 {
			numGroups = 5
		} else {
			numGroups = 6
		}
	}

	out := make([][]string, numGroups)

	for i, name := range all {
		out[i%numGroups] = append(out[i%numGroups], name)
	}

	for len(out) < 7 {
		out = append(out, []string{""})
	}

	return out

}


func GroupsStringsToDisplay(foundGroup database.Groups) [][]string {
	var displayGroups [][]string
	groups := []string{foundGroup.Group1, foundGroup.Group2, foundGroup.Group3, foundGroup.Group4, foundGroup.Group5, foundGroup.Group6, foundGroup.Group7}

	for _, group := range groups {
		if group != "" {
			displayGroups = append(displayGroups, strings.Split(group, ","))
		}
	}
	return displayGroups
}

// func AddAnimationDelay(group [][]string) [][]templ.Attributes {
// 	var out [][]templ.Attributes
// 	for i, arr := range group {
// 		var outArr []templ.Attributes
// 		for j := range arr {
// 			delay := (i + 1) * 350 + (j + 1) * 1700 - (1800 / (i + 1) * (j + 1))
// 			// delay := (i * len(group) + j) * 2000 

// 			styleString := fmt.Sprintf("animation-delay:%dms;", delay)
// 			style := templ.Attributes{"style": styleString}
// 			outArr = append(outArr, style)
// 		}
// 		out = append(out, outArr)
// 	}
// 	return out
// }

func AddAnimationDelay(group [][]string) [][]templ.Attributes {
	var out [][]templ.Attributes
	totalElements := 0
	maxLength := 0
	for _, arr := range group {
			totalElements += len(arr)
			if len(arr) > maxLength {
					maxLength = len(arr)
			}
			out = append(out, make([]templ.Attributes, len(arr)))
	}
	for i := 0; i < maxLength; i++ {
			for j := 0; j < len(group); j++ {
					if i < len(group[j]) {
							delay := (i * len(group) + j) *  625
							styleString := fmt.Sprintf("animation-delay:%dms;", delay)
							style := templ.Attributes{"style": styleString}
							out[j][i] = style
					}
			}
	}
	return out
}

func GroupsToEntry(foundGroups []database.Groups, numOfGroups int, batchName string, projectName string) database.Groups {
	var groupArr [][]string
	for _, group := range foundGroups {
		groupArr = append(groupArr, strings.Split(group.Names, ","))
	}
	newGroup := ShuffleGroups(groupArr)
	groups := SortToGroups(newGroup, numOfGroups)
	

	newGroupString := strings.Join(newGroup, ",")
	group1 := strings.Join(groups[0], ",")
	group2 := strings.Join(groups[1], ",")
	group3 := strings.Join(groups[2], ",")
	group4 := strings.Join(groups[3], ",")
	group5 := strings.Join(groups[4], ",")
	group6 := strings.Join(groups[5], ",")
	group7 := strings.Join(groups[6], ",")

	entry := database.Groups{Batch: batchName, Names: newGroupString, Group1: group1, Group2: group2, Group3: group3, Group4: group4, Group5: group5, Group6: group6, Group7: group7, Project: projectName, IsBase: false}

	return entry
}