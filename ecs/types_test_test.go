package ecs

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
