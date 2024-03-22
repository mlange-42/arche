package geckecs

type Comp9 struct{}

//#region Events
//#endregion

func (e Entity) HasComp9Tag() bool {
	return e.w.comp9Store.Has(e)
}

func (e Entity) TagWithComp9() Entity {
	e.w.comp9Store.Set(e.w.comp9Store.zero, e)
	return e
}

func (e Entity) RemoveComp9Tag() Entity {
	e.w.comp9Store.Remove(e)
	return e
}

func (w *World) RemoveComp9Tags(entities ...Entity) {
	w.comp9Store.Remove(entities...)
}

//#region Iterators

type Comp9ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp9]
}

func (iter *Comp9ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp9ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp9ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp9ReadIter() *Comp9ReadIterator {
	iter := &Comp9ReadIterator{
		w:     w,
		store: w.comp9Store,
	}
	iter.Reset()
	return iter
}

type Comp9WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp9]
}

func (iter *Comp9WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp9WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp9WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp9WriteIter() *Comp9WriteIterator {
	iter := &Comp9WriteIterator{
		w:     w,
		store: w.comp9Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp9Entities() []Entity {
	return w.comp9Store.entities()
}
