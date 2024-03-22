package geckecs

type Name string

func NameFromString(c string) Name {
	return Name(c)
}

func (c Name) ToString() string {
	return string(c)
}

//#region Events
//#endregion

func (e Entity) HasName() bool {
	return e.w.namesStore.Has(e)
}

func (e Entity) ReadName() (Name, bool) {
	return e.w.namesStore.Read(e)
}

func (e Entity) RemoveName() Entity {
	e.w.namesStore.Remove(e)

	return e
}

func (e Entity) WritableName() (*Name, bool) {
	return e.w.namesStore.Writeable(e)
}

func (e Entity) SetName(other Name) Entity {
	e.w.namesStore.Set(other, e)

	return e
}

func (w *World) SetNames(c Name, entities ...Entity) {
	w.namesStore.Set(c, entities...)
}

func (w *World) RemoveNames(entities ...Entity) {
	w.namesStore.Remove(entities...)
}

//#region Resources

// HasName checks if the world has a Name}}
func (w *World) HasNameResource() bool {
	return w.resourceEntity.HasName()
}

// Retrieve the Name resource from the world
func (w *World) NameResource() (Name, bool) {
	return w.resourceEntity.ReadName()
}

// Set the Name resource in the world
func (w *World) SetNameResource(c Name) Entity {
	w.resourceEntity.SetName(c)
	return w.resourceEntity
}

// Remove the Name resource from the world
func (w *World) RemoveNameResource() Entity {
	w.resourceEntity.RemoveName()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type NameReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Name]
}

func (iter *NameReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *NameReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *NameReadIterator) NextName() (Entity, Name) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *NameReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) NameReadIter() *NameReadIterator {
	iter := &NameReadIterator{
		w:     w,
		store: w.namesStore,
	}
	iter.Reset()
	return iter
}

type NameWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Name]
}

func (iter *NameWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *NameWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *NameWriteIterator) NextName() (Entity, *Name) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *NameWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) NameWriteIter() *NameWriteIterator {
	iter := &NameWriteIterator{
		w:     w,
		store: w.namesStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) NameEntities() []Entity {
	return w.namesStore.entities()
}

func (w *World) SetNameSortFn(lessThan func(a, b Entity) bool) {
	w.namesStore.LessThan = lessThan
}

func (w *World) SortNames() {
	w.namesStore.Sort()
}
