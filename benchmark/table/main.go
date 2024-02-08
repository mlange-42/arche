package main

import (
	"testing"
)

type bench struct {
	Name string
	Desc string
	F    func(b *testing.B)
	N    int
	T    float64
}

func main() {
	runBenches("Query", benchesQuery(), toMarkdown)
	runBenches("Entities", benchesEntities(), toMarkdown)
	runBenches("Entities, batched", benchesEntitiesBatch(), toMarkdown)
}
