package ecs

import (
	"fmt"
	"reflect"
	"testing"
)

type label struct{}

type Position struct {
	X int
	Y int
}

type Velocity struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

type testRelationA struct {
	Relation
}

type testRelationB struct {
	Relation
}

type ChildOf struct {
	Relation
}

type testStruct0 struct{ Val int32 }
type testStruct1 struct{ val int32 }
type testStruct2 struct{ val int32 }
type testStruct3 struct{ val int32 }
type testStruct4 struct{ val int32 }
type testStruct5 struct{ val int32 }
type testStruct6 struct{ val int32 }
type testStruct7 struct{ val int32 }
type testStruct8 struct{ val int32 }
type testStruct9 struct{ val int32 }
type testStruct10 struct{ val int32 }
type testStruct11 struct{ val int32 }
type testStruct12 struct{ val int32 }
type testStruct13 struct{ val int32 }
type testStruct14 struct{ val int32 }
type testStruct15 struct{ val int32 }
type testStruct16 struct{ val int32 }

func TestTypeSizes(t *testing.T) {
	printTypeSize[Entity]()
	printTypeSize[entityIndex]()
	printTypeSize[Mask]()
	printTypeSize[World]()
	printTypeSizeName[pagedSlice[archetype]]("pagedArr32")
	printTypeSize[archetype]()
	printTypeSize[archetypeAccess]()
	printTypeSize[archetypeData]()
	printTypeSize[archNode]()
	printTypeSize[nodeData]()
	printTypeSize[layout]()
	printTypeSize[entityPool]()
	printTypeSize[componentRegistry]()
	printTypeSize[bitPool]()
	printTypeSize[Query]()
	printTypeSize[Resources]()
	printTypeSizeName[reflect.Value]("reflect.Value")
	printTypeSize[EntityEvent]()
	printTypeSize[Cache]()
	printTypeSizeName[idMap[uint32]]("idMap")
}

func printTypeSize[T any]() {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", tp.Name(), tp.Size())
}

func printTypeSizeName[T any](name string) {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", name, tp.Size())
}
