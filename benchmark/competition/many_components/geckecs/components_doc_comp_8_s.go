package geckecs

type Comp8 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp8() bool {
	return e.w.comp8SStore.Has(e)
}

func (e Entity) ReadComp8() (Comp8, bool) {
	return e.w.comp8SStore.Read(e)
}

func (e Entity) RemoveComp8() Entity {
	e.w.comp8SStore.Remove(e)

	return e
}

func (e Entity) WritableComp8() (*Comp8, bool) {
	return e.w.comp8SStore.Writeable(e)
}

func (e Entity) SetComp8(other Comp8) Entity {
	e.w.comp8SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp8S(c Comp8, entities ...Entity) {
	w.comp8SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp8S(entities ...Entity) {
	w.comp8SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp8 checks if the world has a Comp8}}
func (w *World) HasComp8Resource() bool {
	return w.resourceEntity.HasComp8()
}

// Retrieve the Comp8 resource from the world
func (w *World) Comp8Resource() (Comp8, bool) {
	return w.resourceEntity.ReadComp8()
}

// Set the Comp8 resource in the world
func (w *World) SetComp8Resource(c Comp8) Entity {
	w.resourceEntity.SetComp8(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp8 resource from the world
func (w *World) RemoveComp8Resource() Entity {
	w.resourceEntity.RemoveComp8()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp8ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp8]
}

func (iter *Comp8ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp8ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp8ReadIterator) NextComp8() (Entity, Comp8) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp8ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp8ReadIter() *Comp8ReadIterator {
	iter := &Comp8ReadIterator{
		w:     w,
		store: w.comp8SStore,
	}
	iter.Reset()
	return iter
}

type Comp8WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp8]
}

func (iter *Comp8WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp8WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp8WriteIterator) NextComp8() (Entity, *Comp8) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp8WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp8WriteIter() *Comp8WriteIterator {
	iter := &Comp8WriteIterator{
		w:     w,
		store: w.comp8SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp8Entities() []Entity {
	return w.comp8SStore.entities()
}

func (w *World) SetComp8SortFn(lessThan func(a, b Entity) bool) {
	w.comp8SStore.LessThan = lessThan
}

func (w *World) SortComp8S() {
	w.comp8SStore.Sort()
}
