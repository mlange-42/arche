// Package generate is for generating the boilerplate code required for the generic API.
package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var typeLetters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
var variableLetters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
var numbers = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
var numberStr = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve"}

type query struct {
	Index         int
	NumberStr     string
	Types         string
	TypesFull     string
	TypesReturn   string
	Variables     string
	ReturnAll     string
	ReturnAllSafe string
	Include       string
	Assign        string
	Arguments     string
	IDTypes       string
	IDAssign      string
	IDAssign2     string
	IDList        string
}

func main() {
	generateQueries()
	generateMaps()
}

func generateMaps() {
	fmt.Println("Generating maps")

	maxIndex := len(typeLetters)
	text := bytes.Buffer{}

	header, err := template.ParseFiles("./_generate/map_header.go.txt")
	if err != nil {
		panic(err)
	}
	err = header.Execute(&text, struct{}{})
	if err != nil {
		panic(err)
	}

	maps, err := template.ParseFiles("./_generate/map.go.txt")
	if err != nil {
		panic(err)
	}

	for i := 1; i <= maxIndex; i++ {
		types := ""
		returnTypes := ""
		fullTypes := ""
		returnAll := ""
		returnAllSafe := ""
		assign := ""
		arguments := ""
		idTypes := ""
		idAssign := ""
		idAssign2 := ""
		idList := ""
		variables := ""

		if i > 0 {
			types = "[" + strings.Join(typeLetters[:i], ", ") + "]"
			variables = strings.Join(variableLetters[:i], ", ")
			returnTypes = "*" + strings.Join(typeLetters[:i], ", *")
			fullTypes = "[" + strings.Join(typeLetters[:i], " any, ") + " any]"
			idTypes = "id" + strings.Join(numbers[:i], " ecs.ID\n\tid") + " ecs.ID"
			idList = "m.id" + strings.Join(numbers[:i], ", m.id")
		}

		for j := 0; j < i; j++ {
			returnAll += fmt.Sprintf("(*%s)(m.world.GetUnchecked(entity, m.id%d))", typeLetters[j], j)
			if j == 0 {
				returnAllSafe += fmt.Sprintf("(*%s)(m.world.Get(entity, m.id%d))", typeLetters[j], j)
			} else {
				returnAllSafe += fmt.Sprintf("(*%s)(m.world.GetUnchecked(entity, m.id%d))", typeLetters[j], j)
			}
			arguments += fmt.Sprintf("%s *%s", strings.ToLower(typeLetters[j]), typeLetters[j])
			idAssign += fmt.Sprintf("	id%d: ecs.ComponentID[%s](w),\n", j, typeLetters[j])
			idAssign2 += fmt.Sprintf("	id%d: m.id%d,\n", j, j)
			if j < i-1 {
				returnAll += ",\n"
				returnAllSafe += ",\n"
				arguments += ", "
			}
			if j == 0 {
				assign += fmt.Sprintf("*(*%s)(m.world.Get(entity, m.id%d)) = *%s\n", typeLetters[j], j, strings.ToLower(typeLetters[j]))
			} else {
				assign += fmt.Sprintf("*(*%s)(m.world.GetUnchecked(entity, m.id%d)) = *%s\n", typeLetters[j], j, strings.ToLower(typeLetters[j]))
			}
		}

		data := query{
			Index:         i,
			NumberStr:     numberStr[i],
			Types:         types,
			TypesReturn:   returnTypes,
			TypesFull:     fullTypes,
			ReturnAll:     returnAll,
			ReturnAllSafe: returnAllSafe,
			Variables:     variables,
			Assign:        assign[:len(assign)-1],
			Arguments:     arguments,
			IDTypes:       idTypes,
			IDAssign:      idAssign,
			IDAssign2:     idAssign2,
			IDList:        idList,
		}
		err = maps.Execute(&text, data)
		if err != nil {
			panic(err)
		}
	}

	if err := os.WriteFile("map_generated.go", text.Bytes(), 0666); err != nil {
		panic(err)
	}
}

func generateQueries() {
	fmt.Println("Generating queries")

	maxIndex := len(typeLetters)

	text := bytes.Buffer{}

	header, err := template.ParseFiles("./_generate/query_header.go.txt")
	if err != nil {
		panic(err)
	}
	err = header.Execute(&text, struct{}{})
	if err != nil {
		panic(err)
	}

	filters, err := template.ParseFiles("./_generate/query.go.txt")
	if err != nil {
		panic(err)
	}

	for i := 0; i <= maxIndex; i++ {
		types := ""
		returnTypes := ""
		fullTypes := ""
		include := ""
		returnAll := ""
		idTypes := ""
		idAssign := ""
		variables := ""
		if i > 0 {
			types = "[" + strings.Join(typeLetters[:i], ", ") + "]"
			variables = strings.Join(variableLetters[:i], ", ")
			returnTypes = "*" + strings.Join(typeLetters[:i], ", *")
			fullTypes = "[" + strings.Join(typeLetters[:i], " any, ") + " any]"
			include = "typeOf[" + strings.Join(typeLetters[:i], "](),\ntypeOf[") + "](),"
			idTypes = "id" + strings.Join(numbers[:i], " ecs.ID\n\tid") + " ecs.ID"
			for j := 0; j < i; j++ {
				returnAll += fmt.Sprintf("(*%s)(q.Query.Get(q.id%d))", typeLetters[j], j)
				idAssign += fmt.Sprintf("	id%d: f.compiled.Ids[%d],\n", j, j)
				if j < i-1 {
					returnAll += ",\n"
				}
			}
		} else {
			include = ""
		}
		data := query{
			Index:       i,
			NumberStr:   numberStr[i],
			Types:       types,
			TypesReturn: returnTypes,
			TypesFull:   fullTypes,
			Variables:   variables,
			ReturnAll:   returnAll,
			Include:     include,
			IDTypes:     idTypes,
			IDAssign:    idAssign,
		}
		err = filters.Execute(&text, data)
		if err != nil {
			panic(err)
		}
		if i == 0 {
			continue
		}
	}
	if err := os.WriteFile("query_generated.go", text.Bytes(), 0666); err != nil {
		panic(err)
	}
}
