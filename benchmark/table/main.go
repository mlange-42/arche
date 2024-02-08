package main

import (
	"fmt"
	"strings"
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
	query := []bench{
		{Name: "Query.Next", Desc: "", F: queryIter_100_000, N: 100_000},
		{Name: "Query.Next + 1x Get", Desc: "", F: queryIterGet_1_100_000, N: 100_000},
		{Name: "Query.Next + 2x Get", Desc: "", F: queryIterGet_2_100_000, N: 100_000},
		{Name: "Query.Next + 5x Get", Desc: "", F: queryIterGet_5_100_000, N: 100_000},

		{Name: "Query.EntityAt, 1 arch", Desc: "", F: querEntityAt_1Arch_1000, N: 1000},
		{Name: "Query.EntityAt, 1 arch", Desc: "registered filter", F: querEntityAtRegistered_1Arch_1000, N: 1000},
		{Name: "Query.EntityAt, 5 arch", Desc: "", F: querEntityAt_5Arch_1000, N: 1000},
		{Name: "Query.EntityAt, 5 arch", Desc: "registered filter", F: querEntityAtRegistered_5Arch_1000, N: 1000},

		{Name: "World.Query", Desc: "", F: queryCreate, N: 1},
		{Name: "World.Query", Desc: "registered filter", F: queryCreateCached, N: 1},
	}

	runBenches("Query", query, toMarkdown)
}

func runBenches(title string, benches []bench, format func([]bench) string) {
	for i := range benches {
		b := &benches[i]
		res := testing.Benchmark(b.F)
		b.T = float64(res.T.Nanoseconds()) / float64(res.N*b.N)
	}
	fmt.Printf("## %s\n\n%s", title, format(benches))
}

func toMarkdown(benches []bench) string {
	b := strings.Builder{}

	b.WriteString("| Operation | Time | Remark |\n")
	b.WriteString("|-----------|------|--------|\n")

	for i := range benches {
		bench := &benches[i]
		b.WriteString(fmt.Sprintf("| %s | %.1f ns | %s |\n", bench.Name, bench.T, bench.Desc))
	}

	return b.String()
}
