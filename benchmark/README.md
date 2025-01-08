# Benchmarks

Arche's [benchmarks](./arche/) and [profiling](./profile/).

See also the internal benchmarks in packages [ecs](../ecs/), [generic](../generic/) and [filter](../filter/).

# Running the benchmarks

From here (`benchmarks/`), run:

```
go test -benchmem -run=^$ -bench ^.*$ ./...
```
