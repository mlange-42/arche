package geckecs

type Comp8 struct{}

//#region Events
//#endregion

func (e Entity) HasComp8Tag() bool {
	return e.w.comp8Store.Has(e)
}

func (e Entity) TagWithComp8() Entity {
	e.w.comp8Store.Set(e.w.comp8Store.zero, e)
	return e
}

func (e Entity) RemoveComp8Tag() Entity {
	e.w.comp8Store.Remove(e)
	return e
}

func (w *World) RemoveComp8Tags(entities ...Entity) {
	w.comp8Store.Remove(entities...)
}

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

func (iter *Comp8ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp8ReadIter() *Comp8ReadIterator {
	iter := &Comp8ReadIterator{
		w:     w,
		store: w.comp8Store,
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

func (iter *Comp8WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp8WriteIter() *Comp8WriteIterator {
	iter := &Comp8WriteIterator{
		w:     w,
		store: w.comp8Store,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp8Entities() []Entity {
	return w.comp8Store.entities()
}
