package main

import "github.com/mlange-42/arche/ecs"

type comp1 struct {
	V int64
	W int64
}

type comp2 struct {
	V int64
	W int64
}

type comp3 struct {
	V int64
	W int64
}

type comp4 struct {
	V int64
	W int64
}

type comp5 struct {
	V int64
	W int64
}

type relComp1 struct {
	ecs.Relation
}
