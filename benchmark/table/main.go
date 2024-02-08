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
	benches := []bench{
		{Name: "Iter", Desc: "", F: queryIter_100_000, N: 100_000},
		{Name: "Iter + Get 1", Desc: "", F: queryIterGet_1_100_000, N: 100_000},
		{Name: "Iter + Get 2", Desc: "", F: queryIterGet_2_100_000, N: 100_000},
		{Name: "Iter + Get 5", Desc: "", F: queryIterGet_5_100_000, N: 100_000},
	}

	for i := range benches {
		b := &benches[i]
		res := testing.Benchmark(b.F)
		b.T = float64(res.T.Nanoseconds()) / float64(res.N*b.N)
	}

	md := toMarkdown(benches)
	fmt.Println(md)
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
