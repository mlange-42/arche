package geckecs

type Comp6 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp6() bool {
	return e.w.comp6SStore.Has(e)
}

func (e Entity) ReadComp6() (Comp6, bool) {
	return e.w.comp6SStore.Read(e)
}

func (e Entity) RemoveComp6() Entity {
	e.w.comp6SStore.Remove(e)

	return e
}

func (e Entity) WritableComp6() (*Comp6, bool) {
	return e.w.comp6SStore.Writeable(e)
}

func (e Entity) SetComp6(other Comp6) Entity {
	e.w.comp6SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp6S(c Comp6, entities ...Entity) {
	w.comp6SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp6S(entities ...Entity) {
	w.comp6SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp6 checks if the world has a Comp6}}
func (w *World) HasComp6Resource() bool {
	return w.resourceEntity.HasComp6()
}

// Retrieve the Comp6 resource from the world
func (w *World) Comp6Resource() (Comp6, bool) {
	return w.resourceEntity.ReadComp6()
}

// Set the Comp6 resource in the world
func (w *World) SetComp6Resource(c Comp6) Entity {
	w.resourceEntity.SetComp6(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp6 resource from the world
func (w *World) RemoveComp6Resource() Entity {
	w.resourceEntity.RemoveComp6()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp6ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp6]
}

func (iter *Comp6ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp6ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp6ReadIterator) NextComp6() (Entity, Comp6) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp6ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp6ReadIter() *Comp6ReadIterator {
	iter := &Comp6ReadIterator{
		w:     w,
		store: w.comp6SStore,
	}
	iter.Reset()
	return iter
}

type Comp6WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp6]
}

func (iter *Comp6WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp6WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp6WriteIterator) NextComp6() (Entity, *Comp6) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp6WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp6WriteIter() *Comp6WriteIterator {
	iter := &Comp6WriteIterator{
		w:     w,
		store: w.comp6SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp6Entities() []Entity {
	return w.comp6SStore.entities()
}

func (w *World) SetComp6SortFn(lessThan func(a, b Entity) bool) {
	w.comp6SStore.LessThan = lessThan
}

func (w *World) SortComp6S() {
	w.comp6SStore.Sort()
}
