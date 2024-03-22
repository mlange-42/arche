package geckecs

type Comp4 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp4() bool {
	return e.w.comp4SStore.Has(e)
}

func (e Entity) ReadComp4() (Comp4, bool) {
	return e.w.comp4SStore.Read(e)
}

func (e Entity) RemoveComp4() Entity {
	e.w.comp4SStore.Remove(e)

	return e
}

func (e Entity) WritableComp4() (*Comp4, bool) {
	return e.w.comp4SStore.Writeable(e)
}

func (e Entity) SetComp4(other Comp4) Entity {
	e.w.comp4SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp4S(c Comp4, entities ...Entity) {
	w.comp4SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp4S(entities ...Entity) {
	w.comp4SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp4 checks if the world has a Comp4}}
func (w *World) HasComp4Resource() bool {
	return w.resourceEntity.HasComp4()
}

// Retrieve the Comp4 resource from the world
func (w *World) Comp4Resource() (Comp4, bool) {
	return w.resourceEntity.ReadComp4()
}

// Set the Comp4 resource in the world
func (w *World) SetComp4Resource(c Comp4) Entity {
	w.resourceEntity.SetComp4(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp4 resource from the world
func (w *World) RemoveComp4Resource() Entity {
	w.resourceEntity.RemoveComp4()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp4ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp4]
}

func (iter *Comp4ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp4ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp4ReadIterator) NextComp4() (Entity, Comp4) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp4ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp4ReadIter() *Comp4ReadIterator {
	iter := &Comp4ReadIterator{
		w:     w,
		store: w.comp4SStore,
	}
	iter.Reset()
	return iter
}

type Comp4WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp4]
}

func (iter *Comp4WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp4WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp4WriteIterator) NextComp4() (Entity, *Comp4) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp4WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp4WriteIter() *Comp4WriteIterator {
	iter := &Comp4WriteIterator{
		w:     w,
		store: w.comp4SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp4Entities() []Entity {
	return w.comp4SStore.entities()
}

func (w *World) SetComp4SortFn(lessThan func(a, b Entity) bool) {
	w.comp4SStore.LessThan = lessThan
}

func (w *World) SortComp4S() {
	w.comp4SStore.Sort()
}
