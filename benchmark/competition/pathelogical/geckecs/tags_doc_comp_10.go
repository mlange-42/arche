package geckecs

type Comp10 struct{}

//#region Events
//#endregion

func (e Entity) HasComp10Tag() bool {
	return e.w.comp10Store.Has(e)
}

func (e Entity) TagWithComp10() Entity {
	e.w.comp10Store.Set(e.w.comp10Store.zero, e)
	return e
}

func (e Entity) RemoveComp10Tag() Entity {
	e.w.comp10Store.Remove(e)
	return e
}

func (w *World) RemoveComp10Tags(entities ...Entity) {
	w.comp10Store.Remove(entities...)
}

//#region Iterators

type Comp10ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp10]
}

func (iter *Comp10ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp10ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp10ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp10ReadIter() *Comp10ReadIterator {
	iter := &Comp10ReadIterator{
		w:     w,
		store: w.comp10Store,
	}
	iter.Reset()
	return iter
}

type Comp10WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp10]
}

func (iter *Comp10WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp10WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp10WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp10WriteIter() *Comp10WriteIterator {
	iter := &Comp10WriteIterator{
		w:     w,
		store: w.comp10Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp10Entities() []Entity {
	return w.comp10Store.entities()
}
