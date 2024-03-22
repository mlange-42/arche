package geckecs

type Comp3 struct{}

//#region Events
//#endregion

func (e Entity) HasComp3Tag() bool {
	return e.w.comp3Store.Has(e)
}

func (e Entity) TagWithComp3() Entity {
	e.w.comp3Store.Set(e.w.comp3Store.zero, e)
	e.w.Comp1Comp2Comp3Set.PossibleUpdate(e)
	return e
}

func (e Entity) RemoveComp3Tag() Entity {
	e.w.comp3Store.Remove(e)
	e.w.Comp1Comp2Comp3Set.PossibleUpdate(e)
	return e
}

func (w *World) RemoveComp3Tags(entities ...Entity) {
	w.comp3Store.Remove(entities...)
	w.Comp1Comp2Comp3Set.PossibleUpdate(entities...)
}

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

func (iter *Comp3ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp3ReadIter() *Comp3ReadIterator {
	iter := &Comp3ReadIterator{
		w:     w,
		store: w.comp3Store,
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

func (iter *Comp3WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp3WriteIter() *Comp3WriteIterator {
	iter := &Comp3WriteIterator{
		w:     w,
		store: w.comp3Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp3Entities() []Entity {
	return w.comp3Store.entities()
}
