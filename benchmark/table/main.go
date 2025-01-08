package main

import (
	"fmt"
	"time"

	"github.com/mlange-42/arche/benchmark"
	"github.com/shirou/gopsutil/v4/cpu"
)

const version = "v0.14.6-dev"

func main() {
	fmt.Printf("Last run: %s  \n", time.Now().Format(time.RFC1123))
	fmt.Printf("Version: Arche %s  \n", version)

	infos, err := cpu.Info()
	if err != nil {
		panic(err)
	}
	for _, info := range infos {
		fmt.Printf("CPU: %s\n\n", info.ModelName)
		break
	}

	benchmark.RunBenchmarks("Query", benchesQuery(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("World access", benchesWorld(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities", benchesEntities(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities, batched", benchesEntitiesBatch(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components", benchesComponents(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components, batched", benchesComponentsBatch(), benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Other", benchesOther(), benchmark.ToMarkdown)
}
