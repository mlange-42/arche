package geckecs

type Comp7 struct{}

//#region Events
//#endregion

func (e Entity) HasComp7Tag() bool {
	return e.w.comp7Store.Has(e)
}

func (e Entity) TagWithComp7() Entity {
	e.w.comp7Store.Set(e.w.comp7Store.zero, e)
	return e
}

func (e Entity) RemoveComp7Tag() Entity {
	e.w.comp7Store.Remove(e)
	return e
}

func (w *World) RemoveComp7Tags(entities ...Entity) {
	w.comp7Store.Remove(entities...)
}

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

func (iter *Comp7ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp7ReadIter() *Comp7ReadIterator {
	iter := &Comp7ReadIterator{
		w:     w,
		store: w.comp7Store,
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

func (iter *Comp7WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp7WriteIter() *Comp7WriteIterator {
	iter := &Comp7WriteIterator{
		w:     w,
		store: w.comp7Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp7Entities() []Entity {
	return w.comp7Store.entities()
}
