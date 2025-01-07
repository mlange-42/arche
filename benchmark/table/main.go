package main

import (
	"fmt"
	"time"

	"github.com/mlange-42/arche/benchmark"
	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {
	infos, err := cpu.Info()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Last run: %s\n", time.Now().Format(time.RFC1123))
	for _, info := range infos {
		fmt.Printf("CPU: %s\n\n", info.ModelName)
	}

	benchmark.RunBenchmarks("Query", benchesQuery(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("World access", benchesWorld(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities", benchesEntities(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities, batched", benchesEntitiesBatch(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components", benchesComponents(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components, batched", benchesComponentsBatch(), benchmark.ToMarkdown)
}
