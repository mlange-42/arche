package geckecs

type IsA Entity

func IsAFromEntity(c Entity) IsA {
	return IsA(c)
}

func (c IsA) ToEntity() Entity {
	return Entity(c)
}

func (c IsA) FromEntity(e Entity) IsA {
	return IsA(e)
}

//#region Events
//#endregion

func (e Entity) HasIsA() bool {
	return e.w.isAStore.Has(e)
}

func (e Entity) ReadIsA() (Entity, bool) {
	val, ok := e.w.isAStore.Read(e)
	if !ok {
		return Entity{}, false
	}
	return Entity(val), true
}

func (e Entity) RemoveIsA() Entity {
	e.w.isAStore.Remove(e)

	return e
}

func (e Entity) WritableIsA() (*IsA, bool) {
	return e.w.isAStore.Writeable(e)
}

func (e Entity) SetIsA(other Entity) Entity {
	e.w.isAStore.Set(IsA(other), e)

	return e
}

func (w *World) SetIsA(c IsA, entities ...Entity) {
	w.isAStore.Set(c, entities...)
}

func (w *World) RemoveIsA(entities ...Entity) {
	w.isAStore.Remove(entities...)
}

//#region Resources

// HasIsA checks if the world has a IsA}}
func (w *World) HasIsAResource() bool {
	return w.resourceEntity.HasIsA()
}

// Retrieve the IsA resource from the world
func (w *World) IsAResource() (Entity, bool) {
	return w.resourceEntity.ReadIsA()
}

// Set the IsA resource in the world
func (w *World) SetIsAResource(c Entity) Entity {
	w.resourceEntity.SetIsA(c)

	return w.resourceEntity
}

// Remove the IsA resource from the world
func (w *World) RemoveIsAResource() Entity {
	w.resourceEntity.RemoveIsA()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type IsAReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[IsA]
}

func (iter *IsAReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *IsAReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *IsAReadIterator) NextIsA() (Entity, IsA) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *IsAReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) IsAReadIter() *IsAReadIterator {
	iter := &IsAReadIterator{
		w:     w,
		store: w.isAStore,
	}
	iter.Reset()
	return iter
}

type IsAWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[IsA]
}

func (iter *IsAWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *IsAWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *IsAWriteIterator) NextIsA() (Entity, *IsA) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *IsAWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) IsAWriteIter() *IsAWriteIterator {
	iter := &IsAWriteIterator{
		w:     w,
		store: w.isAStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) IsAEntities() []Entity {
	return w.isAStore.entities()
}

func (w *World) SetIsASortFn(lessThan func(a, b Entity) bool) {
	w.isAStore.LessThan = lessThan
}

func (w *World) SortIsA() {
	w.isAStore.Sort()
}
