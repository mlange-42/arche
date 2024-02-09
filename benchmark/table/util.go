package main

import (
	"fmt"
	"strings"
	"testing"
)

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

	b.WriteString(fmt.Sprintf("| %-32s | %-12s | %-28s |\n", "Operation", "Time", "Remark"))
	b.WriteString(fmt.Sprintf("|%s|%s:|%s|\n", strings.Repeat("-", 34), strings.Repeat("-", 13), strings.Repeat("-", 30)))

	for i := range benches {
		bench := &benches[i]
		t := fmt.Sprintf("%.1f ns", bench.T)
		b.WriteString(fmt.Sprintf("| %-32s | %12s | %-28s |\n", bench.Name, t, bench.Desc))
	}
	b.WriteString("\n")

	return b.String()
}
