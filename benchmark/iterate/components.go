package iterate

import "github.com/mlange-42/arche/ecs"

type position struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

type testStruct0 struct{ val int32 }
type testStruct1 struct{ val int32 }
type testStruct2 struct{ val int32 }
type testStruct3 struct{ val int32 }
type testStruct4 struct{ val int32 }
type testStruct5 struct{ val int32 }
type testStruct6 struct{ val int32 }
type testStruct7 struct{ val int32 }
type testStruct8 struct{ val int32 }
type testStruct9 struct{ val int32 }

func registerAll(w *ecs.World) []ecs.ID {
	_ = testStruct0{1}
	_ = testStruct1{1}
	_ = testStruct2{1}
	_ = testStruct3{1}
	_ = testStruct4{1}
	_ = testStruct5{1}
	_ = testStruct6{1}
	_ = testStruct7{1}
	_ = testStruct8{1}
	_ = testStruct9{1}

	ids := make([]ecs.ID, 10)
	ids[0] = ecs.ComponentID[testStruct0](w)
	ids[1] = ecs.ComponentID[testStruct1](w)
	ids[2] = ecs.ComponentID[testStruct2](w)
	ids[3] = ecs.ComponentID[testStruct3](w)
	ids[4] = ecs.ComponentID[testStruct4](w)
	ids[5] = ecs.ComponentID[testStruct5](w)
	ids[6] = ecs.ComponentID[testStruct6](w)
	ids[7] = ecs.ComponentID[testStruct7](w)
	ids[8] = ecs.ComponentID[testStruct8](w)
	ids[9] = ecs.ComponentID[testStruct9](w)

	return ids
}
