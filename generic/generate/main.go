package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var typeLetters = []string{"A", "B", "C", "D", "E", "F", "G", "H"}

type query struct {
	Index       int
	Types       string
	TypesFull   string
	TypesReturn string
	ReturnAll   string
	Include     string
}
type getter struct {
	Query query
	ID    int
	Index int
	Type  string
}

func main() {
	fmt.Println("Generating code")

	maxIndex := len(typeLetters)

	text := bytes.Buffer{}

	header, err := template.ParseFiles("./generate/header.gotxt")
	if err != nil {
		panic(err)
	}
	err = header.Execute(&text, struct{}{})
	if err != nil {
		panic(err)
	}

	filters, err := template.ParseFiles("./generate/filter.gotxt")
	if err != nil {
		panic(err)
	}
	queryGetAll, err := template.ParseFiles("./generate/query_getall.gotxt")
	if err != nil {
		panic(err)
	}
	queryGet, err := template.ParseFiles("./generate/query_get.gotxt")
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
			returnTypes = "(*" + strings.Join(typeLetters[:i], ", *") + ")"
			fullTypes = "[" + strings.Join(typeLetters[:i], " any, ") + " any]"
			include = "include: []reflect.Type{typeOf[" + strings.Join(typeLetters[:i], "](), typeOf[") + "]()},"
			for j := 0; j < i; j++ {
				returnAll += fmt.Sprintf("(*%s)(q.Query.Get(q.ids[%d]))", typeLetters[j], j)
				if j < i-1 {
					returnAll += ", "
				}
			}
		}
		data := query{
			Index:       i,
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
				Query: data,
				ID:    j + 1,
				Index: j,
				Type:  typeLetters[j],
			}
			err = queryGet.Execute(&text, getterData)
			if err != nil {
				panic(err)
			}
		}
	}
	if err := os.WriteFile("filter.go", text.Bytes(), 0666); err != nil {
		panic(err)
	}

}
