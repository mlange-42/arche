package geckecs

type Comp4 struct{}

//#region Events
//#endregion

func (e Entity) HasComp4Tag() bool {
	return e.w.comp4Store.Has(e)
}

func (e Entity) TagWithComp4() Entity {
	e.w.comp4Store.Set(e.w.comp4Store.zero, e)
	return e
}

func (e Entity) RemoveComp4Tag() Entity {
	e.w.comp4Store.Remove(e)
	return e
}

func (w *World) RemoveComp4Tags(entities ...Entity) {
	w.comp4Store.Remove(entities...)
}

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

func (iter *Comp4ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp4ReadIter() *Comp4ReadIterator {
	iter := &Comp4ReadIterator{
		w:     w,
		store: w.comp4Store,
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

func (iter *Comp4WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp4WriteIter() *Comp4WriteIterator {
	iter := &Comp4WriteIterator{
		w:     w,
		store: w.comp4Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp4Entities() []Entity {
	return w.comp4Store.entities()
}
