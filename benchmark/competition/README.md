# Competition

Comparative benchmarks against other Go ECS implementations, as well as [against Array of Structs](./array_of_structs/).

```
go test -benchmem -run="^$" -bench "^.*$" ./competition/...
```

### Modifying components

From `/arche/benchmark/competition/pos_vel` run:
```
go test -benchmem -run="^$" -bench "^.*$" . --count 5 
```

### Iterating entities

From `/arche/benchmark/competition/add_remove` run:
```
go test -benchmem -run="^$" -bench "^.*$" . --count 5 
```
