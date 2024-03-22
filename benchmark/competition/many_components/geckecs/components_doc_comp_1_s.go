package geckecs

type Comp1 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp1() bool {
	return e.w.comp1SStore.Has(e)
}

func (e Entity) ReadComp1() (Comp1, bool) {
	return e.w.comp1SStore.Read(e)
}

func (e Entity) RemoveComp1() Entity {
	e.w.comp1SStore.Remove(e)

	return e
}

func (e Entity) WritableComp1() (*Comp1, bool) {
	return e.w.comp1SStore.Writeable(e)
}

func (e Entity) SetComp1(other Comp1) Entity {
	e.w.comp1SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp1S(c Comp1, entities ...Entity) {
	w.comp1SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp1S(entities ...Entity) {
	w.comp1SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp1 checks if the world has a Comp1}}
func (w *World) HasComp1Resource() bool {
	return w.resourceEntity.HasComp1()
}

// Retrieve the Comp1 resource from the world
func (w *World) Comp1Resource() (Comp1, bool) {
	return w.resourceEntity.ReadComp1()
}

// Set the Comp1 resource in the world
func (w *World) SetComp1Resource(c Comp1) Entity {
	w.resourceEntity.SetComp1(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp1 resource from the world
func (w *World) RemoveComp1Resource() Entity {
	w.resourceEntity.RemoveComp1()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp1ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp1]
}

func (iter *Comp1ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp1ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp1ReadIterator) NextComp1() (Entity, Comp1) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp1ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp1ReadIter() *Comp1ReadIterator {
	iter := &Comp1ReadIterator{
		w:     w,
		store: w.comp1SStore,
	}
	iter.Reset()
	return iter
}

type Comp1WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp1]
}

func (iter *Comp1WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp1WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp1WriteIterator) NextComp1() (Entity, *Comp1) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp1WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp1WriteIter() *Comp1WriteIterator {
	iter := &Comp1WriteIterator{
		w:     w,
		store: w.comp1SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp1Entities() []Entity {
	return w.comp1SStore.entities()
}

func (w *World) SetComp1SortFn(lessThan func(a, b Entity) bool) {
	w.comp1SStore.LessThan = lessThan
}

func (w *World) SortComp1S() {
	w.comp1SStore.Sort()
}
