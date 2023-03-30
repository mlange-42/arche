# Benchmarks

[Arche benchmarks](./arche/) and [competition](./competition/) against other Go ECS implementations. [Profiling](./profile/).

See also the internal benchmarks in packages [ecs](../ecs/), [generic](../generic/) and [filter](../filter/).

# Running the benchmarks

From here (`benchmarks/`), run:

```
go test -benchmem -run=^$ -bench ^.*$ ./...
```
