package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var typeLetters = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
var numbers = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight"}
var indices = []string{"zeroth", "first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth"}

type query struct {
	Index       int
	NumberStr   string
	Types       string
	TypesFull   string
	TypesReturn string
	ReturnAll   string
	Include     string
	Components  string
	Arguments   string
}
type getter struct {
	Query    query
	ID       int
	Index    int
	IndexStr string
	Type     string
}

func main() {
	generateQueries()
	generateMutates()
}

func generateMutates() {
	fmt.Println("Generating mutates")

	maxIndex := len(typeLetters)
	text := bytes.Buffer{}

	header, err := template.ParseFiles("./generate/mutate_header.go.txt")
	if err != nil {
		panic(err)
	}
	err = header.Execute(&text, struct{}{})
	if err != nil {
		panic(err)
	}

	mutate, err := template.ParseFiles("./generate/mutate.go.txt")
	if err != nil {
		panic(err)
	}

	for i := 0; i <= maxIndex; i++ {
		types := ""
		returnTypes := ""
		fullTypes := ""
		returnAll := ""
		include := ""
		components := ""
		arguments := ""

		if i > 0 {
			types = "[" + strings.Join(typeLetters[:i], ", ") + "]"
			returnTypes = "*" + strings.Join(typeLetters[:i], ", *")
			fullTypes = "[" + strings.Join(typeLetters[:i], " any, ") + " any]"
			include = "[]ecs.ID{ecs.ComponentID[" + strings.Join(typeLetters[:i], "](w), ecs.ComponentID[") + "](w)}"
		} else {
			include = "[]ecs.ID{}"
		}

		for j := 0; j < i; j++ {
			returnAll += fmt.Sprintf("(*%s)(m.world.Get(entity, m.ids[%d]))", typeLetters[j], j)
			arguments += fmt.Sprintf("%s *%s", strings.ToLower(typeLetters[j]), typeLetters[j])
			if j < i-1 {
				returnAll += ", "
				arguments += ", "
			}
			components += fmt.Sprintf("ecs.Component{ID: m.ids[%d], Component: %s},\n", j, strings.ToLower(typeLetters[j]))
		}

		data := query{
			Index:       i,
			NumberStr:   numbers[i],
			Types:       types,
			TypesReturn: returnTypes,
			TypesFull:   fullTypes,
			ReturnAll:   returnAll,
			Include:     include,
			Components:  components,
			Arguments:   arguments,
		}
		err = mutate.Execute(&text, data)
		if err != nil {
			panic(err)
		}
	}

	if err := os.WriteFile("mutate_generated.go", text.Bytes(), 0666); err != nil {
		panic(err)
	}
}

func generateQueries() {
	fmt.Println("Generating queries")

	maxIndex := len(typeLetters)

	text := bytes.Buffer{}

	header, err := template.ParseFiles("./generate/query_header.go.txt")
	if err != nil {
		panic(err)
	}
	err = header.Execute(&text, struct{}{})
	if err != nil {
		panic(err)
	}

	filters, err := template.ParseFiles("./generate/query.go.txt")
	if err != nil {
		panic(err)
	}
	queryGetAll, err := template.ParseFiles("./generate/query_getall.go.txt")
	if err != nil {
		panic(err)
	}
	queryGet, err := template.ParseFiles("./generate/query_get.go.txt")
	if err != nil {
		panic(err)
	}

	for i := 0; i <= maxIndex; i++ {
		types := ""
		returnTypes := ""
		fullTypes := ""
		include := ""
		returnAll := ""
		if i > 0 {
			types = "[" + strings.Join(typeLetters[:i], ", ") + "]"
			returnTypes = "*" + strings.Join(typeLetters[:i], ", *")
			fullTypes = "[" + strings.Join(typeLetters[:i], " any, ") + " any]"
			include = "[]Comp{typeOf[" + strings.Join(typeLetters[:i], "](), typeOf[") + "]()}"
			for j := 0; j < i; j++ {
				returnAll += fmt.Sprintf("(*%s)(q.Query.Get(q.ids[%d]))", typeLetters[j], j)
				if j < i-1 {
					returnAll += ", "
				}
			}
		} else {
			include = "[]Comp{}"
		}
		data := query{
			Index:       i,
			NumberStr:   numbers[i],
			Types:       types,
			TypesReturn: returnTypes,
			TypesFull:   fullTypes,
			ReturnAll:   returnAll,
			Include:     include,
		}
		err = filters.Execute(&text, data)
		if err != nil {
			panic(err)
		}
		if i == 0 {
			continue
		}

		err = queryGetAll.Execute(&text, data)
		if err != nil {
			panic(err)
		}

		for j := 0; j < i; j++ {
			getterData := getter{
				Query:    data,
				ID:       j + 1,
				Index:    j,
				IndexStr: indices[j+1],
				Type:     typeLetters[j],
			}
			err = queryGet.Execute(&text, getterData)
			if err != nil {
				panic(err)
			}
		}
	}
	if err := os.WriteFile("query_generated.go", text.Bytes(), 0666); err != nil {
		panic(err)
	}
}
