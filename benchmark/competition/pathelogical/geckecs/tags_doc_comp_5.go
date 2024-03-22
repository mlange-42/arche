package geckecs

type Comp5 struct{}

//#region Events
//#endregion

func (e Entity) HasComp5Tag() bool {
	return e.w.comp5Store.Has(e)
}

func (e Entity) TagWithComp5() Entity {
	e.w.comp5Store.Set(e.w.comp5Store.zero, e)
	return e
}

func (e Entity) RemoveComp5Tag() Entity {
	e.w.comp5Store.Remove(e)
	return e
}

func (w *World) RemoveComp5Tags(entities ...Entity) {
	w.comp5Store.Remove(entities...)
}

//#region Iterators

type Comp5ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp5]
}

func (iter *Comp5ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp5ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp5ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp5ReadIter() *Comp5ReadIterator {
	iter := &Comp5ReadIterator{
		w:     w,
		store: w.comp5Store,
	}
	iter.Reset()
	return iter
}

type Comp5WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp5]
}

func (iter *Comp5WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp5WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp5WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp5WriteIter() *Comp5WriteIterator {
	iter := &Comp5WriteIterator{
		w:     w,
		store: w.comp5Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp5Entities() []Entity {
	return w.comp5Store.entities()
}
