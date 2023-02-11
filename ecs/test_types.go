package ecs

type label struct{}

type position struct {
	X int
	Y int
}

type velocity struct {
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

func registerAll(w *World) []ID {
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

	ids := make([]ID, 10)
	ids[0] = ComponentID[testStruct0](w)
	ids[1] = ComponentID[testStruct1](w)
	ids[2] = ComponentID[testStruct2](w)
	ids[3] = ComponentID[testStruct3](w)
	ids[4] = ComponentID[testStruct4](w)
	ids[5] = ComponentID[testStruct5](w)
	ids[6] = ComponentID[testStruct6](w)
	ids[7] = ComponentID[testStruct7](w)
	ids[8] = ComponentID[testStruct8](w)
	ids[9] = ComponentID[testStruct9](w)

	return ids
}
