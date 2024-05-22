package geckecs

type ChildOf Entity

func ChildOfFromEntity(c Entity) ChildOf {
	return ChildOf(c)
}

func (c ChildOf) ToEntity() Entity {
	return Entity(c)
}

func (c ChildOf) FromEntity(e Entity) ChildOf {
	return ChildOf(e)
}

//#region Events
//#endregion

func (e Entity) HasChildOf() bool {
	return e.w.childOfStore.Has(e)
}

func (e Entity) ReadChildOf() (Entity, bool) {
	val, ok := e.w.childOfStore.Read(e)
	if !ok {
		return Entity{}, false
	}
	return Entity(val), true
}

func (e Entity) RemoveChildOf() Entity {
	e.w.childOfStore.Remove(e)

	return e
}

func (e Entity) WritableChildOf() (*ChildOf, bool) {
	return e.w.childOfStore.Writeable(e)
}

func (e Entity) SetChildOf(other Entity) Entity {
	e.w.childOfStore.Set(ChildOf(other), e)

	return e
}

func (w *World) SetChildOf(c ChildOf, entities ...Entity) {
	w.childOfStore.Set(c, entities...)
}

func (w *World) RemoveChildOf(entities ...Entity) {
	w.childOfStore.Remove(entities...)
}

//#region Resources

// HasChildOf checks if the world has a ChildOf}}
func (w *World) HasChildOfResource() bool {
	return w.resourceEntity.HasChildOf()
}

// Retrieve the ChildOf resource from the world
func (w *World) ChildOfResource() (Entity, bool) {
	return w.resourceEntity.ReadChildOf()
}

// Set the ChildOf resource in the world
func (w *World) SetChildOfResource(c Entity) Entity {
	w.resourceEntity.SetChildOf(c)

	return w.resourceEntity
}

// Remove the ChildOf resource from the world
func (w *World) RemoveChildOfResource() Entity {
	w.resourceEntity.RemoveChildOf()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type ChildOfReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[ChildOf]
}

func (iter *ChildOfReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *ChildOfReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *ChildOfReadIterator) NextChildOf() (Entity, ChildOf) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *ChildOfReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) ChildOfReadIter() *ChildOfReadIterator {
	iter := &ChildOfReadIterator{
		w:     w,
		store: w.childOfStore,
	}
	iter.Reset()
	return iter
}

type ChildOfWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[ChildOf]
}

func (iter *ChildOfWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *ChildOfWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *ChildOfWriteIterator) NextChildOf() (Entity, *ChildOf) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *ChildOfWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) ChildOfWriteIter() *ChildOfWriteIterator {
	iter := &ChildOfWriteIterator{
		w:     w,
		store: w.childOfStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) ChildOfEntities() []Entity {
	return w.childOfStore.entities()
}

func (w *World) SetChildOfSortFn(lessThan func(a, b Entity) bool) {
	w.childOfStore.LessThan = lessThan
}

func (w *World) SortChildOf() {
	w.childOfStore.Sort()
}
