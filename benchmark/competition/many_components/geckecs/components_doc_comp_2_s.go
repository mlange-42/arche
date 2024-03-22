package geckecs

type Comp2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp2() bool {
	return e.w.comp2SStore.Has(e)
}

func (e Entity) ReadComp2() (Comp2, bool) {
	return e.w.comp2SStore.Read(e)
}

func (e Entity) RemoveComp2() Entity {
	e.w.comp2SStore.Remove(e)

	return e
}

func (e Entity) WritableComp2() (*Comp2, bool) {
	return e.w.comp2SStore.Writeable(e)
}

func (e Entity) SetComp2(other Comp2) Entity {
	e.w.comp2SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp2S(c Comp2, entities ...Entity) {
	w.comp2SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp2S(entities ...Entity) {
	w.comp2SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp2 checks if the world has a Comp2}}
func (w *World) HasComp2Resource() bool {
	return w.resourceEntity.HasComp2()
}

// Retrieve the Comp2 resource from the world
func (w *World) Comp2Resource() (Comp2, bool) {
	return w.resourceEntity.ReadComp2()
}

// Set the Comp2 resource in the world
func (w *World) SetComp2Resource(c Comp2) Entity {
	w.resourceEntity.SetComp2(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp2 resource from the world
func (w *World) RemoveComp2Resource() Entity {
	w.resourceEntity.RemoveComp2()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp2ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp2]
}

func (iter *Comp2ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp2ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp2ReadIterator) NextComp2() (Entity, Comp2) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp2ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp2ReadIter() *Comp2ReadIterator {
	iter := &Comp2ReadIterator{
		w:     w,
		store: w.comp2SStore,
	}
	iter.Reset()
	return iter
}

type Comp2WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp2]
}

func (iter *Comp2WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp2WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp2WriteIterator) NextComp2() (Entity, *Comp2) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp2WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp2WriteIter() *Comp2WriteIterator {
	iter := &Comp2WriteIterator{
		w:     w,
		store: w.comp2SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp2Entities() []Entity {
	return w.comp2SStore.entities()
}

func (w *World) SetComp2SortFn(lessThan func(a, b Entity) bool) {
	w.comp2SStore.LessThan = lessThan
}

func (w *World) SortComp2S() {
	w.comp2SStore.Sort()
}
