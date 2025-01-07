package main

import (
	"fmt"
	"time"

	"github.com/mlange-42/arche/benchmark"
)

func main() {
	fmt.Printf("Last run: %s\n\n", time.Now().Format(time.RFC1123))

	benchmark.RunBenchmarks("Query", benchesQuery(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("World access", benchesWorld(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities", benchesEntities(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities, batched", benchesEntitiesBatch(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components", benchesComponents(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components, batched", benchesComponentsBatch(), benchmark.ToMarkdown)
}
