package geckecs

type Comp2 struct{}

//#region Events
//#endregion

func (e Entity) HasComp2Tag() bool {
	return e.w.comp2Store.Has(e)
}

func (e Entity) TagWithComp2() Entity {
	e.w.comp2Store.Set(e.w.comp2Store.zero, e)
	e.w.Comp1Comp2Comp3Set.PossibleUpdate(e)
	return e
}

func (e Entity) RemoveComp2Tag() Entity {
	e.w.comp2Store.Remove(e)
	e.w.Comp1Comp2Comp3Set.PossibleUpdate(e)
	return e
}

func (w *World) RemoveComp2Tags(entities ...Entity) {
	w.comp2Store.Remove(entities...)
	w.Comp1Comp2Comp3Set.PossibleUpdate(entities...)
}

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

func (iter *Comp2ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp2ReadIter() *Comp2ReadIterator {
	iter := &Comp2ReadIterator{
		w:     w,
		store: w.comp2Store,
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

func (iter *Comp2WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp2WriteIter() *Comp2WriteIterator {
	iter := &Comp2WriteIterator{
		w:     w,
		store: w.comp2Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp2Entities() []Entity {
	return w.comp2Store.entities()
}
