package geckecs

type Comp3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp3() bool {
	return e.w.comp3SStore.Has(e)
}

func (e Entity) ReadComp3() (Comp3, bool) {
	return e.w.comp3SStore.Read(e)
}

func (e Entity) RemoveComp3() Entity {
	e.w.comp3SStore.Remove(e)

	return e
}

func (e Entity) WritableComp3() (*Comp3, bool) {
	return e.w.comp3SStore.Writeable(e)
}

func (e Entity) SetComp3(other Comp3) Entity {
	e.w.comp3SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp3S(c Comp3, entities ...Entity) {
	w.comp3SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp3S(entities ...Entity) {
	w.comp3SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp3 checks if the world has a Comp3}}
func (w *World) HasComp3Resource() bool {
	return w.resourceEntity.HasComp3()
}

// Retrieve the Comp3 resource from the world
func (w *World) Comp3Resource() (Comp3, bool) {
	return w.resourceEntity.ReadComp3()
}

// Set the Comp3 resource in the world
func (w *World) SetComp3Resource(c Comp3) Entity {
	w.resourceEntity.SetComp3(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp3 resource from the world
func (w *World) RemoveComp3Resource() Entity {
	w.resourceEntity.RemoveComp3()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp3ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp3]
}

func (iter *Comp3ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp3ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp3ReadIterator) NextComp3() (Entity, Comp3) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp3ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp3ReadIter() *Comp3ReadIterator {
	iter := &Comp3ReadIterator{
		w:     w,
		store: w.comp3SStore,
	}
	iter.Reset()
	return iter
}

type Comp3WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp3]
}

func (iter *Comp3WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp3WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp3WriteIterator) NextComp3() (Entity, *Comp3) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp3WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp3WriteIter() *Comp3WriteIterator {
	iter := &Comp3WriteIterator{
		w:     w,
		store: w.comp3SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp3Entities() []Entity {
	return w.comp3SStore.entities()
}

func (w *World) SetComp3SortFn(lessThan func(a, b Entity) bool) {
	w.comp3SStore.LessThan = lessThan
}

func (w *World) SortComp3S() {
	w.comp3SStore.Sort()
}
