package geckecs

type Comp6 struct{}

//#region Events
//#endregion

func (e Entity) HasComp6Tag() bool {
	return e.w.comp6Store.Has(e)
}

func (e Entity) TagWithComp6() Entity {
	e.w.comp6Store.Set(e.w.comp6Store.zero, e)
	return e
}

func (e Entity) RemoveComp6Tag() Entity {
	e.w.comp6Store.Remove(e)
	return e
}

func (w *World) RemoveComp6Tags(entities ...Entity) {
	w.comp6Store.Remove(entities...)
}

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

func (iter *Comp6ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp6ReadIter() *Comp6ReadIterator {
	iter := &Comp6ReadIterator{
		w:     w,
		store: w.comp6Store,
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

func (iter *Comp6WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp6WriteIter() *Comp6WriteIterator {
	iter := &Comp6WriteIterator{
		w:     w,
		store: w.comp6Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp6Entities() []Entity {
	return w.comp6Store.entities()
}
