package main

import (
	"fmt"
	"testing"
	"time"
)

type bench struct {
	Name string
	Desc string
	F    func(b *testing.B)
	N    int
	T    float64
}

func main() {
	fmt.Printf("Last run: %s\n\n", time.Now().Format(time.RFC1123))

	runBenches("Query", benchesQuery(), toMarkdown)
	runBenches("World access", benchesWorld(), toMarkdown)
	runBenches("Entities", benchesEntities(), toMarkdown)
	runBenches("Entities, batched", benchesEntitiesBatch(), toMarkdown)
	runBenches("Components", benchesComponents(), toMarkdown)
	runBenches("Components, batched", benchesComponentsBatch(), toMarkdown)
}
