package generic

import "github.com/mlange-42/arche/ecs"

type Position struct {
	X int
	Y int
}

type Velocity struct {
	X int
	Y int
}

type testRelationA struct {
	ecs.Relation
}

type testRelationB struct {
	ecs.Relation
}
