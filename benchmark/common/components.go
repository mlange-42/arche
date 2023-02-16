package common

import "github.com/mlange-42/arche/ecs"

type Position struct {
	X int
	Y int
}

type Rotation struct {
	Angle int
}

type TestStruct0 struct{ val int32 }
type TestStruct1 struct{ val int32 }
type TestStruct2 struct{ val int32 }
type TestStruct3 struct{ val int32 }
type TestStruct4 struct{ val int32 }
type TestStruct5 struct{ val int32 }
type TestStruct6 struct{ val int32 }
type TestStruct7 struct{ val int32 }
type TestStruct8 struct{ val int32 }
type TestStruct9 struct{ val int32 }

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

	ids := make([]ecs.ID, 10)
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

	return ids
}
