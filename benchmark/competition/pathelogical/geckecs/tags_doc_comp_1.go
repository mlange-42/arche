package geckecs

type Comp1 struct{}

//#region Events
//#endregion

func (e Entity) HasComp1Tag() bool {
	return e.w.comp1Store.Has(e)
}

func (e Entity) TagWithComp1() Entity {
	e.w.comp1Store.Set(e.w.comp1Store.zero, e)
	e.w.Comp1Comp2Comp3Set.PossibleUpdate(e)
	return e
}

func (e Entity) RemoveComp1Tag() Entity {
	e.w.comp1Store.Remove(e)
	e.w.Comp1Comp2Comp3Set.PossibleUpdate(e)
	return e
}

func (w *World) RemoveComp1Tags(entities ...Entity) {
	w.comp1Store.Remove(entities...)
	w.Comp1Comp2Comp3Set.PossibleUpdate(entities...)
}

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

func (iter *Comp1ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp1ReadIter() *Comp1ReadIterator {
	iter := &Comp1ReadIterator{
		w:     w,
		store: w.comp1Store,
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

func (iter *Comp1WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp1WriteIter() *Comp1WriteIterator {
	iter := &Comp1WriteIterator{
		w:     w,
		store: w.comp1Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp1Entities() []Entity {
	return w.comp1Store.entities()
}
