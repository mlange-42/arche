package geckecs

type Comp7 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp7() bool {
	return e.w.comp7SStore.Has(e)
}

func (e Entity) ReadComp7() (Comp7, bool) {
	return e.w.comp7SStore.Read(e)
}

func (e Entity) RemoveComp7() Entity {
	e.w.comp7SStore.Remove(e)

	return e
}

func (e Entity) WritableComp7() (*Comp7, bool) {
	return e.w.comp7SStore.Writeable(e)
}

func (e Entity) SetComp7(other Comp7) Entity {
	e.w.comp7SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp7S(c Comp7, entities ...Entity) {
	w.comp7SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp7S(entities ...Entity) {
	w.comp7SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp7 checks if the world has a Comp7}}
func (w *World) HasComp7Resource() bool {
	return w.resourceEntity.HasComp7()
}

// Retrieve the Comp7 resource from the world
func (w *World) Comp7Resource() (Comp7, bool) {
	return w.resourceEntity.ReadComp7()
}

// Set the Comp7 resource in the world
func (w *World) SetComp7Resource(c Comp7) Entity {
	w.resourceEntity.SetComp7(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp7 resource from the world
func (w *World) RemoveComp7Resource() Entity {
	w.resourceEntity.RemoveComp7()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp7ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp7]
}

func (iter *Comp7ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp7ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp7ReadIterator) NextComp7() (Entity, Comp7) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp7ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp7ReadIter() *Comp7ReadIterator {
	iter := &Comp7ReadIterator{
		w:     w,
		store: w.comp7SStore,
	}
	iter.Reset()
	return iter
}

type Comp7WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp7]
}

func (iter *Comp7WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp7WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp7WriteIterator) NextComp7() (Entity, *Comp7) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp7WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp7WriteIter() *Comp7WriteIterator {
	iter := &Comp7WriteIterator{
		w:     w,
		store: w.comp7SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp7Entities() []Entity {
	return w.comp7SStore.entities()
}

func (w *World) SetComp7SortFn(lessThan func(a, b Entity) bool) {
	w.comp7SStore.LessThan = lessThan
}

func (w *World) SortComp7S() {
	w.comp7SStore.Sort()
}
