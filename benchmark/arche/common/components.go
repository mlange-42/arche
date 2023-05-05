package common

import "github.com/mlange-42/arche/ecs"

type Position struct {
	X int
	Y int
}

type Rotation struct {
	Angle int
}

type ChildOf struct {
	ecs.Relation
}

type TestStruct0 struct{ Val int32 }
type TestStruct1 struct{ Val int32 }
type TestStruct2 struct{ Val int32 }
type TestStruct3 struct{ Val int32 }
type TestStruct4 struct{ Val int32 }
type TestStruct5 struct{ Val int32 }
type TestStruct6 struct{ Val int32 }
type TestStruct7 struct{ Val int32 }
type TestStruct8 struct{ Val int32 }
type TestStruct9 struct{ Val int32 }
type TestStruct10 struct{ Val int32 }

func RegisterAll(w *ecs.World) []ecs.ID {
	_ = TestStruct0{1}
	_ = TestStruct1{1}
	_ = TestStruct2{1}
	_ = TestStruct3{1}
	_ = TestStruct4{1}
	_ = TestStruct5{1}
	_ = TestStruct6{1}
	_ = TestStruct7{1}
	_ = TestStruct8{1}
	_ = TestStruct9{1}
	_ = TestStruct10{1}

	ids := make([]ecs.ID, 11)
	ids[0] = ecs.ComponentID[TestStruct0](w)
	ids[1] = ecs.ComponentID[TestStruct1](w)
	ids[2] = ecs.ComponentID[TestStruct2](w)
	ids[3] = ecs.ComponentID[TestStruct3](w)
	ids[4] = ecs.ComponentID[TestStruct4](w)
	ids[5] = ecs.ComponentID[TestStruct5](w)
	ids[6] = ecs.ComponentID[TestStruct6](w)
	ids[7] = ecs.ComponentID[TestStruct7](w)
	ids[8] = ecs.ComponentID[TestStruct8](w)
	ids[9] = ecs.ComponentID[TestStruct9](w)
	ids[10] = ecs.ComponentID[TestStruct10](w)

	return ids
}
